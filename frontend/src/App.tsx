import React, { useCallback, useEffect, useState } from "react";
import CodeEditor from "@/components/CodeEditor";
import ListBuilder from "@/components/ListBuilder";
import Terminal from "@/components/Terminal";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import "@/App.css";
import { useLanguages } from '@/hooks/useLanguages';
import { useCodeExecution } from '@/hooks/useCodeExecution';
import { useSystemPackages } from '@/hooks/useSystemPackages';
import { useLanguagePackages } from '@/hooks/useLanguagePackages';
import { LanguageSelector } from '@/components/LanguageSelector';
import { usePackagesExist } from "@/hooks/usePackageExists";
import HelpModal from "@/components/HelpModal";
import RequestPackageModal from "@/components/RequestModal";
import {
  enable as enableDarkMode,
} from 'darkreader';


const App: React.FC = () => {
  const { languages, selectedLanguage, setSelectedLanguage } = useLanguages();
  const { terminalOutput, executeCode, isLoading } = useCodeExecution();
  const [terminalHeight, setTerminalHeight] = useState<number>(300);
  const [codeContent, setCodeContent] = useState<string>("");
  const [args, setArgs] = useState<string>("");
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] = useState<string[]>([]);
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<string[]>([]);
  const [systemSearchString, setSystemSearchString] = useState<string>("");
  const [languageSearchString, setLanguageSearchString] = useState<string>("");
  const [resetLanguageDependencies, setResetLanguageDependencies] = useState({});
  const { systemPackages } = useSystemPackages(systemSearchString);
  const { languagePackages, resetLanguagePackages } = useLanguagePackages(languageSearchString, selectedLanguage?.searchquery);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [selectedLanguagePrefix, setSelectedLanguagePrefix] = useState<string>("");
  const [pendingLanguageChange, setPendingLanguageChange] = useState<any>(null);
  const { existsResponse } = usePackagesExist(
    pendingLanguageChange?.searchquery || "",
    selectedLanguageDependencies
  );
  const [isRequestModalOpen, setIsRequestModalOpen] = useState<boolean>(false);
  const [isHelpModalOpen, setIsHelpModalOpen] = useState<boolean>(false);



  useEffect(() => {
    enableDarkMode({
      brightness: 100,
      contrast: 90,
      sepia: 10,
    });
    if (pendingLanguageChange && existsResponse) {
      handleLanguageChangeEffect(pendingLanguageChange, existsResponse);
      setPendingLanguageChange(null);
    }
  }, [pendingLanguageChange, existsResponse]);

  const handleTerminalResize = useCallback((height: number) => {
    setTerminalHeight(height);
  }, []);

  const handleEditorChange = useCallback((content: string) => {
    setCodeContent(content);
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
  }, [selectedLanguage, codeContent, selectedSystemDependencies, selectedLanguageDependencies, args, executeCode]);

  const resetOnNewLanguage = useCallback((language: any) => {
    setSelectedLanguage(language);
    setCodeContent(language.defaultcode);
    setLanguageSearchString("");
    setSelectedLanguageDependencies([]);
    setResetLanguageDependencies({});
    resetLanguagePackages();
  }, [setSelectedLanguage, resetLanguagePackages]);

  const handleLanguageChangeEffect = useCallback((language: any, existsResponse: any) => {
    const newPrefix = language.name.split("-")[0];

    if (newPrefix !== selectedLanguagePrefix) {
      setSelectedLanguagePrefix(newPrefix);
      resetOnNewLanguage(language);
    } else if (language.name !== selectedLanguage?.name) {
      if (existsResponse.exists) {
        setSelectedLanguage(language);
        setCodeContent(language.defaultcode);
      } else {
        const nonExistingPackages = existsResponse.nonExistingPackages;
        setSelectedLanguageDependencies(prev => prev.filter(dep => !nonExistingPackages.includes(dep)));
        setSelectedLanguage(language);
        setCodeContent(language.defaultcode);
      }
    }
  }, [selectedLanguagePrefix, selectedLanguage, resetOnNewLanguage, setSelectedLanguage]);

  const handleLanguageChange = useCallback((language: any) => {
    setPendingLanguageChange(language);
  }, []);

  return (
    <div className="flex h-screen">
      {/* Editor Container */}
      <div className={`editor-container flex-1 transition-all duration-300 ${isSidebarOpen ? "w-2/3" : "w-full"}`}>
        <div className="top-bar flex justify-between items-center h-4 pt-7 pr-0 mr-11 bg-transparent">
          <div className="language-args-container flex items-center ">
            <LanguageSelector
              languages={languages}
              selectedLanguage={selectedLanguage}
              onLanguageChange={handleLanguageChange}
            />
            <Input
              type="text"
              placeholder="Args"
              value={args}
              onChange={(e) => setArgs(e.target.value)}
              className="args-input w-50 bg-neutral-900 text-white border-none ml-2"
            />
          </div>
          <Button
            onClick={handleRunCode}
            disabled={isLoading}
            className={`run-code-btn ${isLoading ? 'loading' : ''}`}
          >
            {!isLoading && 'Run Code'}
          </Button>

          {/* Menu Toggle Button */}
          <Button
            onClick={() => setIsSidebarOpen(!isSidebarOpen)} // Toggle sidebar visibility
            className="menu-toggle-btn mx-2 min-w-300"
          >
            {isSidebarOpen ? "Menu" : "Menu"}
          </Button>
        </div>

        {/* Resizable Editor */}
        <div className="flex-grow overflow-hidden" style={{ height: `calc(100% - ${terminalHeight}px - 4rem)` }}>
          <CodeEditor
            languages={languages}
            selectedLanguage={selectedLanguage}
            onChange={handleEditorChange}
            editorOptions={{ wordWrap: "on" }}
          />
        </div>

        {/* Resizable Terminal */}
        <div className="relative pt-4 pr-1 pb-0" style={{ height: `${terminalHeight}px` }}>
          <div
            className="absolute top-0 left-0 right-0 h-1 bg-gray-600 cursor-n-resize"
            onMouseDown={(e) => {
              const startY = e.clientY;
              const startHeight = terminalHeight;
              const handleMouseMove = (e: MouseEvent) => {
                const deltaY = startY - e.clientY;
                const newHeight = Math.max(100, Math.min(startHeight + deltaY, window.innerHeight - 200));
                handleTerminalResize(newHeight);
              };
              const handleMouseUp = () => {
                document.removeEventListener('mousemove', handleMouseMove);
                document.removeEventListener('mouseup', handleMouseUp);
              };
              document.addEventListener('mousemove', handleMouseMove);
              document.addEventListener('mouseup', handleMouseUp);
            }}
          />
          <div className="terminal-container border-none" style={{ height: 'calc(100% - 20px)' }}>
            <Terminal output={terminalOutput} tabName="Output" />
          </div>
        </div>
      </div>

      {/* Collapsible Sidebar */}
      {
        isSidebarOpen && (
          <div className="sidebar w-1/3 text-white p-6 pt-1 flex flex-col justify-between transition-all duration-300">
            {/* Sidebar Content */}
            <div className="flex-1">
              <div className="flex flex-col gap-4 h-full">
                {/* System Dependencies Section */}
                <div className="flex-1 flex flex-col rounded-md shadow-md p-4">
                  <span className="">System Dependencies</span>
                  <div className="flex-1 mt-2 rounded-md min-h-56">
                    <ListBuilder
                      items={systemPackages}
                      onSelectionChange={setSelectedSystemDependencies}
                      onSearchChange={setSystemSearchString}
                    />
                  </div>
                </div>
                <div className="border border-t border-zinc-700"></div>

                {/* Language Dependencies Section */}
                <div className="flex-1 flex flex-col bg-transparent rounded-md shadow-md p-4">
                  <span className="">Language Dependencies</span>
                  <div className="flex-1 overflow-y-auto mt-2 rounded-md min-h-56">
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
            </div>

            {/* Sidebar Buttons */}
            <div className="flex space-x-2 mt-6">
              <button
                className={`flex-grow cursor-pointer px-3 py-2 bg-neutral-900 rounded-md transition-colors hover:bg-stone-600`}
                onClick={() => setIsRequestModalOpen(true)}
              >
                Request Package
              </button>
              <button
                className={`flex-grow cursor-pointer px-3 py-2 bg-neutral-900 rounded-md transition-colors hover:bg-stone-600`}
                onClick={() => setIsHelpModalOpen(true)}
              >
                Help
              </button>
            </div>


            {/* Sidebar Footer */}
            <div className="text-sm text-white border-t border-gray-700 pt-1">
              Valkyrie
            </div>

          </div>
        )
      }

      <HelpModal
        isOpen={isHelpModalOpen}
        onClose={() => setIsHelpModalOpen(false)}
      />

      <RequestPackageModal
        isOpen={isRequestModalOpen}
        onClose={() => setIsRequestModalOpen(false)}
      />
    </div>
  );

};

export default App;