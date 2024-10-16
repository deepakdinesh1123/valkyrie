import React, { useCallback, useEffect, useState } from "react";
import CodeEditor from "@/components/CodeEditor";
import ListBuilder from "@/components/ListBuilder";
import Terminal from "@/components/Terminal";
import { Button } from "@/components/ui/button";
import "@/App.css";
import { useLanguages } from '@/hooks/useLanguages';
import { useCodeExecution } from '@/hooks/useCodeExecution';
import { useSystemPackages } from '@/hooks/useSystemPackages';
import { useLanguagePackages } from '@/hooks/useLanguagePackages';
import { LanguageSelector } from '@/components/LanguageSelector';
import { usePackagesExist } from "@/hooks/usePackageExists";

const App: React.FC = () => {
  const { languages, selectedLanguage, setSelectedLanguage } = useLanguages();
  const { terminalOutput, executeCode, isLoading } = useCodeExecution();

  const [codeContent, setCodeContent] = useState<string>("");
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] = useState<string[]>([]);
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<string[]>([]);
  const [systemSearchString, setSystemSearchString] = useState<string>("");
  const [languageSearchString, setLanguageSearchString] = useState<string>("");
  const [resetLanguageDependencies, setResetLanguageDependencies] = useState({});
  const { systemPackages, loading: loadingSystemPackages, error: systemPackagesError } = useSystemPackages(systemSearchString);
  const { languagePackages, loading: loadingLanguagePackages, error: languagePackagesError, resetLanguagePackages } = useLanguagePackages(languageSearchString, selectedLanguage?.searchquery);

  const [selectedLanguagePrefix, setSelectedLanguagePrefix] = useState<string>("");
  const [pendingLanguageChange, setPendingLanguageChange] = useState<any>(null);
  const { existsResponse } = usePackagesExist(
    pendingLanguageChange?.searchquery || "",
    selectedLanguageDependencies
  );

  useEffect(() => {
    if (pendingLanguageChange && existsResponse) {
      handleLanguageChangeEffect(pendingLanguageChange, existsResponse);
      setPendingLanguageChange(null);
    }
  }, [pendingLanguageChange, existsResponse]);

  const handleEditorChange = (content: string) => {
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

  const handleLanguageChangeEffect = useCallback((language: any, existsResponse: any) => {
    const newPrefix = language.name.split("-")[0];

    if (newPrefix !== selectedLanguagePrefix) {
      setSelectedLanguagePrefix(newPrefix);
      resetOnNewLanguage(language);
    } else {
      if (language.name !== selectedLanguage?.name) {
        if (existsResponse.exists) {
          setSelectedLanguage(language);
          setCodeContent(language.defaultcode);
        } else {
          const nonExistingPackages = existsResponse.nonExistingPackages;

          const updatedDependencies = selectedLanguageDependencies.filter(
            dep => !nonExistingPackages.includes(dep)
          );

          setSelectedLanguageDependencies(updatedDependencies);

          setSelectedLanguage(language);
          setCodeContent(language.defaultcode);
        }
      }
    }
  }, [selectedLanguagePrefix, selectedLanguage, selectedLanguageDependencies, setSelectedLanguage]);

  const handleLanguageChange = useCallback((language: any) => {
    setPendingLanguageChange(language);
  }, []);

  const resetOnNewLanguage = useCallback((language: any) => {
    setSelectedLanguage(language);
    setCodeContent(language.defaultcode);
    setLanguageSearchString("");
    setSelectedLanguageDependencies([]);
    setResetLanguageDependencies({});
    resetLanguagePackages();
  }, [setSelectedLanguage, resetLanguagePackages]);

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
          onLanguageChange={handleLanguageChange}
        />

        <div className="run-code-button">
          <Button
            onClick={handleRunCode}
            disabled={isLoading}
            className={`run-code-btn ${isLoading ? 'loading' : ''}`}
          >
            {!isLoading && 'Run'}
          </Button>
        </div>

        {loadingSystemPackages && <div>Loading system packages...</div>}
        {systemPackagesError && <div>{systemPackagesError}</div>}
        {loadingLanguagePackages && <div>Loading language packages...</div>}
        {languagePackagesError && <div>{languagePackagesError}</div>}
      </div>

      <div className="sidebar">
        <div className="flex flex-col gap-4 h-full overflow-y-auto">
          <span>System Dependencies</span>
          <ListBuilder
            items={systemPackages}
            onSelectionChange={setSelectedSystemDependencies}
            onSearchChange={setSystemSearchString}
          />

          <span>Language Dependencies</span>
          <ListBuilder
            items={languagePackages}
            onSelectionChange={setSelectedLanguageDependencies}
            onSearchChange={setLanguageSearchString}
            resetTrigger={resetLanguageDependencies}
            nonExistingPackages={existsResponse?.nonExistingPackages || []}
          />
        </div>
      </div>
    </div>
  );
};

export default App;
