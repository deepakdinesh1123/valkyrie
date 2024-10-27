package config

var Languages map[string]map[string]string = map[string]map[string]string{
	"python-3.10": {
		"nixPackageName": "python310",
		"version":        "3.10",
		"extension":      "py",
		"template":       "python/python.tmpl",
		"monacoLanguage": "python",
		"defaultCode":    "# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == '__main__':\n    main()",
		"searchquery":    "python310Packages",
	},
	"python-3.11": {
		"nixPackageName": "python311",
		"version":        "3.11",
		"extension":      "py",
		"template":       "python/python.tmpl",
		"monacoLanguage": "python",
		"defaultCode":    "# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == '__main__':\n    main()",
		"searchquery":    "python311Packages",
	},
	"python-3.12": {
		"nixPackageName": "python312",
		"version":        "3.10",
		"extension":      "py",
		"template":       "python/python.tmpl",
		"monacoLanguage": "python",
		"defaultCode":    "# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == '__main__':\n    main()",
		"searchquery":    "python312Packages",
	},
	// "rust": {
	// 	"nixPackageName": "rustc",
	// 	"version":        "1.77.2",
	// 	"extension":      "rs",
	// 	"template":       "rust/rust.tmpl",
	// 	"monacoLanguage": "rust",
	// 	"defaultCode":    "fn main() {\n    // Type your Rust code here\n}",
	// 	"searchquery":    "rust",
	// },
	"go-1.19": {
		"nixPackageName": "go_1_19",
		"version":        "1.19.1",
		"extension":      "go",
		"template":       "go/go.tmpl",
		"monacoLanguage": "go",
		"defaultCode":    "package main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Type your Go code here\n}",
		"searchquery":    "go",
	},
	"go-1.23": {
		"nixPackageName": "go_1_23",
		"version":        "1.23.1",
		"extension":      "go",
		"template":       "go/go.tmpl",
		"monacoLanguage": "go",
		"defaultCode":    "package main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Type your Go code here\n}",
		"searchquery":    "go",
	},
	"go-1.21": {
		"nixPackageName": "go_1_21",
		"version":        "1.21.13",
		"extension":      "go",
		"template":       "go/go.tmpl",
		"monacoLanguage": "go",
		"defaultCode":    "package main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Type your Go code here\n}",
		"searchquery":    "go",
	},
	"go-1.22": {
		"nixPackageName": "go_1_22",
		"version":        "1.22.6",
		"extension":      "go",
		"template":       "go/go.tmpl",
		"monacoLanguage": "go",
		"defaultCode":    "package main\n\nimport \"fmt\"\n\nfunc main() {\n\t// Type your Go code here\n}",
		"searchquery":    "go",
	},
	"ada-14": {
		"nixPackageName": "gnat14",
		"version":        "14.1.0",
		"extension":      "adb",
		"template":       "ada/ada.tmpl",
		"monacoLanguage": "ada",
		"defaultCode":    "with Ada.Text_IO; use Ada.Text_IO;\n\nprocedure Hello is\nbegin\n    Put_Line(\"Hello, World!\");\nend Hello;",
		"searchquery":    "ada",
	},
	"ada-13": {
		"nixPackageName": "gnat13",
		"version":        "13.2.0",
		"extension":      "adb",
		"template":       "ada/ada.tmpl",
		"monacoLanguage": "ada",
		"defaultCode":    "with Ada.Text_IO; use Ada.Text_IO;\n\nprocedure Hello is\nbegin\n    Put_Line(\"Hello, World!\");\nend Hello;",
		"searchquery":    "ada",
	},
	"assembly": {
		"nixPackageName": "nasm",
		"version":        "2.16.03",
		"extension":      "asm",
		"template":       "assembly/assembly.tmpl",
		"monacoLanguage": "assembly",
		"defaultCode":    "",
		"searchquery":    "assembly",
	},
	"bash": {
		"nixPackageName": "bash",
		"version":        "5.2p32",
		"extension":      "bash",
		"template":       "bash/bash.tmpl",
		"monacoLanguage": "bash",
		"defaultCode":    "echo Hello World",
		"searchquery":    "bash",
	},
	"bun": {
		"nixPackageName": "bun",
		"version":        "1.18",
		"extension":      "js",
		"template":       "bun/bun.tmpl",
		"monacoLanguage": "javascript",
		"defaultCode":    "console.log('hello world')",
		"searchquery":    "bun",
	},
	"cobol": {
		"nixPackageName": "gnu-cobol",
		"version":        "3.1.2",
		"extension":      "cob",
		"template":       "cobol/cobol.tmpl",
		"monacoLanguage": "cobol",
		"defaultCode":    "IDENTIFICATION DIVISION.\nPROGRAM-ID. HelloWorld.\nDATA DIVISION.\nWORKING-STORAGE SECTION.\nPROCEDURE DIVISION.\n    DISPLAY \"Hello, World!\".\n    STOP RUN.",
		"searchquery":    "cobol",
	},
	"crystal-1.11": {
		"nixPackageName": "crystal",
		"version":        "1.11.2",
		"extension":      "cr",
		"template":       "crystal/crystal.tmpl",
		"monacoLanguage": "crystal",
		"defaultCode":    "puts \"Hello, World!\"",
		"searchquery":    "crystal",
	},
	"crystal-1.2": {
		"nixPackageName": "crystal_1_2",
		"version":        "1.2.2",
		"extension":      "cr",
		"template":       "crystal/crystal.tmpl",
		"monacoLanguage": "crystal",
		"defaultCode":    "puts \"Hello, World!\"",
		"searchquery":    "crystal",
	},
	"crystal-1.9": {
		"nixPackageName": "crystal_1_9",
		"version":        "1.9.2",
		"extension":      "cr",
		"template":       "crystal/crystal.tmpl",
		"monacoLanguage": "crystal",
		"defaultCode":    "puts \"Hello, World!\"",
		"searchquery":    "crystal",
	},
	"crystal-1.8": {
		"nixPackageName": "crystal_1_8",
		"version":        "1.8.2",
		"extension":      "cr",
		"template":       "crystal/crystal.tmpl",
		"monacoLanguage": "crystal",
		"defaultCode":    "puts \"Hello, World!\"",
		"searchquery":    "crystal",
	},
	"crystal-1.7": {
		"nixPackageName": "crystal_1_7",
		"version":        "1.7.3",
		"extension":      "cr",
		"template":       "crystal/crystal.tmpl",
		"monacoLanguage": "crystal",
		"defaultCode":    "puts \"Hello, World!\"",
		"searchquery":    "crystal",
	},
	"dart": {
		"nixPackageName": "dart",
		"version":        "3.3.4",
		"extension":      "dart",
		"template":       "dart/dart.tmpl",
		"monacoLanguage": "dart",
		"defaultCode":    "void main() {\n  print('Hello, World!');\n}",
		"searchquery":    "dart",
	},
	"deno": {
		"nixPackageName": "deno",
		"version":        "1.44.3",
		"extension":      "ts",
		"template":       "deno/deno.tmpl",
		"monacoLanguage": "typescript",
		"defaultCode":    "console.log('Hello World');",
		"searchquery":    "deno",
	},
	"fortran-13": {
		"nixPackageName": "gfortran",
		"version":        "13.2.0",
		"extension":      "f90",
		"template":       "fortran/fortran.tmpl",
		"monacoLanguage": "fortran",
		"defaultCode":    "program hello\n    print *, \"Hello, World!\"\nend program hello",
		"searchquery":    "fortran",
	},
	"fortran-12": {
		"nixPackageName": "gfortran12",
		"version":        "12.3.0",
		"extension":      "f90",
		"template":       "fortran/fortran.tmpl",
		"monacoLanguage": "fortran",
		"defaultCode":    "program hello\n    print *, \"Hello, World!\"\nend program hello",
		"searchquery":    "fortran",
	},
	"groovy": {
		"nixPackageName": "groovy",
		"version":        "3.0.11",
		"extension":      "groovy",
		"template":       "groovy/groovy.tmpl",
		"monacoLanguage": "groovy",
		"defaultCode":    "println 'Hello, World!'",
		"searchquery":    "groovy",
	},
	"julia": {
		"nixPackageName": "julia",
		"version":        "1.10.3",
		"extension":      "jl",
		"template":       "julia/julia.tmpl",
		"monacoLanguage": "julia",
		"defaultCode":    "println(\"Hello, World!\")",
		"searchquery":    "julia",
	},
	"lua": {
		"nixPackageName": "lua",
		"version":        "5.2.4",
		"extension":      "lua",
		"template":       "lua/lua.tmpl",
		"monacoLanguage": "lua",
		"defaultCode":    "print(\"Hello, World!\")",
		"searchquery":    "lua",
	},
	"nim": {
		"nixPackageName": "nim",
		"version":        "2.0.4",
		"extension":      "nim",
		"template":       "nim/nim.tmpl",
		"monacoLanguage": "nim",
		"defaultCode":    "echo \"Hello, World!\"",
		"searchquery":    "nim",
	},
	"node-22": {
		"nixPackageName": "nodejs_22",
		"version":        "22.4.1",
		"extension":      "js",
		"template":       "node/node.tmpl",
		"monacoLanguage": "javascript",
		"defaultCode":    "// Type your JavaScript code here\n\nfunction main() {\n    // Your code here\n}\n\nmain();",
		"searchquery":    "nodePackages",
	},
	"node-20": {
		"nixPackageName": "nodejs_20",
		"version":        "20.15.1",
		"extension":      "js",
		"template":       "node/node.tmpl",
		"monacoLanguage": "javascript",
		"defaultCode":    "// Type your JavaScript code here\n\nfunction main() {\n    // Your code here\n}\n\nmain();",
		"searchquery":    "nodePackages",
	},
	"node-18": {
		"nixPackageName": "nodejs_18",
		"version":        "18.20.4",
		"extension":      "js",
		"template":       "node/node.tmpl",
		"monacoLanguage": "javascript",
		"defaultCode":    "// Type your JavaScript code here\n\nfunction main() {\n    // Your code here\n}\n\nmain();",
		"searchquery":    "nodePackages",
	},
	"perl": {
		"nixPackageName": "perl",
		"version":        "5.38.2",
		"extension":      "pl",
		"template":       "perl/perl.tmpl",
		"monacoLanguage": "perl",
		"defaultCode":    "#!/usr/bin/perl\nuse strict;\nuse warnings;\n\nprint \"Hello, World!\\n\";\n",
		"searchquery":    "perl",
	},
	"php-8.2": {
		"nixPackageName": "php",
		"version":        "8.2.24",
		"extension":      "php",
		"template":       "php/php.tmpl",
		"monacoLanguage": "php",
		"defaultCode":    "<?php\necho \"Hello, World!\";\n?>\n",
		"searchquery":    "php",
	},
	"php-8.3": {
		"nixPackageName": "php",
		"version":        "8.3.12",
		"extension":      "php",
		"template":       "php/php.tmpl",
		"monacoLanguage": "php",
		"defaultCode":    "<?php\necho \"Hello, World!\";\n?>\n",
		"searchquery":    "php",
	},
	"php-8.1": {
		"nixPackageName": "php",
		"version":        "8.1.30",
		"extension":      "php",
		"template":       "php/php.tmpl",
		"monacoLanguage": "php",
		"defaultCode":    "<?php\necho \"Hello, World!\";\n?>\n",
		"searchquery":    "php",
	},
	"sql": {
		"nixPackageName": "sql",
		"version":        "3.45.3",
		"extension":      "sql",
		"template":       "sql/sql.tmpl",
		"monacoLanguage": "sql",
		"defaultCode":    "CREATE TABLE employees (id INT PRIMARY KEY, name VARCHAR(100), salary DECIMAL(10, 2));",
		"searchquery":    "sql",
	},
	"swift": {
		"nixPackageName": "swift",
		"version":        "5.8",
		"extension":      "swift",
		"template":       "swift/swift.tmpl",
		"monacoLanguage": "swift",
		"defaultCode":    "print(\"Hello, World!\")",
		"searchquery":    "swift",
	},
	"zig": {
		"nixPackageName": "zig",
		"version":        "3.45.3",
		"extension":      "zig",
		"template":       "zig/zig.tmpl",
		"monacoLanguage": "zig",
		"defaultCode":    "const std = @import(\"std\");\n\npub fn main() !void {\n    const stdout = std.io.getStdOut().writer();\n    try stdout.print(\"Hello, World!\\n\", .{});\n}",
		"searchquery":    "zig",
	},
}
