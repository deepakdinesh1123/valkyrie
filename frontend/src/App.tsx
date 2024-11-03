import React, { useCallback, useEffect, useState } from "react";
import LanguageSelector from "@/components/LanguageSelector";
import '@/App.css';
import { Input } from "@/components/ui/input";
import CodeEditor from "@/components/CodeEditor";
import { useLanguages } from "@/hooks/useLanguages";
import { useLanguageVersions } from "@/hooks/useLanguageVersions";
import { LanguageResponse, LanguageVersion } from "@/api-client";
import { Button } from "@/components/ui/button";
import { useCodeExecution } from "@/hooks/useCodeExecution";

const App: React.FC = () => {
  const [args, setArgs] = useState<string>("");
  const { selectedLanguage, setSelectedLanguage} = useLanguages();
  const { selectedLanguageVersion, setSelectedLanguageVersion} = useLanguageVersions(selectedLanguage?.id);
  const [codeContent, setCodeContent] = useState<string>("");
  const { terminalOutput, executeCode, isLoading } = useCodeExecution();
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] = useState<string[]>([]);
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<string[]>([]);

  const handleEditorChange = useCallback((content: string) => {
    setCodeContent(content);
}, []);

  const handleLanguageChange = useCallback((language: LanguageResponse) => {
    setSelectedLanguage(language);
  }, []);

  const handleVersionChange = useCallback((version: LanguageVersion) => {
    setSelectedLanguageVersion(version);
  }, []);

  const handleRunCode = useCallback(() => {
    if (selectedLanguage) {
      executeCode({
        language: selectedLanguage.name,
        code: codeContent,
        environment: {
          systemDependencies: selectedSystemDependencies,
          languageDependencies: selectedLanguageDependencies,
          args: args
        }
      });
    }
  }, [selectedLanguage, codeContent, args, executeCode]);

return (
    <div className="editor-container flex-1 w-full">
        <div className="top-bar flex flex-wrap justify-between items-center p-2 bg-transparent mr-20">
            <div className="flex flex-wrap items-center w-full sm:w-auto mb-2 sm:mb-0">
                <div className="w-full sm:w-auto mb-2 sm:mb-0 sm:mr-2 border-none">
                    <LanguageSelector
                      onLanguageChange={handleLanguageChange}
                      selectedLanguage={selectedLanguage}
                      onVersionChange={handleVersionChange}
                      selectedLanguageVersion={selectedLanguageVersion}
                    />
                </div>
                <Input
                    type="text"
                    placeholder="Args"
                    value={args}
                    onChange={(e) => setArgs(e.target.value)}
                    className="args-input w-full sm:w-36 mr-1 bg-neutral-900 text-white border-opacity-100 focus:ring-0"
                />
            </div>
            <div className="flex items-center w-full sm:w-auto justify-end">
            <Button
              onClick={handleRunCode}
              disabled={isLoading}
              className={`run-code-btn mr-2 ${isLoading ? 'loading' : ''} w-1/2 sm:w-auto bg-neutral-900 transition-colors hover:bg-stone-600 text-sm active:bg-neutral-900`}
            >
              {isLoading ? '' : 'Run Code'}
            </Button>
            {/* <Button
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              className="menu-toggle-btn w-1/2 sm:w-auto bg-neutral-900 transition-colors hover:bg-stone-600 text-sm active:bg-neutral-900"
            >
              {isSidebarOpen ? "Menu" : "Menu"}
            </Button> */}
          </div>
        </div>
        <div className="flex-grow overflow-hidden">
            <CodeEditor
                selectedLanguage={selectedLanguage}
                selectedLanguageVersion={selectedLanguageVersion}
                onChange={handleEditorChange}
                editorOptions={{ wordWrap: "on" }}
            />
        </div>
    </div>
);
};

export default App;