package config

var Languages map[string]map[string]string = map[string]map[string]string{
	"python-3.10": {
		"nixPackageName": "python310",
		"version":        "3.10",
		"extension":      "py",
	},
	"rust": {
		"nixPackageName": "rust",
		"version":        "1.64.0",
		"extension":      "rs",
	},
	"go-1.19": {
		"nixPackageName": "go_1_19",
		"version":        "1.19.1",
		"extension":      "go",
	},
	"java": {
		"nixPackageName": "jdk",
		"version":        "17",
		"extension":      "java",
	},
	"javascript": {
		"nixPackageName": "nodejs",
		"version":        "18.12.1",
		"extension":      "js",
	},
	"typescript": {
		"nixPackageName": "nodejs",
		"version":        "18.12.1",
		"extension":      "ts",
	},
}
