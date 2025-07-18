package makejudge

import (
	"encoding/json"
	"fmt"
	langtype "share-cat/Nepeta-cataria/pkg/language"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
	"strings"
)

const pythonExecutableTemplate = `{{.Import}}
import json
import sys
from typing import Any, List, Dict

{{.UserCode}}

def run_single_test():
    test_case = {{.TestData}}
    try:
        input_data = json.loads(test_case["input"])
        result = main(input_data)
        if isinstance(result, str):
            print(result)
        else:
            print(json.dumps(result))
    except Exception as e:
        print(e, file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    run_single_test()
`

const cppExecutableTemplate = `{{.Import}}
#include <iostream>
#include <string>
#include <vector>
#include <nlohmann/json.hpp>

using json = nlohmann::json;
using namespace std;

{{.UserCode}}

int main() {
    json test_case = R"({{.TestData}})"_json;
    
    try {
        auto input_data = json::parse(test_case["input"].get<string>());
        auto result = solve(input_data);
        cout << result << endl;
    } catch (const exception& e) {
        cerr << e.what() << endl;
        return 1;
    }
    
    return 0;
}
`

const jsExecutableTemplate = `{{.Import}}

{{.UserCode}}

function runSingleTest() {
    const testCase = {{.TestData}};
    
    try {
        const inputData = JSON.parse(testCase.input);
        const result = main(inputData);
        
        if (typeof result === 'string') {
            console.log(result);
        } else {
            console.log(JSON.stringify(result));
        }
    } catch (error) {
        console.error(error.message);
        process.exit(1);
    }
}

runSingleTest();
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

	langType, err := langtype.StringToLanguageType(language)
	if err != nil {
		return "", fmt.Errorf("不支援的語言: %s", language)
	}

	switch langType {
	case langtype.LanguageTypePython:
		template = pythonExecutableTemplate
		imports = generateImportsForPython(inputType, outputType)
	case langtype.LanguageTypeCPP:
		template = cppExecutableTemplate
		imports = generateIncludesForCpp(inputType, outputType)
	case langtype.LanguageTypeJavaScript:
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

	langType, err := langtype.StringToLanguageType(language)
	if err != nil {
		return "", fmt.Errorf("不支援的語言: %s", language)
	}

	switch langType {
	case langtype.LanguageTypePython:
		template = pythonExecutableTemplate
		imports = generateImportsForPython(inputType, outputType)
	case langtype.LanguageTypeCPP:
		template = cppExecutableTemplate
		imports = generateIncludesForCpp(inputType, outputType)
	case langtype.LanguageTypeJavaScript:
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
		for _, imp := range inputType.GetRequiredImports(string(langtype.LanguageTypePython)) {
			imports[imp] = true
		}
	}

	if outputType != nil {
		for _, imp := range outputType.GetRequiredImports(string(langtype.LanguageTypePython)) {
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
		for _, inc := range inputType.GetRequiredImports(string(langtype.LanguageTypeCPP)) {
			includes[fmt.Sprintf("#include %s", inc)] = true
		}
	}

	if outputType != nil {
		for _, inc := range outputType.GetRequiredImports(string(langtype.LanguageTypeCPP)) {
			includes[fmt.Sprintf("#include %s", inc)] = true
		}
	}

	var includeLines []string
	for inc := range includes {
		includeLines = append(includeLines, inc)
	}

	return strings.Join(includeLines, "\n")
}
