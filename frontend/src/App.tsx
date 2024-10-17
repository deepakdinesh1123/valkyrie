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
              className="args-input w-50 bg-neutral-900 text-white border-none"
            />
          </div>
          <Button
            onClick={handleRunCode}
            disabled={isLoading}
            className={`run-code-btn ${isLoading ? 'loading' : ''}`}
          >
            {!isLoading && 'Run Code'}
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

      <div className="sidebar w-70  text-white h-screen p-6 flex flex-col justify-between">
        {/* Sidebar Header */}
        <div>
          <h2 className="font-bold border-b border-gray-700 pb-2">Menu</h2>
          <ul className="mt-6 space-y-2">
            <li
              className={`cursor-pointer px-3 py-2 rounded-md transition-colors ${sidebarOption === "addDependencies" ? "bg-neutral-800" : "hover:bg-stone-600"
                }`}
              onClick={() => setSidebarOption("addDependencies")}
            >
              Add Dependencies
            </li>
            <li
              className={`cursor-pointer px-3 py-2 rounded-md transition-colors ${isRequestModalOpen ? "bg-neutral-800" : "hover:bg-stone-600"
                }`}
              onClick={() => setIsRequestModalOpen(true)}
            >
              Request Package
            </li>
            <li
              className={`cursor-pointer px-3 py-2 rounded-md transition-colors ${sidebarOption === "help" ? "bg-neutral-800" : "hover:bg-stone-600"
                }`}
              onClick={() => setSidebarOption("help")}
            >
              Help
            </li>
          </ul>
        </div>

        {/* Sidebar Content */}
        <div className="flex-1 overflow-y-auto mt-1 border-t border-gray-700 pt-6">
          {sidebarOption === "addDependencies" && (
            <div className="flex flex-col gap-6">
              <span className="text-lg font-semibold">System Dependencies</span>
              <ListBuilder
                items={systemPackages}
                onSelectionChange={setSelectedSystemDependencies}
                onSearchChange={setSystemSearchString}
              />

              <span className="text-lg font-semibold">Language Dependencies</span>
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
        {/* Sidebar Footer */}
        <div className="text-sm text-neutral-500 border-t border-gray-700 pt-4">
          Valkyrie
        </div>
      </div>



      <Dialog open={isRequestModalOpen} onOpenChange={setIsRequestModalOpen}>
        <DialogContent className=" bg-black">
          <DialogHeader>
            <DialogTitle className="text-white">Request a Package</DialogTitle>
            <DialogDescription>
              If you need a package that's not available, please submit a request to our team. We'll review it and add it to our system if possible. Currently we are supporting only packages that are already available as <a className="underline text-blue-500" href="https://nixos.org/manual/nixpkgs/stable/#overview-of-nixpkgs">nixpkgs</a>, so if you are not sure if your package is available, head over to <a className="underline text-blue-500" href="https://search.nixos.org/packages">NixOS search</a> to check.
            </DialogDescription>
          </DialogHeader>
          <p className="text-white">To request a package, please fill out this <a className="underline text-blue-500" href="https://forms.gle/XpSVTpf3ix4rAjrr9">form</a> with the following information:</p>
          <ul className="list-disc pl-5 text-white">
            <li>Package details</li>
            <li>Nix channel version</li>
          </ul>
          <DialogClose asChild>
            <Button className="mt-4 border border-transparent hover:border-white transition-colors">Close</Button>
          </DialogClose>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default App;
