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
        input_data = test_case["input"]
        args = []
        for item in input_data:
            for key, value in item.items():
                args.append(value)

        result = main(*args)
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
        auto input_data = test_case["input"];
        {{.ParameterExtraction}}
        auto result = solve({{.ParameterList}});
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
        const inputData = testCase.input;
        const args = [];
        for (const item of inputData) {
            for (const key in item) {
                args.push(item[key]);
            }
        }

        const result = main(...args);

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

	if len(testData.InputFormat) == 0 {
		return "", fmt.Errorf("至少需要一個輸入格式")
	}

	if testData.OutputFormat.Type == "" {
		return "", fmt.Errorf("需要輸出格式類型")
	}

	outputType, err := loadjudge.ParseType(testData.OutputFormat.Type)
	if err != nil {
		return "", fmt.Errorf("無法解析輸出類型: %w", err)
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
		imports = generateImportsForPython(testData.InputFormat, outputType, string(langtype.LanguageTypePython))
	case langtype.LanguageTypeCPP:
		template = cppExecutableTemplate
		imports = generateIncludesForCpp(testData.InputFormat, outputType, string(langtype.LanguageTypeCPP))

		paramExtraction, paramList := generateCppParameterHandling(testData.InputFormat)
		template = strings.ReplaceAll(template, "{{.ParameterExtraction}}", paramExtraction)
		template = strings.ReplaceAll(template, "{{.ParameterList}}", paramList)
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

	if len(testData.TestConfig.InputFormat) == 0 {
		return "", fmt.Errorf("至少需要一個輸入格式")
	}

	if testData.TestConfig.OutputFormat.Type == "" {
		return "", fmt.Errorf("需要輸出格式類型")
	}

	outputType, err := loadjudge.ParseType(testData.TestConfig.OutputFormat.Type)
	if err != nil {
		return "", fmt.Errorf("無法解析輸出類型: %w", err)
	}

	var template string
	var imports string

	langType, err := langtype.StringToLanguageType(language)
	if err != nil {
		return "", fmt.Errorf("不支援的語言: %s", language)
	}

	var inputFormats []struct {
		Name string
		Type string
	}
	for _, inputFormat := range testData.TestConfig.InputFormat {
		inputFormats = append(inputFormats, struct {
			Name string
			Type string
		}{
			Name: inputFormat.Name,
			Type: inputFormat.Type,
		})
	}

	switch langType {
	case langtype.LanguageTypePython:
		template = pythonExecutableTemplate
		imports = generateImportsForPython(inputFormats, outputType, string(langtype.LanguageTypePython))
	case langtype.LanguageTypeCPP:
		template = cppExecutableTemplate
		imports = generateIncludesForCpp(inputFormats, outputType, string(langtype.LanguageTypeCPP))

		paramExtraction, paramList := generateCppParameterHandling(inputFormats)
		template = strings.ReplaceAll(template, "{{.ParameterExtraction}}", paramExtraction)
		template = strings.ReplaceAll(template, "{{.ParameterList}}", paramList)
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

func generateImportsForPython(inputFormats interface{}, outputType *loadjudge.TypeInfo, language string) string {
	imports := make(map[string]bool)

	var processFormat = func(name, typeStr string) {
		inputType, err := loadjudge.ParseType(typeStr)
		if err != nil {
			return
		}

		for _, imp := range inputType.GetRequiredImports(language) {
			imports[imp] = true
		}
	}

	switch v := inputFormats.(type) {
	case []struct {
		Name string
		Type string
	}:
		for _, inputFormat := range v {
			processFormat(inputFormat.Name, inputFormat.Type)
		}
	case []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}:
		for _, inputFormat := range v {
			processFormat(inputFormat.Name, inputFormat.Type)
		}
	}

	// 處理輸出類型的 imports
	if outputType != nil {
		for _, imp := range outputType.GetRequiredImports(language) {
			imports[imp] = true
		}
	}

	var importLines []string
	for imp := range imports {
		importLines = append(importLines, imp)
	}

	return strings.Join(importLines, "\n")
}

func generateIncludesForCpp(inputFormats interface{}, outputType *loadjudge.TypeInfo, language string) string {
	includes := map[string]bool{
		"#include <iostream>": true,
		"#include <vector>":   true,
		"#include <string>":   true,
		"#include <chrono>":   true,
	}

	var processFormat = func(name, typeStr string) {
		inputType, err := loadjudge.ParseType(typeStr)
		if err != nil {
			return
		}

		for _, inc := range inputType.GetRequiredImports(language) {
			includes[fmt.Sprintf("#include %s", inc)] = true
		}
	}

	switch v := inputFormats.(type) {
	case []struct {
		Name string
		Type string
	}:
		for _, inputFormat := range v {
			processFormat(inputFormat.Name, inputFormat.Type)
		}
	case []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}:
		for _, inputFormat := range v {
			processFormat(inputFormat.Name, inputFormat.Type)
		}
	}

	if outputType != nil {
		for _, inc := range outputType.GetRequiredImports(language) {
			includes[fmt.Sprintf("#include %s", inc)] = true
		}
	}

	var includeLines []string
	for inc := range includes {
		includeLines = append(includeLines, inc)
	}

	return strings.Join(includeLines, "\n")
}

func generateCppParameterHandling(inputFormats interface{}) (string, string) {
	var extractionLines []string
	var parameterNames []string

	var processFormat = func(name, typeStr string, index int) {
		inputType, err := loadjudge.ParseType(typeStr)
		if err != nil {
			return
		}

		cppType, err := inputType.GenerateLanguageType(string(langtype.LanguageTypeCPP))
		if err != nil {
			return
		}

		extractionLines = append(extractionLines, fmt.Sprintf("        auto %s = input_data[%d][\"%s\"].get<%s>();", name, index, name, cppType))
		parameterNames = append(parameterNames, name)
	}

	index := 0
	switch v := inputFormats.(type) {
	case []struct {
		Name string
		Type string
	}:
		for _, inputFormat := range v {
			processFormat(inputFormat.Name, inputFormat.Type, index)
			index++
		}
	case []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}:
		for _, inputFormat := range v {
			processFormat(inputFormat.Name, inputFormat.Type, index)
			index++
		}
	}

	extraction := strings.Join(extractionLines, "\n")
	paramList := strings.Join(parameterNames, ", ")

	return extraction, paramList
}
