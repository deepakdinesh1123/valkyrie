import React, { useCallback, useEffect, useState } from "react";
import LanguageSelector from "@/components/LanguageSelector";
import '@/App.css';
import { Input } from "@/components/ui/input";
import CodeEditor from "@/components/CodeEditor";
import { ServerStatusProvider } from "@/components/ServerStatusProvider";
import AppContainer from "@/components/AppContainer";
import { useLanguages } from "@/hooks/useLanguages";
import { useLanguageVersions } from "@/hooks/useLanguageVersions";
import { LanguageResponse, LanguageVersion } from "@/api-client";
import { Button } from "@/components/ui/button";
import { useCodeExecution } from "@/hooks/useCodeExecution";
import Terminal from "@/components/Terminal";
import HelpModal from "@/components/HelpModal";
import RequestPackageModal from "@/components/RequestModal";
import Announcement from "@/components/Announcement";
import PackageIcon from '@/assets/package.svg'
import HelpIcon from '@/assets/help.svg'
import ValkyrieIcon from '@/assets/valkyrie.svg'
import ListBuilder from "@/components/ListBuilder";
import { useSystemPackages } from "@/hooks/useSystemPackages";
import { useLanguagePackages } from "@/hooks/useLanguagePackages";
import { usePackagesExist } from "@/hooks/usePackageExists";
import { useDefaultPackages } from "@/hooks/useDefaultPackages";
import Editor from '@monaco-editor/react';

