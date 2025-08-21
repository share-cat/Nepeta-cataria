package makejudge

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	langtype "share-cat/Nepeta-cataria/pkg/language"
	"share-cat/Nepeta-cataria/pkg/loadjudge"
)

var funcMap = template.FuncMap{
	"IndentLines": func(s string, indent string) string {
		lines := strings.Split(s, "\n")
		indentedLines := make([]string, len(lines))
		for i, line := range lines {
			indentedLines[i] = indent + line
		}
		return strings.Join(indentedLines, "\n")
	},
	"Join": strings.Join,
}

const pythonMainTemplateStr = `
# Problem: {{.ProblemName}}
# Sample Input:
# {{.SampleInput | IndentLines "  # "}}
# Sample Output:
# {{.SampleOutput | IndentLines "  # "}}

{{.FunctionSignature}}
    # {{.UserCodePlaceholder}}

    raise NotImplementedError("Please implement the solution in the 'main' function.")
`

// --- C++ Template ---
const cppMainTemplateStr = `
// Problem: {{.ProblemName}}
// Sample Input:
// {{.SampleInput | IndentLines "// "}}
// Sample Output:
// {{.SampleOutput | IndentLines "// "}}

{{.FunctionSignature}}
    // {{.UserCodePlaceholder}}

    throw std::runtime_error("Please implement the solution in the 'solve' function.");
}

`

// --- JavaScript Template ---
const javascriptMainTemplateStr = `
// Problem: {{.ProblemName}}
// Sample Input:
// {{.SampleInput | IndentLines "// "}}
// Sample Output:
// {{.SampleOutput | IndentLines "// "}}

{{.FunctionSignature}}
    // {{.UserCodePlaceholder}}

    throw new Error("Please implement the solution in the 'main' function.");
}

`

func fillTemplateDataWithTypes(td *TemplateData, language string, testData *loadjudge.TestDataConfig) error {
	if testData == nil {
		return fmt.Errorf("testData is required for template generation")
	}

	td.UserCodePlaceholder = "Replace this with your solution."

	allImports := make(map[string]bool)

	var outputType *loadjudge.TypeInfo
	if testData.TestConfig.OutputFormat.Type == "" {
		return fmt.Errorf("output format type is required in testdata.yaml")
	}

	var err error
	outputType, err = loadjudge.ParseType(testData.TestConfig.OutputFormat.Type)
	if err != nil {
		return fmt.Errorf("failed to parse output type: %w", err)
	}

	imports := outputType.GetRequiredImports(language)
	for _, imp := range imports {
		allImports[imp] = true
	}

	if len(testData.TestConfig.InputFormat) == 0 {
		return fmt.Errorf("at least one input format is required in testdata.yaml")
	}

	var inputFormats []struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	}

	for _, inputFormat := range testData.TestConfig.InputFormat {
		if inputFormat.Name == "" || inputFormat.Type == "" {
			return fmt.Errorf("input format must have both name and type")
		}

		inputFormats = append(inputFormats, struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
		}{
			Name: inputFormat.Name,
			Type: inputFormat.Type,
		})

		inputType, err := loadjudge.ParseType(inputFormat.Type)
		if err != nil {
			return fmt.Errorf("failed to parse input type '%s' for parameter '%s': %w", inputFormat.Type, inputFormat.Name, err)
		}

		imports := inputType.GetRequiredImports(language)
		for _, imp := range imports {
			allImports[imp] = true
		}
	}

	var importLines []string
	for imp := range allImports {
		importLines = append(importLines, imp)
	}

	switch language {
	case string(langtype.LanguageTypePython):
		td.RequiredImports = strings.Join(importLines, "\n")
		err := generatePythonTemplate(td, inputFormats, outputType)
		if err != nil {
			return err
		}

	case string(langtype.LanguageTypeCPP):
		var includeLines []string
		for _, inc := range importLines {
			includeLines = append(includeLines, fmt.Sprintf("#include %s", inc))
		}
		standardIncludes := []string{
			"#include <iostream>",
			"#include <sstream>",
			"#include <stdexcept>",
		}
		includeLines = append(standardIncludes, includeLines...)
		td.RequiredIncludes = strings.Join(includeLines, "\n")

		err := generateCppTemplate(td, inputFormats, outputType)
		if err != nil {
			return err
		}

	case string(langtype.LanguageTypeJavaScript):
		err := generateJavaScriptTemplate(td, inputFormats, outputType)
		if err != nil {
			return err
		}
	}

	return nil
}

