import React, { useState } from "react";
import CodeEditor from "@/components/CodeEditor";
import ListBuilder from "@/components/ListBuilder";
import Terminal from "@/components/Terminal";
import { Button } from "@/components/ui/button";
import "@/App.css";
import { useLanguages } from '@/hooks/useLanguages';
import { useCodeExecution } from '@/hooks/useCodeExecution';
import { LanguageSelector } from '@/components/LanguageSelector';

const systemDependencies = ["git", "postgres", "nodejs", "python", "docker", "nginx", "redis"];
const languageDependencies = ["git", "postgres", "nodejs", "python", "docker", "nginx", "redis"];

const App: React.FC = () => {
  const { languages, selectedLanguage, setSelectedLanguage } = useLanguages();
  const { terminalOutput, executeCode } = useCodeExecution();
  const [codeContent, setCodeContent] = useState<string>("");
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<string[]>([]);
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] = useState<string[]>([]);

  const handleEditorChange = (_tab: string, content: string) => {
    setCodeContent(content);
  };

  const handleRunCode = () => {
    if (selectedLanguage) {
      executeCode({
        language: selectedLanguage.name,
        code: codeContent,
        systemDependencies: selectedSystemDependencies,
        languageDependencies: selectedLanguageDependencies,
      });
    }
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
          <span className="terminal-label">Output</span>
          <div className="output-terminal">
            <Terminal output={terminalOutput} />
          </div>
        </div>

        <LanguageSelector
          languages={languages}
          selectedLanguage={selectedLanguage}
          onLanguageChange={(language) => {
            setSelectedLanguage(language);
            setCodeContent(language.defaultcode);
          }}
        />

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
            onSelectionChange={setSelectedSystemDependencies}
          />

          <div className="my-2"></div>
          <span>Language Dependencies</span>
          <ListBuilder
            items={languageDependencies}
            onSelectionChange={setSelectedLanguageDependencies}
          />
        </div>
      </div>
    </div>
  );
};

export default App;