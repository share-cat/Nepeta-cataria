package loadjudge

import (
	"fmt"
	"regexp"
	"strings"

	langtype "share-cat/Nepeta-cataria/pkg/language"
)

type TypeInfo struct {
	BaseType    string   // 基本類型: int, str, float, double
	Dimensions  []string // 維度信息: ["list", "list"] 表示 list[list[...]]
	IsContainer bool     // 是否為容器類型
}

type LanguageTypeMapping struct {
	TypeMap map[langtype.LanguageType]string
	Import  string
	Include string
}

// 輔助函數來獲取語言類型映射
func (ltm LanguageTypeMapping) GetType(lang langtype.LanguageType) string {
	if typeStr, exists := ltm.TypeMap[lang]; exists {
		return typeStr
	}
	return ""
}

var baseTypeMap = map[string]LanguageTypeMapping{
	"int": {
		TypeMap: map[langtype.LanguageType]string{
			langtype.LanguageTypePython:     "int",
			langtype.LanguageTypeCPP:        "int",
			langtype.LanguageTypeJavaScript: "number",
		},
	},
	"str": {
		TypeMap: map[langtype.LanguageType]string{
			langtype.LanguageTypePython:     "str",
			langtype.LanguageTypeCPP:        "std::string",
			langtype.LanguageTypeJavaScript: "string",
		},
		Include: "<string>",
	},
	"float": {
		TypeMap: map[langtype.LanguageType]string{
			langtype.LanguageTypePython:     "float",
			langtype.LanguageTypeCPP:        "float",
			langtype.LanguageTypeJavaScript: "number",
		},
	},
	"double": {
		TypeMap: map[langtype.LanguageType]string{
			langtype.LanguageTypePython:     "float",
			langtype.LanguageTypeCPP:        "double",
			langtype.LanguageTypeJavaScript: "number",
		},
	},
	"bool": {
		TypeMap: map[langtype.LanguageType]string{
			langtype.LanguageTypePython:     "bool",
			langtype.LanguageTypeCPP:        "bool",
			langtype.LanguageTypeJavaScript: "boolean",
		},
	},
}

var containerTypeMap = map[string]map[string]string{
	"list": {
		"python_template":     "List[%s]",
		"cpp_template":        "std::vector<%s>",
		"javascript_template": "%s[]",
		"python_import":       "from typing import List",
		"cpp_include":         "<vector>",
	},
	"set": {
		"python_template":     "Set[%s]",
		"cpp_template":        "std::set<%s>",
		"javascript_template": "Set<%s>",
		"python_import":       "from typing import Set",
		"cpp_include":         "<set>",
	},
	"dict": {
		"python_template":     "Dict[str, %s]",
		"cpp_template":        "std::map<std::string, %s>",
		"javascript_template": "Record<string, %s>",
		"python_import":       "from typing import Dict",
		"cpp_include":         "<map>",
	},
}

func ParseType(typeStr string) (*TypeInfo, error) {
	typeStr = strings.TrimSpace(typeStr)

	containerRegex := regexp.MustCompile(`^(list|set|dict)\[(.+)\]$`)

	info := &TypeInfo{}
	current := typeStr

	for {
		matches := containerRegex.FindStringSubmatch(current)
		if len(matches) != 3 {

			break
		}

		containerType := matches[1]
		innerType := matches[2]

		info.Dimensions = append(info.Dimensions, containerType)
		info.IsContainer = true
		current = innerType
	}

	if _, exists := baseTypeMap[current]; !exists {
		return nil, fmt.Errorf("unsupported base type: %s", current)
	}

	info.BaseType = current
	return info, nil
}