func generatePythonTemplate(td *TemplateData, inputFormats []struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}, outputType *loadjudge.TypeInfo) error {
	if len(inputFormats) == 0 {
		return fmt.Errorf("input formats are required for Python template generation")
	}

	if outputType == nil {
		return fmt.Errorf("output type is required for Python template generation")
	}

	var parameters []string

	// 生成每個輸入參數
	for _, inputFormat := range inputFormats {
		inputType, err := loadjudge.ParseType(inputFormat.Type)
		if err != nil {
			return fmt.Errorf("failed to parse input type '%s' for parameter '%s': %w", inputFormat.Type, inputFormat.Name, err)
		}

		inputParam, err := inputType.GenerateFunctionParameter(string(langtype.LanguageTypePython), inputFormat.Name)
		if err != nil {
			return fmt.Errorf("failed to generate Python input parameter for '%s': %w", inputFormat.Name, err)
		}
		parameters = append(parameters, inputParam)
	}

	outputTypeStr, err := outputType.GenerateLanguageType(string(langtype.LanguageTypePython))
	if err != nil {
		return fmt.Errorf("failed to generate Python output type: %w", err)
	}

	td.FunctionSignature = fmt.Sprintf("def main(%s) -> %s:", strings.Join(parameters, ", "), outputTypeStr)
	return nil
}

func generateCppTemplate(td *TemplateData, inputFormats []struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}, outputType *loadjudge.TypeInfo) error {
	if len(inputFormats) == 0 {
		return fmt.Errorf("input formats are required for C++ template generation")
	}

	if outputType == nil {
		return fmt.Errorf("output type is required for C++ template generation")
	}

	var parameters []string

	for _, inputFormat := range inputFormats {
		inputType, err := loadjudge.ParseType(inputFormat.Type)
		if err != nil {
			return fmt.Errorf("failed to parse input type '%s' for parameter '%s': %w", inputFormat.Type, inputFormat.Name, err)
		}

		inputParam, err := inputType.GenerateFunctionParameter(string(langtype.LanguageTypeCPP), inputFormat.Name)
		if err != nil {
			return fmt.Errorf("failed to generate C++ input parameter for '%s': %w", inputFormat.Name, err)
		}
		parameters = append(parameters, inputParam)
	}

	outputTypeStr, err := outputType.GenerateLanguageType(string(langtype.LanguageTypeCPP))
	if err != nil {
		return fmt.Errorf("failed to generate C++ output type: %w", err)
	}

	td.FunctionSignature = fmt.Sprintf("%s solve(%s) {", outputTypeStr, strings.Join(parameters, ", "))
	return nil
}

func generateJavaScriptTemplate(td *TemplateData, inputFormats []struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}, outputType *loadjudge.TypeInfo) error {
	if len(inputFormats) == 0 {
		return fmt.Errorf("input formats are required for JavaScript template generation")
	}

	if outputType == nil {
		return fmt.Errorf("output type is required for JavaScript template generation")
	}

	var parameterNames []string
	var parameterComments []string

	for _, inputFormat := range inputFormats {
		parameterNames = append(parameterNames, inputFormat.Name)

		inputType, err := loadjudge.ParseType(inputFormat.Type)
		if err != nil {
			return fmt.Errorf("failed to parse input type '%s' for parameter '%s': %w", inputFormat.Type, inputFormat.Name, err)
		}

		inputTypeStr, err := inputType.GenerateLanguageType(string(langtype.LanguageTypeJavaScript))
		if err != nil {
			return fmt.Errorf("failed to generate JavaScript input type for '%s': %w", inputFormat.Name, err)
		}
		parameterComments = append(parameterComments, fmt.Sprintf(" * @param {%s} %s", inputTypeStr, inputFormat.Name))
	}

	outputTypeStr, err := outputType.GenerateLanguageType(string(langtype.LanguageTypeJavaScript))
	if err != nil {
		return fmt.Errorf("failed to generate JavaScript output type: %w", err)
	}

	td.FunctionSignature = fmt.Sprintf(`/**
%s
 * @returns {%s}
 */
function main(%s) {`, strings.Join(parameterComments, "\n"), outputTypeStr, strings.Join(parameterNames, ", "))
	return nil
}

func getTemplate(language string) (*template.Template, error) {
	var tplStr string
	templateName := "mainSolution"
	switch language {
	case string(langtype.LanguageTypePython):
		tplStr = pythonMainTemplateStr
	case string(langtype.LanguageTypeCPP):
		tplStr = cppMainTemplateStr
		templateName = "cppSolution"
	case string(langtype.LanguageTypeJavaScript):
		tplStr = javascriptMainTemplateStr
	default:
		return nil, fmt.Errorf("unsupported language for template: %s", language)
	}
	return template.New(templateName).Funcs(funcMap).Parse(tplStr)
}

func executeCodeTemplateWithTypes(language string, data TemplateData, testData *loadjudge.TestDataConfig) (string, error) {
	err := fillTemplateDataWithTypes(&data, language, testData)
	if err != nil {
		return "", fmt.Errorf("failed to fill template data: %w", err)
	}

	tmpl, err := getTemplate(language)
	if err != nil {
		return "", err
	}

	var tplBuffer bytes.Buffer
	if err := tmpl.Execute(&tplBuffer, data); err != nil {
		return "", fmt.Errorf("failed to execute %s template: %w", language, err)
	}
	return tplBuffer.String(), nil
}
