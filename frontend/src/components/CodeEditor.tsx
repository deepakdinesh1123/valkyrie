import React from "react";
import Editor, { EditorProps } from "@monaco-editor/react";
import { Language, LanguageVersion } from "@/api-client";

interface CodeEditorProps {
  selectedLanguage: Language;
  selectedLanguageVersion: LanguageVersion;
  onChange?: (content: string) => void;
  editorOptions?: EditorProps["options"];
  value?: string; 
}

const CodeEditor: React.FC<CodeEditorProps> = ({
  selectedLanguage,
  onChange,
  editorOptions,
  value, 
}) => {
  const handleEditorChange = (newValue: string | undefined) => {
    const newContent = newValue ?? "";
    onChange?.(newContent);
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E1E] text-white">
      <div className="flex h-14 px-4 border-b border-stone-700 pb-0 mb-0 pt-3">
        <div
          className="inline-block px-4 py-2 border border-stone-700 min-w-20"
          style={{ marginBottom: '-1px' }}
        >
          <button className="text-sm text-white bg-transparent border-none cursor-pointer focus:outline-none pb-0 mb-0">
            {`main.${selectedLanguage.extension}`}
          </button>
        </div>
      </div>
      <div className="flex-grow mt-0">
        <Editor
          height="100%"
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