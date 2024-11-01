import React, { useCallback, useEffect, useState } from "react";
import LanguageSelector from "@/components/LanguageSelector";
import '@/App.css';
import { Input } from "@/components/ui/input";
import CodeEditor from "@/components/CodeEditor";
import { useLanguages } from "@/hooks/useLanguages";
import { useLanguageVersions } from "@/hooks/useLanguageVersions";

const App: React.FC = () => {
  const [args, setArgs] = useState<string>("");
  const { selectedLanguage} = useLanguages();
  const { selectedLanguageVersion } = useLanguageVersions(selectedLanguage?.id);
  const [codeContent, setCodeContent] = useState<string>("");

  const handleEditorChange = useCallback((content: string) => {
    setCodeContent(content);
  }, []);

  // Log selected language and version when they change
  useEffect(() => {
    if (selectedLanguage) {
      console.log(`Selected Language: ${selectedLanguage.name}, Version: ${selectedLanguageVersion.version}`);
    }
  }, [selectedLanguage, selectedLanguageVersion]);

  return (
    <div className="editor-container flex-1 w-full">
      <div className="top-bar flex flex-wrap justify-between items-center p-2 bg-transparent mr-20">
        <div className="flex flex-wrap items-center w-full sm:w-auto mb-2 sm:mb-0">
          <div className="w-full sm:w-auto mb-2 sm:mb-0 sm:mr-2 border-none">
            <LanguageSelector />
          </div>
          <Input
            type="text"
            placeholder="Args"
            value={args}
            onChange={(e) => setArgs(e.target.value)}
            className="args-input w-full sm:w-36 mr-1 bg-neutral-900 text-white border-opacity-100 focus:ring-0"
          />
        </div>
      </div>
      <div className="flex-grow overflow-hidden">
        <CodeEditor
          selectedLanguage={selectedLanguage}
          selectedVersion={selectedLanguageVersion}
          onChange={handleEditorChange}
          editorOptions={{ wordWrap: "on" }}
        />
      </div>
    </div>
  );
};

export default App;
