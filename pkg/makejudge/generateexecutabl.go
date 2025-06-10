package makejudge

import (
	"encoding/json"
	"fmt"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
	"strings"
)

const pythonExecutableTemplate = `{{.Import}}
import json
import sys
import time
from typing import Any, List, Dict

{{.UserCode}}

def run_single_test():
    test_case = {{.TestData}}
    try:
        start_time = time.time()
        input_data = json.loads(test_case["input"])
        result = main(input_data)
        actual_output = json.dumps(result) if not isinstance(result, str) else result
        result_data = {
            "name": test_case["name"],            
            "input": test_case["input"],
            "actual": actual_output,
            "error": "",
        }
    except Exception as e:
        result_data = {
            "name": test_case["name"],
            "input": test_case["input"],
            "actual": "",
            "error": str(e)
        }
    return result_data
if __name__ == "__main__":
    result = run_single_test()
    print(json.dumps(result, ensure_ascii=False, indent=2))
`

const cppExecutableTemplate = `{{.Import}}
#include <iostream>
#include <string>
#include <vector>
#include <chrono>
#include <nlohmann/json.hpp>

using json = nlohmann::json;
using namespace std;

{{.UserCode}}

int main() {
    json test_case = R"({{.TestData}})"_json;
    json result_data;
    
    try {
        auto start = chrono::high_resolution_clock::now();
        auto input_data = json::parse(test_case["input"].get<string>());
        auto result = solve(input_data);
        auto end = chrono::high_resolution_clock::now();
        auto duration = chrono::duration_cast<chrono::milliseconds>(end - start);
        
        json actual_json = result;
        string actual_output = actual_json.dump();
        
        result_data = {
            {"name", test_case["name"]},
            {"input", test_case["input"]},
            {"actual", actual_output},
            {"error", ""}
        };
        
    } catch (const exception& e) {
        result_data = {
            {"name", test_case["name"]},
            {"input", test_case["input"]},
            {"actual", ""},
            {"error", e.what()}
        };
    }
    
    cout << result_data.dump(2) << endl;
    return 0;
}
`

const jsExecutableTemplate = `{{.Import}}

{{.UserCode}}

function runSingleTest() {
    const testCase = {{.TestData}};
    
    try {
        const startTime = Date.now();
        const inputData = JSON.parse(testCase.input);
        const result = main(inputData);
        const actualOutput = typeof result === 'string' ? result : JSON.stringify(result);
        
        const resultData = {
            name: testCase.name,
            input: testCase.input,
            actual: actualOutput,
            error: ""
        };
        
        return resultData;
        
    } catch (error) {
        return {
            name: testCase.name,
            input: testCase.input,
            actual: "",
            error: error.message
        };
    }
}

const result = runSingleTest();
console.log(JSON.stringify(result, null, 2));
`

func GenerateExecutableCodeWithSettingData(testData SettingData, testCase loadjudge.TestCase, userCode string, language string) (string, error) {
	testCaseJSON, err := json.Marshal(testCase)
	if err != nil {
		return "", fmt.Errorf("序列化失敗: %w", err)
	}

	var inputType, outputType *loadjudge.TypeInfo
	inputType, outputType, err = loadjudge.ParseTestDataTypesString(testData.InputFormat.Type, testData.OutputFormat.Type)
	if err != nil {
		fmt.Printf("無法解析類型資訊: %v\n", err)
	}

	var template string
	var imports string

	switch strings.ToLower(language) {
	case loadjudge.Python:
		template = pythonExecutableTemplate
		imports = generateImportsForPython(inputType, outputType)
	case loadjudge.Cpp:
		template = cppExecutableTemplate
		imports = generateIncludesForCpp(inputType, outputType)
	case loadjudge.JavaScript:
		template = jsExecutableTemplate
		imports = ""
	default:
		return "", fmt.Errorf("不支援的語言: %s", language)
	}

	code := strings.ReplaceAll(template, "{{.Import}}", imports)
	code = strings.ReplaceAll(code, "{{.UserCode}}", userCode)
	code = strings.ReplaceAll(code, "{{.TestData}}", string(testCaseJSON))

	return code, nil
}

func GenerateExecutableCodeWithTestData(testData loadjudge.TestDataConfig, testCase loadjudge.TestCase, userCode string, language string) (string, error) {
	testCaseJSON, err := json.Marshal(testCase)
	if err != nil {
		return "", fmt.Errorf("序列化失敗: %w", err)
	}

	var inputType, outputType *loadjudge.TypeInfo
	inputType, outputType, err = loadjudge.ParseTestDataTypes(testData)
	if err != nil {
		fmt.Printf("無法解析類型資訊: %v\n", err)
	}

	var template string
	var imports string

	switch strings.ToLower(language) {
	case loadjudge.Python:
		template = pythonExecutableTemplate
		imports = generateImportsForPython(inputType, outputType)
	case loadjudge.Cpp:
		template = cppExecutableTemplate
		imports = generateIncludesForCpp(inputType, outputType)
	case loadjudge.JavaScript:
		template = jsExecutableTemplate
		imports = ""
	default:
		return "", fmt.Errorf("不支援的語言: %s", language)
	}

	code := strings.ReplaceAll(template, "{{.Import}}", imports)
	code = strings.ReplaceAll(code, "{{.UserCode}}", userCode)
	code = strings.ReplaceAll(code, "{{.TestData}}", string(testCaseJSON))

	return code, nil
}

func generateImportsForPython(inputType, outputType *loadjudge.TypeInfo) string {
	imports := make(map[string]bool)

	if inputType != nil {
		for _, imp := range inputType.GetRequiredImports(loadjudge.Python) {
			imports[imp] = true
		}
	}

	if outputType != nil {
		for _, imp := range outputType.GetRequiredImports(loadjudge.Python) {
			imports[imp] = true
		}
	}

	var importLines []string
	for imp := range imports {
		importLines = append(importLines, imp)
	}

	return strings.Join(importLines, "\n")
}

func generateIncludesForCpp(inputType, outputType *loadjudge.TypeInfo) string {
	includes := map[string]bool{
		"#include <iostream>": true,
		"#include <vector>":   true,
		"#include <string>":   true,
		"#include <chrono>":   true,
	}

	if inputType != nil {
		for _, inc := range inputType.GetRequiredImports(loadjudge.Cpp) {
			includes[fmt.Sprintf("#include %s", inc)] = true
		}
	}

	if outputType != nil {
		for _, inc := range outputType.GetRequiredImports(loadjudge.Cpp) {
			includes[fmt.Sprintf("#include %s", inc)] = true
		}
	}

	var includeLines []string
	for inc := range includes {
		includeLines = append(includeLines, inc)
	}

	return strings.Join(includeLines, "\n")
}
