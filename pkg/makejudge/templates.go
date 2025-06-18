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

func fillTemplateDataWithTypes(td *TemplateData, language string, inputType, outputType *loadjudge.TypeInfo) error {
	td.UserCodePlaceholder = "Replace this with your solution."

	allImports := make(map[string]bool)

	if inputType != nil {
		imports := inputType.GetRequiredImports(language)
		for _, imp := range imports {
			allImports[imp] = true
		}
	}

	if outputType != nil {
		imports := outputType.GetRequiredImports(language)
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
		err := generatePythonTemplate(td, inputType, outputType)
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

		err := generateCppTemplate(td, inputType, outputType)
		if err != nil {
			return err
		}

	case string(langtype.LanguageTypeJavaScript):
		err := generateJavaScriptTemplate(td, inputType, outputType)
		if err != nil {
			return err
		}
	}

	return nil
}

func generatePythonTemplate(td *TemplateData, inputType, outputType *loadjudge.TypeInfo) error {

	if inputType != nil {
		inputParam, err := inputType.GenerateFunctionParameter(string(langtype.LanguageTypePython), "input_data")
		if err != nil {
			return fmt.Errorf("failed to generate Python input parameter: %w", err)
		}

		var outputTypeStr string
		if outputType != nil {
			outputTypeStr, err = outputType.GenerateLanguageType(string(langtype.LanguageTypePython))
			if err != nil {
				return fmt.Errorf("failed to generate Python output type: %w", err)
			}
		} else {
			outputTypeStr = "str"
		}

		td.FunctionSignature = fmt.Sprintf("def main(%s) -> %s:", inputParam, outputTypeStr)
	} else {

		td.FunctionSignature = "def main(input_lines: List[str]) -> str:"
		if td.RequiredImports == "" {
			td.RequiredImports = "from typing import List"
		} else if !strings.Contains(td.RequiredImports, "from typing import List") {
			td.RequiredImports += "\nfrom typing import List"
		}
	}
	return nil
}

func generateCppTemplate(td *TemplateData, inputType, outputType *loadjudge.TypeInfo) error {

	if inputType != nil {
		inputParam, err := inputType.GenerateFunctionParameter(string(langtype.LanguageTypeCPP), "input_data")
		if err != nil {
			return fmt.Errorf("failed to generate C++ input parameter: %w", err)
		}

		var outputTypeStr string
		if outputType != nil {
			outputTypeStr, err = outputType.GenerateLanguageType(string(langtype.LanguageTypeCPP))
			if err != nil {
				return fmt.Errorf("failed to generate C++ output type: %w", err)
			}
		} else {
			outputTypeStr = "std::string"
		}

		td.FunctionSignature = fmt.Sprintf("%s solve(%s) {", outputTypeStr, inputParam)
	} else {

		td.FunctionSignature = "std::string solve(const std::vector<std::string>& input_lines) {"
	}
	return nil
}

func generateJavaScriptTemplate(td *TemplateData, inputType, outputType *loadjudge.TypeInfo) error {

	if inputType != nil {
		inputTypeStr, err := inputType.GenerateLanguageType(string(langtype.LanguageTypeJavaScript))
		if err != nil {
			return fmt.Errorf("failed to generate JavaScript input type: %w", err)
		}

		var outputTypeStr string
		if outputType != nil {
			outputTypeStr, err = outputType.GenerateLanguageType(string(langtype.LanguageTypeJavaScript))
			if err != nil {
				return fmt.Errorf("failed to generate JavaScript output type: %w", err)
			}
		} else {
			outputTypeStr = "string"
		}

		td.FunctionSignature = fmt.Sprintf(`/**
 * @param {%s} input_data
 * @returns {%s}
 */
function main(input_data) {`, inputTypeStr, outputTypeStr)
	} else {

		td.FunctionSignature = "function main(input_lines) {"
	}
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

func executeCodeTemplateWithTypes(language string, data TemplateData) (string, error) {
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
