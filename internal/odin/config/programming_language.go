package config

var Languages map[string]map[string]string = map[string]map[string]string{
	"python-3.10": {
		"nixPackageName": "python310",
		"version":        "3.10",
		"extension":      "py",
		"template":       "python.tmpl",
	},
	"rust": {
		"nixPackageName": "rust",
		"version":        "1.64.0",
		"extension":      "rs",
		"template":       "rust.tmpl",
	},
	"go-1.19": {
		"nixPackageName": "go_1_19",
		"version":        "1.19.1",
		"extension":      "go",
		"template":       "go.tmpl",
	},
	"java": {
		"nixPackageName": "jdk",
		"version":        "17",
		"extension":      "java",
		"template":       "java.tmpl",
	},
	"javascript": {
		"nixPackageName": "nodejs",
		"version":        "18.12.1",
		"extension":      "js",
		"template":       "javascript.tmpl",
	},
	"typescript": {
		"nixPackageName": "nodejs",
		"version":        "18.12.1",
		"extension":      "ts",
		"template":       "typescript.tmpl",
	},
}
