import { useState } from "react";
import CodeEditor from "./components/CodeEditor";
import ListBuilder from "./components/ListBuilder";
import Terminal from "./components/Terminal";
import { Button } from "@/components/ui/button";
import "./App.css";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const App = () => {
  const [selectedLanguage, setSelectedLanguage] = useState({
    name: "Python",
    extension: "py",
    monacoLanguage: "python",
    defaultCode:
      '# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == "__main__":\n    main()',
  });
  const [codeContent, setCodeContent] = useState(selectedLanguage.defaultCode);
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<
    string[]
  >([]);
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] =
    useState<string[]>([]);
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]); // Added state for terminal output

  const systemDependencies = [
    "git",
    "postgres",
    "nodejs",
    "python",
    "docker",
    "nginx",
    "redis",
  ];
  const languageDependencies = [
    "git",
    "postgres",
    "nodejs",
    "python",
    "docker",
    "nginx",
    "redis",
  ];

  const languages = [
    {
      name: "Python",
      extension: "py",
      monacoLanguage: "python",
      defaultCode:
        '# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == "__main__":\n    main()',
    },
    {
      name: "JavaScript",
      extension: "js",
      monacoLanguage: "javascript",
      defaultCode:
        "// Type your JavaScript code here\n\nfunction main() {\n    // Your code here\n}\n\nmain();",
    },
    {
      name: "Rust",
      extension: "rs",
      monacoLanguage: "rust",
      defaultCode: "fn main() {\n    // Type your Rust code here\n}",
    },
    {
      name: "Go",
      extension: "go",
      monacoLanguage: "go",
      defaultCode:
        'package main\n\nimport "fmt"\n\nfunc main() {\n\t// Type your Go code here\n}',
    },
    {
      name: "C++",
      extension: "cpp",
      monacoLanguage: "cpp",
      defaultCode:
        "#include <iostream>\n\nint main() {\n    // Type your C++ code here\n    return 0;\n}",
    },
  ];

  const handleEditorChange = (tabName: string, content: string) => {
    setCodeContent(content);
  };

  const handleSystemSelectionChange = (dependencies: string[]) =>
    setSelectedSystemDependencies(dependencies);

  const handleLanguageSelectionChange = (dependencies: string[]) =>
    setSelectedLanguageDependencies(dependencies);

  const handleRunCode = () => {
    const runData = {
      language: selectedLanguage.monacoLanguage,
      code: codeContent,
      systemDependencies: selectedSystemDependencies,
      languageDependencies: selectedLanguageDependencies,
    };

    const output = `Running code: 
      System: ${runData.systemDependencies.join(", ") || "None"}\n
      Language: ${runData.languageDependencies.join(", ") || "None"}\n
      Code: ${runData.code}`;

    setTerminalOutput((prev) => [...prev, output]);
  };

  return (
    <div className="App">
      <div className="editor-container">
        <CodeEditor
          languages={languages}
          selectedLanguage={selectedLanguage}
          onChange={handleEditorChange}
          editorOptions={{ wordWrap: "on" }}
        />
        <div className="terminal-container">
          <span className="terminal-label ">Output</span>
          <div className="output-terminal">
            <Terminal output={terminalOutput} />
          </div>
        </div>

        <div className="language-selector">
          <Select
            value={selectedLanguage.name}
            onValueChange={(value) => {
              const language = languages.find((lang) => lang.name === value);
              if (language) {
                setSelectedLanguage(language);
                setCodeContent(language.defaultCode);
              }
            }}
          >
            <SelectTrigger className="w-[180px] bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 rounded-md px-2 py-1 transition-colors duration-200 ease-in-out">
              <SelectValue placeholder="Select a language" />
            </SelectTrigger>
            <SelectContent className="dark:bg-gray-800">
              {languages.map((lang) => (
                <SelectItem key={lang.name} value={lang.name}>
                  {lang.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="run-code-button">
          <Button onClick={handleRunCode}>Run Code</Button>
        </div>
      </div>

      <div className="sidebar">
        <div className="flex flex-col gap-4 h-full overflow-y-auto">
          <div className="my-1"></div>

          <span>System Dependencies</span>
          <ListBuilder
            items={systemDependencies}
            onSelectionChange={handleSystemSelectionChange}
          />

          <div className="my-2"></div>
          <span>Language Dependencies</span>
          <ListBuilder
            items={languageDependencies}
            onSelectionChange={handleLanguageSelectionChange}
          />
        </div>
      </div>
    </div>
  );
};

export default App;