const App: React.FC = () => {
  const [args, setArgs] = useState<string>("");
  const [compilerArgs, setCompilerArgs] = useState<string>("");
  const { selectedLanguage, setSelectedLanguage } = useLanguages();
  const { selectedLanguageVersion, setSelectedLanguageVersion } = useLanguageVersions(selectedLanguage?.id || 0);
  const [codeContent, setCodeContent] = useState<string>("");
  const { terminalOutput, executeCode, isLoading } = useCodeExecution();
  const [selectedLanguageDependencies, setSelectedLanguageDependencies] = useState<string[]>([]);
  const [selectedSystemDependencies, setSelectedSystemDependencies] = useState<string[]>([]);
  const [terminalHeight, setTerminalHeight] = useState<number>(300);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const [systemSearchString, setSystemSearchString] = useState<string>("");
  const [languageSearchString, setLanguageSearchString] = useState<string>("");
  const [isRequestModalOpen, setIsRequestModalOpen] = useState<boolean>(false);
  const [isHelpModalOpen, setIsHelpModalOpen] = useState<boolean>(false);
  const calculatedHeight = `calc(100% - ${terminalHeight}px)`;
  const [isAnnouncmentModalOpen, setIsAnnouncementModalOpen] = useState<boolean>(false);
  const { systemPackages } = useSystemPackages(systemSearchString);
  const { languagePackages, resetLanguagePackages } = useLanguagePackages(languageSearchString, selectedLanguageVersion?.search_query || "");
  const { defaultSystemPackages, defaultLanguagePackages } = useDefaultPackages(selectedLanguageVersion?.search_query || "");
  const [pendingVersionChange, setPendingVersionChange] = useState<any>(null);
  const [resetLanguageDependencies, setResetLanguageDependencies] = useState({});
  const { existsResponse } = usePackagesExist(
    pendingVersionChange?.search_query || "",
    selectedLanguageDependencies
  );
  const finalSystemPackages = systemPackages.length > 0 ? systemPackages : defaultSystemPackages;
  const finalLanguagePackages = languagePackages.length > 0 ? languagePackages : defaultLanguagePackages;
  const [currentInput, setCurrentInput] = useState('');
  const [setup, setSetup] = useState('');
  const [command, setCommand] = useState('');
  const [activeTab, setActiveTab] = useState('dependencies');

  useEffect(() => {
    if (pendingVersionChange && existsResponse) {
      handleLanguageChangeEffect(pendingVersionChange, existsResponse);
      setPendingVersionChange(null);
    }
  }, [pendingVersionChange, existsResponse]);


  const resetOnNewLanguage = useCallback((language: LanguageResponse) => {
    setSelectedLanguage(language);
    setCodeContent(language.default_code);
    setSelectedLanguageDependencies([]);
    resetLanguagePackages();
    setResetLanguageDependencies({});
  }, [setSelectedLanguage, resetLanguagePackages]);

  const handleLanguageChangeEffect = useCallback(
    (version: any, existsResponse: any) => {
      if (version !== selectedLanguageVersion) {
        if (existsResponse.exists) {
          setSelectedLanguageVersion(version);
        } else {
          const nonExistingPackages = existsResponse.nonExistingPackages;
          setSelectedLanguageDependencies(prev =>
            prev.filter(dep => !nonExistingPackages.includes(dep))
          );
          setSelectedLanguageVersion(version);
        }
      }
    },
    [
      selectedLanguage,
      selectedLanguageVersion,
      resetOnNewLanguage,
      setSelectedLanguage,
      setSelectedLanguageVersion
    ]
  );


  useEffect(() => {
    if (selectedLanguage) {
      setCodeContent(selectedLanguage.default_code);
    }
  }, [selectedLanguage]);

  const handleEditorChange = useCallback((content: string) => {
    setCodeContent(content);
  }, []);

  const handleLanguageChange = useCallback((language: LanguageResponse) => {
    resetOnNewLanguage(language);
  }, [resetOnNewLanguage]);

  const handleVersionChange = useCallback((version: LanguageVersion) => {
    if (version !== selectedLanguageVersion) {
      setPendingVersionChange(version);
      const nonExistingPackages = existsResponse?.nonExistingPackages || [];
      setSelectedLanguageVersion(version);
      setSelectedLanguageDependencies((prev) => prev.filter(dep => !nonExistingPackages.includes(dep)));
    }
  }, [selectedLanguageVersion, existsResponse]);


  const handleRunCode = useCallback(() => {
    console.log(currentInput);

    if (selectedLanguage && codeContent) {
      executeCode({
        language: selectedLanguage.name,
        version: selectedLanguageVersion?.version || "",
        code: codeContent,
        environment: {
          systemDependencies: selectedSystemDependencies,
          languageDependencies: selectedLanguageDependencies,
          setup: setup,
        },
        cmdLineArgs: args,
        input: currentInput,
        compilerArgs: compilerArgs,
        command: command
      });
    } else {
    }
  }, [selectedLanguage, codeContent, args, compilerArgs, setup, command, executeCode, selectedLanguageDependencies, selectedSystemDependencies]);

  const handleTerminalResize = useCallback((height: number) => {
    setTerminalHeight(height);
  }, []);

  const handleInputChange = (input: string) => {
    setCurrentInput(input);
  };


  return (
    <ServerStatusProvider>
      <AppContainer>
    <div className="flex h-screen relative">
      <div className="absolute top-0 right-0 z-10 ">
        <img
          src={ValkyrieIcon}
          className="h-12 object-cove mr-2 pt-3 pr-20"
          alt="Valkyrie"
          style={{ objectPosition: 'top' }}
        />
      </div>
      <div className="editor-container flex-1 w-full">
        <div className="top-bar flex flex-wrap justify-between items-center p-2 pt-4 bg-transparent mr-32">
          <div className="flex flex-wrap items-center w-full sm:w-auto  sm:mb-0">
            <div className="w-full sm:w-auto  sm:mb-0 sm:mr-2 border-none">
              <LanguageSelector
                onLanguageChange={handleLanguageChange}
                selectedLanguage={selectedLanguage}
                selectedLanguageVersion={selectedLanguageVersion}
                onVersionChange={handleVersionChange}
              />
            </div>
            <Input
              type="text"
              placeholder="Args"
              value={args}
              onChange={(e) => setArgs(e.target.value)}
              className="args-input w-full sm:w-36 mr-1 bg-neutral-900 text-white border-opacity-100 focus:ring-0"
            />
            <Input
              type="text"
              placeholder="Compiler Args"
              value={compilerArgs}
              onChange={(e) => setCompilerArgs(e.target.value)}
              className="args-input w-full sm:w-36 mr-1 bg-neutral-900 text-white border-opacity-100 focus:ring-0"
            />
          </div>
          <div className="flex items-center w-full sm:w-auto justify-end">
            <Button
              onClick={handleRunCode}
              disabled={isLoading || !codeContent}
              className={`run-code-btn mr-2 ${isLoading ? 'loading' : ''} w-1/2 sm:w-auto bg-neutral-900 transition-colors hover:bg-stone-600 text-sm active:bg-neutral-900`}
            >
              {isLoading ? "" : 'Run Code'}
            </Button>
            <Button
              onClick={() => setIsSidebarOpen(!isSidebarOpen)}
              className="menu-toggle-btn w-1/2 sm:w-auto bg-neutral-900 transition-colors hover:bg-stone-600 text-sm active:bg-neutral-900"
            >
              {isSidebarOpen ? "Menu" : "Menu"}
            </Button>
          </div>
        </div>
        <div
          className="flex-grow overflow-hidden transition-all duration-300 ease-in-out"
          style={{
            height: calculatedHeight,
            width: isSidebarOpen ? '74%' : '100%',
          }}
        >
          <CodeEditor
            selectedLanguage={selectedLanguage}
            selectedLanguageVersion={selectedLanguageVersion}
            onChange={handleEditorChange}
            value={codeContent}
            editorOptions={{ wordWrap: "on" }}
            height={calculatedHeight}
          />
        </div>
        <div className="relative" style={{ height: `${terminalHeight}px` }}>
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
          <div className="terminal-container transition-all duration-300 ease-in-out mr-0" style={{ height: 'calc(100%)', width: isSidebarOpen ? '75%' : '100%' }}>
            <Terminal output={terminalOutput} tabName="Output" onInputChange={handleInputChange} />
          </div>
        </div>
      </div>

      <div
        className={`fixed top-0 right-0 h-full w-full md:w-1/3 lg:w-1/4 bg-neutral-800 text-white p-2 
        flex flex-col justify-between transition-transform duration-300 ease-in-out z-50
        ${isSidebarOpen ? 'translate-x-0' : 'translate-x-full'}
        shadow-2xl`}
      >
        {/* Tab Navigation */}
        <div className="flex border-b  border-zinc-700 mb-2 ">
          <button
            className={`flex-1 py-2 bg-neutral-900 ${activeTab === 'dependencies' ? 'border-b-2 border-blue-500' : ''}`}
            onClick={() => setActiveTab('dependencies')}
          >
            Dependencies
          </button>
          <button
            className={`flex-1 py-2 bg-neutral-900 ${activeTab === 'setup' ? 'border-b-2 border-blue-500' : ''}`}
            onClick={() => setActiveTab('setup')}
          >
            Setup
          </button>
        </div>

        {/* Dependencies Tab Content */}
        {activeTab === 'dependencies' && (
          <div className="flex-1">
            <div className="flex flex-col gap-2 h-full">
              <div className="flex-1 flex flex-col rounded-md shadow-md p-2">
                <span className="">System Dependencies</span>
                <div className="flex-1 mt-2 rounded-md min-h-[14rem]">
                  <ListBuilder
                    items={finalSystemPackages}
                    onSelectionChange={setSelectedSystemDependencies}
                    onSearchChange={setSystemSearchString}
                  />
                </div>
              </div>
              <div className="border border-t border-zinc-700"></div>
              <div className="flex-1 flex flex-col bg-transparent rounded-md shadow-md p-2">
                <span className="">Language Dependencies</span>
                <div className="flex-1 mt-2 rounded-md min-h-[14rem]">
                  <ListBuilder
                    items={finalLanguagePackages}
                    onSelectionChange={setSelectedLanguageDependencies}
                    onSearchChange={setLanguageSearchString}
                    resetTrigger={resetLanguageDependencies}
                    nonExistingPackages={existsResponse?.nonExistingPackages || []}
                  />
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Setup Tab Content */}
        {activeTab === 'setup' && (
          <div className="flex-1 flex flex-col gap-4 p-2">
            <div className="flex flex-col flex-1">
              <label className="mb-2 text-sm">Setup</label>
              <Editor
                className="flex-1 w-full bg-[#1E1E1E] text-white p-2 rounded-md resize-none"
                width="100%"
                language="shell" 
                theme="vs-dark"  
                value={setup}
                onChange={(value) => setSetup(value || '')}
                options={{
                  readOnly: false,                 
                  lineNumbers: 'off',              
                  minimap: { enabled: false },    
                  scrollBeyondLastLine: false,     
                  wordWrap: 'off',                 
                  folding: false,                  
                  lineDecorationsWidth: 0,         
                  lineNumbersMinChars: 0,          
                  guides: {
                    indentation: false             
                  },
                  renderLineHighlight: 'none',     
                  overviewRulerLanes: 0,           
                  contextmenu: false,              
                  suggestOnTriggerCharacters: false, 
                  quickSuggestions: false,         
                  automaticLayout: true,         
                  cursorBlinking: 'solid',         
                  cursorStyle: 'line',             
                }}
              />
            </div>
            <div className="flex flex-col flex-1">
              <label className="mb-2 text-sm">Commands</label>
              <Editor
                className="flex-1 w-full bg-[#1E1E1E] text-white p-2 rounded-md resize-none"
                width="100%"
                language="shell" 
                theme="vs-dark"  
                value={command}
                onChange={(value) => setCommand(value || '')}
                options={{
                  readOnly: false,                 
                  lineNumbers: 'off',              
                  minimap: { enabled: false },    
                  scrollBeyondLastLine: false,     
                  wordWrap: 'off',                 
                  folding: false,                  
                  lineDecorationsWidth: 0,         
                  lineNumbersMinChars: 0,          
                  guides: {
                    indentation: false             
                  },
                  renderLineHighlight: 'none',     
                  overviewRulerLanes: 0,           
                  contextmenu: false,              
                  suggestOnTriggerCharacters: false, 
                  quickSuggestions: false,         
                  automaticLayout: true,         
                  cursorBlinking: 'solid',         
                  cursorStyle: 'line',             
                }}
              />
            </div>
          </div>
        )}

        <div className="flex flex-col sm:flex-row space-y-2 sm:space-y-0 sm:space-x-2 mt-4 pb-2">
          <button
            className="flex-grow flex items-center justify-center cursor-pointer px-3 py-2 bg-neutral-900 rounded-md transition-colors hover:bg-stone-600 text-sm"
            onClick={() => setIsRequestModalOpen(true)}
          >
            <img src={PackageIcon} alt="Package" className="h-4 w-4 mr-2" />
            Request Package
          </button>
          <button
            className="flex-grow flex items-center justify-center cursor-pointer px-3 py-2 bg-neutral-900 rounded-md transition-colors hover:bg-stone-600 text-sm"
            onClick={() => setIsHelpModalOpen(true)}
          >
            <img src={HelpIcon} alt="Help" className="h-4 w-4 mr-2" />
            Help
          </button>
        </div>
        <div className="border border-t border-zinc-700"></div>
        <div className="text-sm text-white pt-1 flex items-center space-x-2 justify-between">
          <span>Valkyrie</span>
          <a
            href="https://discord.gg/9VMKq3rt5C"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:opacity-80 transition-opacity"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="currentColor"
              className="w-5 h-5 text-white"
            >
              <path
                d="M20.317 4.369a19.791 19.791 0 00-4.885-1.527.074.074 0 00-.079.037c-.21.375-.444.864-.608 1.248a18.292 18.292 0 00-5.487 0 12.327 12.327 0 00-.617-1.248.079.079 0 00-.079-.037A19.425 19.425 0 003.68 4.369a.07.07 0 00-.032.027C.533 9.39-.32 14.313.099 19.163a.082.082 0 00.031.058 19.875 19.875 0 005.996 3.03.079.079 0 00.084-.027c.464-.637.873-1.312 1.226-2.016a.074.074 0 00-.041-.105 13.12 13.12 0 01-1.872-.9.076.076 0 01-.008-.126c.125-.094.25-.191.371-.292a.073.073 0 01.077-.01c3.927 1.793 8.18 1.793 12.061 0a.073.073 0 01.079.009c.122.1.247.198.372.292a.076.076 0 01-.007.125 12.663 12.663 0 01-1.873.901.075.075 0 00-.04.105c.366.704.776 1.379 1.224 2.016a.079.079 0 00.084.028 19.875 19.875 0 005.997-3.03.08.08 0 00.031-.058c.5-5.192-.83-10.058-3.575-14.767a.061.061 0 00-.03-.028zM8.02 15.331c-1.182 0-2.158-1.085-2.158-2.419 0-1.333.953-2.418 2.158-2.418 1.21 0 2.174 1.09 2.158 2.418 0 1.334-.953 2.419-2.158 2.419zm7.974 0c-1.182 0-2.158-1.085-2.158-2.419 0-1.333.953-2.418 2.158-2.418 1.21 0 2.174 1.09 2.158 2.418 0 1.334-.953 2.419-2.158 2.419z"
              />
            </svg>
          </a>
        </div>
      </div>
      <HelpModal
        isOpen={isHelpModalOpen}
        onClose={() => setIsHelpModalOpen(false)}
      />

      <RequestPackageModal
        isOpen={isRequestModalOpen}
        onClose={() => setIsRequestModalOpen(false)}
      />

      <Announcement
        isOpen={isAnnouncmentModalOpen}
        onClose={() => setIsAnnouncementModalOpen(false)}
      />
    </div>
    </AppContainer>
    </ServerStatusProvider>
  );
};

export default App;