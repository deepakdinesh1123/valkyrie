package config

var Languages map[string]map[string]string = map[string]map[string]string{
	"python-3.10": {
		"nixPackageName":  "python310",
		"version":         "3.10",
		"extension":       "py",
		"flake_template":  "python/python.flake.tmpl",
		"script_template": "python/python.script.tmpl",
		"monacoLanguage":  "python",
		"defaultCode":     "# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == '__main__':\n    main()",
		"searchquery":     "python310Packages",
	},
	"python-3.11": {
		"nixPackageName":  "python311",
		"version":         "3.11",
		"extension":       "py",
		"flake_template":  "python/python.flake.tmpl",
		"script_template": "python/python.script.tmpl",
		"monacoLanguage":  "python",
		"defaultCode":     "# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == '__main__':\n    main()",
		"searchquery":     "python311Packages",
	},
	"python-3.12": {
		"nixPackageName":  "python312",
		"version":         "3.10",
		"extension":       "py",
		"flake_template":  "python/python.flake.tmpl",
		"script_template": "python/python.script.tmpl",
		"monacoLanguage":  "python",
		"defaultCode":     "# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == '__main__':\n    main()",
		"searchquery":     "python312Packages",
	},
	"rust": {
		"nixPackageName": "rust",
		"version":        "1.64.0",
		"extension":      "rs",
		"template":       "rust.tmpl",
		"monacoLanguage": "rust",
		"defaultCode":    "fn main() {\n    // Type your Rust code here\n}",
		"searchquery":    "rust",
	},
	"go-1.19": {
		"nixPackageName": "go_1_19",
		"version":        "1.19.1",
		"extension":      "go",
		"template":       "go.tmpl",
		"monacoLanguage": "go",
		"defaultCode":    "package main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Type your Go code here\n}",
		"searchquery":    "go",
	},
	"java": {
		"nixPackageName": "jdk",
		"version":        "17",
		"extension":      "java",
		"template":       "java.tmpl",
		"monacoLanguage": "java",
		"defaultCode":    "public class Main {\n    public static void main(String[] args) {\n        // Type your Java code here\n    }\n}",
		"searchquery":    "javaPackages",
	},
	"javascript": {
		"nixPackageName": "nodejs",
		"version":        "18.12.1",
		"extension":      "js",
		"template":       "javascript.tmpl",
		"monacoLanguage": "javascript",
		"defaultCode":    "// Type your JavaScript code here\n\nfunction main() {\n    // Your code here\n}\n\nmain();",
		"searchquery":    "nodePackages",
	},
	"typescript": {
		"nixPackageName": "nodejs",
		"version":        "18.12.1",
		"extension":      "ts",
		"template":       "typescript.tmpl",
		"monacoLanguage": "typescript",
		"defaultCode":    "// Type your TypeScript code here\n\nfunction main() {\n    // Your code here\n}\n\nmain();",
		"searchquery":    "nodePackages",
	},
}