func (ti *TypeInfo) GenerateLanguageType(language string) (string, error) {
	baseMapping, exists := baseTypeMap[ti.BaseType]
	if !exists {
		return "", fmt.Errorf("unsupported base type: %s", ti.BaseType)
	}

	langType := stringToLanguageType(language)
	currentType := baseMapping.GetType(langType)
	if currentType == "" {
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	for i := len(ti.Dimensions) - 1; i >= 0; i-- {
		container := ti.Dimensions[i]
		containerMap, exists := containerTypeMap[container]
		if !exists {
			return "", fmt.Errorf("unsupported container type: %s", container)
		}

		var template string
		switch language {
		case string(langtype.LanguageTypePython):
			template = containerMap["python_template"]
		case string(langtype.LanguageTypeCPP):
			template = containerMap["cpp_template"]
		case string(langtype.LanguageTypeJavaScript):
			template = containerMap["javascript_template"]
		}

		currentType = fmt.Sprintf(template, currentType)
	}

	return currentType, nil
}

func (ti *TypeInfo) GetRequiredImports(language string) []string {
	var imports []string
	seen := make(map[string]bool)

	if baseMapping, exists := baseTypeMap[ti.BaseType]; exists {
		switch language {
		case string(langtype.LanguageTypePython):
			if baseMapping.Import != "" && !seen[baseMapping.Import] {
				imports = append(imports, baseMapping.Import)
				seen[baseMapping.Import] = true
			}
		case string(langtype.LanguageTypeCPP):
			if baseMapping.Include != "" && !seen[baseMapping.Include] {
				imports = append(imports, baseMapping.Include)
				seen[baseMapping.Include] = true
			}
		}
	}

	for _, container := range ti.Dimensions {
		if containerMap, exists := containerTypeMap[container]; exists {
			var importStr string
			switch language {
			case string(langtype.LanguageTypePython):
				importStr = containerMap["python_import"]
			case string(langtype.LanguageTypeCPP):
				importStr = containerMap["cpp_include"]
			}

			if importStr != "" && !seen[importStr] {
				imports = append(imports, importStr)
				seen[importStr] = true
			}
		}
	}

	return imports
}

func (ti *TypeInfo) GenerateFunctionParameter(language, paramName string) (string, error) {
	langType, err := ti.GenerateLanguageType(language)
	if err != nil {
		return "", err
	}

	switch language {
	case string(langtype.LanguageTypePython):
		return fmt.Sprintf("%s: %s", paramName, langType), nil
	case string(langtype.LanguageTypeCPP):
		if ti.IsContainer {
			return fmt.Sprintf("const %s& %s", langType, paramName), nil
		} else {
			return fmt.Sprintf("%s %s", langType, paramName), nil
		}
	case string(langtype.LanguageTypeJavaScript):
		return paramName, nil
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}
}

func ParseTestDataTypesString(Inputstring string, Outputstring string) (inputType *TypeInfo, outputType *TypeInfo, err error) {
	inputType, err = ParseType(Inputstring)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse input type '%s': %w", Inputstring, err)
	}
	outputType, err = ParseType(Outputstring)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse output type '%s': %w", Outputstring, err)
	}

	return inputType, outputType, nil
}

func ParseTestDataTypes(testConfig TestDataConfig) (inputType *TypeInfo, outputType *TypeInfo, err error) {
	if testConfig.TestConfig.InputFormat.Type != "" {
		inputType, err = ParseType(testConfig.TestConfig.InputFormat.Type)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse input type '%s': %w", testConfig.TestConfig.InputFormat.Type, err)
		}
	}

	if testConfig.TestConfig.OutputFormat.Type != "" {
		outputType, err = ParseType(testConfig.TestConfig.OutputFormat.Type)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse output type '%s': %w", testConfig.TestConfig.OutputFormat.Type, err)
		}
	}

	return inputType, outputType, nil
}

// 將字符串轉換為語言類型
func stringToLanguageType(langStr string) langtype.LanguageType {
	langType, err := langtype.StringToLanguageType(langStr)
	if err != nil {
		return langtype.LanguageTypeUnknown
	}
	return langType
}

func Test_Type_Parser() {
	examples := []string{
		"int",
		"str",
		"list[int]",
		"list[str]",
		"list[list[int]]",
		"list[list[str]]",
		"set[int]",
		"dict[int]",
	}

	for _, typeStr := range examples {
		fmt.Printf("\n=== Type: %s ===\n", typeStr)

		typeInfo, err := ParseType(typeStr)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", typeStr, err)
			continue
		}

		for _, lang := range []string{string(langtype.LanguageTypePython), string(langtype.LanguageTypeCPP), string(langtype.LanguageTypeJavaScript)} {
			langType, err := typeInfo.GenerateLanguageType(lang)
			if err != nil {
				fmt.Printf("%s: Error - %v\n", lang, err)
			} else {
				fmt.Printf("%s: %s\n", lang, langType)
			}

			param, err := typeInfo.GenerateFunctionParameter(lang, "data")
			if err != nil {
				fmt.Printf("%s param: Error - %v\n", lang, err)
			} else {
				fmt.Printf("%s param: %s\n", lang, param)
			}
		}

		pythonImports := typeInfo.GetRequiredImports(string(langtype.LanguageTypePython))
		if len(pythonImports) > 0 {
			fmt.Printf("Python imports: %v\n", pythonImports)
		}

		cppIncludes := typeInfo.GetRequiredImports(string(langtype.LanguageTypeCPP))
		if len(cppIncludes) > 0 {
			fmt.Printf("C++ includes: %v\n", cppIncludes)
		}
	}
}
