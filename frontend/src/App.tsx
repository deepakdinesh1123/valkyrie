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
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogClose } from "@/components/ui/dialog";
import { HelpComponent } from "@/components/HelpComponent";
import { Input } from "./components/ui/input";


const App: React.FC = () => {
  const { languages, selectedLanguage, setSelectedLanguage } = useLanguages();
  const { terminalOutput, executeCode, isLoading } = useCodeExecution();

  const [codeContent, setCodeContent] = useState<string>("");
  const [args, setArgs] = useState<string>("");
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] = useState<string[]>([]);
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<string[]>([]);
  const [systemSearchString, setSystemSearchString] = useState<string>("");
  const [languageSearchString, setLanguageSearchString] = useState<string>("");
  const [resetLanguageDependencies, setResetLanguageDependencies] = useState({});
  const { systemPackages } = useSystemPackages(systemSearchString);
  const { languagePackages, resetLanguagePackages } = useLanguagePackages(languageSearchString, selectedLanguage?.searchquery);

  const [selectedLanguagePrefix, setSelectedLanguagePrefix] = useState<string>("");
  const [pendingLanguageChange, setPendingLanguageChange] = useState<any>(null);
  const { existsResponse } = usePackagesExist(
    pendingLanguageChange?.searchquery || "",
    selectedLanguageDependencies
  );

  const [sidebarOption, setSidebarOption] = useState<string>("addDependencies");
  const [isRequestModalOpen, setIsRequestModalOpen] = useState<boolean>(false);

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
        environment: {
          systemDependencies: selectedSystemDependencies,
          languageDependencies: selectedLanguageDependencies,
          args: args
        }

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
        <div className="top-bar">
          <div className="language-args-container">
            <div className="language-selector">
              <LanguageSelector
                languages={languages}
                selectedLanguage={selectedLanguage}
                onLanguageChange={handleLanguageChange}
              />
            </div>
            <Input
              type="text"
              placeholder="Args"
              value={args}
              onChange={(e) => setArgs(e.target.value)}
              className="args-input w-50"
            />
          </div>
          <Button
            onClick={handleRunCode}
            disabled={isLoading}
            className={`run-code-btn ${isLoading ? 'loading' : ''}`}
          >
            {!isLoading && 'Run'}
          </Button>
        </div>

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
      </div>

      <div className="sidebar">
        <div className="flex flex-col gap-4 h-full">
          <div className="flex flex-col gap-2">
            <Button onClick={() => setSidebarOption("addDependencies")}>Add Dependencies</Button>
            <Button onClick={() => setIsRequestModalOpen(true)}>Request Package</Button>
            <Button onClick={() => setSidebarOption("help")}>Help</Button>
          </div>

          {sidebarOption === "addDependencies" && (
            <div className="flex flex-col gap-4 overflow-y-auto">
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
          )}

          {sidebarOption === "help" && <HelpComponent />}
        </div>
      </div>

      <Dialog open={isRequestModalOpen} onOpenChange={setIsRequestModalOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Request a Package</DialogTitle>
            <DialogDescription>
              If you need a package that's not available, please submit a request to our team. We'll review it and add it to our system if possible.
            </DialogDescription>
          </DialogHeader>
          <p>To request a package, please contact our support team at support@example.com with the following information:</p>
          <ul className="list-disc pl-5">
            <li>Package name</li>
            <li>Programming language</li>
            <li>Brief description of why you need this package</li>
          </ul>
          <DialogClose asChild>
            <Button className="mt-4">Close</Button>
          </DialogClose>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default App;
