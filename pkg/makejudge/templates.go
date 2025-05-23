package makejudge

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
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

func getTemplate(language string) (*template.Template, error) {
	var tplStr string
	templateName := "mainSolution"
	switch language {
	case Python:
		tplStr = pythonMainTemplateStr
	case Cpp:
		tplStr = cppMainTemplateStr
		templateName = "cppSolution"
	case JavaScript:
		tplStr = javascriptMainTemplateStr
	default:
		return nil, fmt.Errorf("unsupported language for template: %s", language)
	}
	return template.New(templateName).Funcs(funcMap).Parse(tplStr)
}

func fillTemplateDataForJudge(td *TemplateData, language string, types []string) {
	td.UserCodePlaceholder = "Replace this with your solution."
	switch language {
	case Python:
		if len(types) > 0 {
			var params []string
			for i, typ := range types {
				switch typ {
				case "int":
					params = append(params, fmt.Sprintf("param%d: int", i+1))
				case "float":
					params = append(params, fmt.Sprintf("param%d: float", i+1))
				case "string", "str":
					params = append(params, fmt.Sprintf("param%d: str", i+1))
				case "list":
					params = append(params, fmt.Sprintf("param%d: List[int]", i+1))
				case "list_str":
					params = append(params, fmt.Sprintf("param%d: List[str]", i+1))
				default:
					params = append(params, fmt.Sprintf("param%d: Any", i+1))
				}
			}

			td.FunctionSignature = fmt.Sprintf("def main(%s) -> str:", strings.Join(params, ", "))
		} else {
			td.FunctionSignature = "def main(input_datas: List[str]) -> str:"
		}
		td.ExecutionBoilerplate = `
		if __name__ == "__main__":
		    import sys
		    input_lines = [line.rstrip('\n') for line in sys.stdin]
		    try:
		        output_str = main(input_lines)
		        if not isinstance(output_str, str):
		            output_str = str(output_str) # Ensure the final output is a string
		        sys.stdout.write(output_str)
		        # Judge server might handle final newline, or expect it from solution.
		        # For consistency, let's ensure one if output is not empty and doesn't have one.
		        if output_str and not output_str.endswith('\n'):
		            sys.stdout.write('\n')
		    except NotImplementedError as e:
		        sys.stderr.write(f"Solution not implemented: {e}\n")
		        sys.exit(1) # Non-zero exit code for error
		    except Exception as e:
		        sys.stderr.write(f"Runtime Error: {e}\n")
		        # import traceback
		        # traceback.print_exc(file=sys.stderr) # Useful for user debugging locally
		        sys.exit(1) # Non-zero exit code for error
		`
	case Cpp:
		if len(types) > 0 {
			var params []string
			for i, typ := range types {
				switch typ {
				case "int":
					params = append(params, fmt.Sprintf("int param%d", i+1))
				case "float":
					params = append(params, fmt.Sprintf("double param%d", i+1))
				case "string", "str":
					params = append(params, fmt.Sprintf("const std::string& param%d", i+1))
				case "list":
					params = append(params, fmt.Sprintf("const std::vector<int>& param%d", i+1))
				case "list_str":
					params = append(params, fmt.Sprintf("const std::vector<std::string>& param%d", i+1))
				default:
					params = append(params, fmt.Sprintf("const std::string& param%d", i+1))
				}
			}
			td.FunctionSignature = fmt.Sprintf("std::string solve(%s) {", strings.Join(params, ", "))
		} else {
			td.FunctionSignature = "std::string solve(const std::vector<std::string>& input_datas) {"
		}
		td.ExecutionBoilerplate = `
int main() {
    std::ios_base::sync_with_stdio(false);
    std::cin.tie(NULL);
    std::vector<std::string> input_lines;
    std::string line;
    while (std::getline(std::cin, line)) {
        input_lines.push_back(line);
    }
    try {
        std::string output_result_str = solve(input_lines);
        std::cout << output_result_str;
        // Judge server might handle final newline.
        // For consistency, if output is not empty and doesn't have one:
        if (!output_result_str.empty() && output_result_str.back() != '\n') {
            std::cout << std::endl;
        }
    } catch (const std::runtime_error& e) {
        std::cerr << "Runtime Error: " << e.what() << std::endl;
        return 1; // Non-zero exit code for error
    } catch (const std::exception& e) {
        std::cerr << "An unexpected error occurred: " << e.what() << std::endl;
        return 1;
    } catch (...) {
        std::cerr << "An unknown error occurred." << std::endl;
        return 1;
    }
    return 0;
}
`
	case JavaScript:
		if len(types) > 0 {
			var jsDocParams []string
			var params []string
			for i, typ := range types {
				paramName := fmt.Sprintf("param%d", i+1)
				switch typ {
				case "int", "float":
					jsDocParams = append(jsDocParams, fmt.Sprintf(" * @param {number} %s", paramName))
				case "string", "str":
					jsDocParams = append(jsDocParams, fmt.Sprintf(" * @param {string} %s", paramName))
				case "list":
					jsDocParams = append(jsDocParams, fmt.Sprintf(" * @param {number[]} %s", paramName))
				case "list_str":
					jsDocParams = append(jsDocParams, fmt.Sprintf(" * @param {string[]} %s", paramName))
				default:
					jsDocParams = append(jsDocParams, fmt.Sprintf(" * @param {any} %s", paramName))
				}
				params = append(params, paramName)
			}

			td.FunctionSignature = fmt.Sprintf("/**\n%s\n * @returns {string}\n */\nfunction main(%s) {",
				strings.Join(jsDocParams, "\n"), strings.Join(params, ", "))
		} else {
			td.FunctionSignature = "function main(input_datas) {"
		}
		td.ExecutionBoilerplate = `
function run() {
    const readline = require('readline');
    const rl = readline.createInterface({
        input: process.stdin,
        crlfDelay: Infinity
    });
    const input_lines = [];
    rl.on('line', (line) => {
        input_lines.push(line);
    });
    rl.on('close', () => {
        try {
            let output_str = main(input_lines);
            if (typeof output_str !== 'string') {
                output_str = String(output_str);
            }
            process.stdout.write(output_str);
            if (output_str && !output_str.endsWith('\n')) {
                 process.stdout.write('\\n');
            }
        } catch (e) {
            console.error(e.message);
            process.exit(1); // Non-zero exit code for error
        }
    });
}
run();
`
	}
}

func executeCodeTemplate(language string, data TemplateData) (string, error) {
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
