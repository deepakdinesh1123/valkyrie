import React from "react";
import Editor, { EditorProps } from "@monaco-editor/react";
import { Language, LanguageVersion } from "@/api-client";
import TabIcon from '@/assets/tabicon.svg'

interface CodeEditorProps {
  selectedLanguage: Language;
  selectedLanguageVersion: LanguageVersion;
  onChange?: (content: string) => void;
  editorOptions?: EditorProps["options"];
  value?: string;
  height: string;
}

const CodeEditor: React.FC<CodeEditorProps> = ({
  selectedLanguage,
  onChange,
  editorOptions,
  value,
  height,
}) => {
  const handleEditorChange = (newValue: string | undefined) => {
    const newContent = newValue ?? "";
    onChange?.(newContent);
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E1E] text-white">
      <div className="flex h-16  border-b border-stone-700 pb-0 mb-0 pt-3">
        <img src={TabIcon} className="ml-2" alt="Valkyrie" />
        <div
          className="inline-flex items-center px-3 py-2 border border-stone-700 min-w-20 ml-2"

        >
          <button className="text-sm text-white bg-transparent border-none cursor-pointer focus:outline-none pb-0 mb-0">
            {`main.${selectedLanguage.extension}`}
          </button>
        </div>
      </div>

      <div className="flex-grow mt-0">
        <Editor
          height={height}
          width="100%"
          language={selectedLanguage.monaco_language}
          value={value ?? selectedLanguage.default_code}
          onChange={handleEditorChange}
          theme="vs-dark"
          options={{
            minimap: { enabled: false },
            scrollBeyondLastLine: false,
            fontSize: 14,
            tabSize: 2,
            ...editorOptions,
          }}
        />
      </div>
    </div>
  );
};

export default CodeEditor;
