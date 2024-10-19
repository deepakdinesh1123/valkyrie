import React, { useState, useEffect } from "react";
import Editor, { EditorProps } from "@monaco-editor/react";
import { Language } from "@/api-client";
import { getLanguagePrefix } from "@/utils/prefix";


interface CodeEditorProps {
  languages: Language[];
  selectedLanguage: Language;
  onChange?: (content: string) => void;
  editorOptions?: EditorProps["options"];
}

const CodeEditor: React.FC<CodeEditorProps> = ({
  selectedLanguage,
  onChange,
  editorOptions,
}) => {
  const [content, setContent] = useState(selectedLanguage.defaultcode);
  const [previousPrefix, setPreviousPrefix] = useState(
    getLanguagePrefix(selectedLanguage.name)
  );

  useEffect(() => {
    const currentPrefix = getLanguagePrefix(selectedLanguage.name);

    if (currentPrefix !== previousPrefix) {
      setContent(selectedLanguage.defaultcode);
      setPreviousPrefix(currentPrefix);
    }
  }, [selectedLanguage]);

  const handleEditorChange = (newValue: string | undefined) => {
    const newContent = newValue ?? "";
    setContent(newContent);
    onChange?.(newContent);
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E1E] text-white">
      {/* VS Code-style filename tab */}
      <div className="flex h-14 px-4 border-b border-stone-700 pb-0 mb-0 pt-3">
        <div
          className="inline-block px-4 py-2 border min-w-20"
          style={{ marginBottom: '-1px' }}
        >
          <button className="text-sm text-white bg-transparent border-none cursor-pointer focus:outline-none pb-0 mb-0">
            {`main.${selectedLanguage.extension}`}
          </button>
        </div>

      </div>

      {/* Editor Section */}
      <div className="flex-grow mt-0">
        <Editor
          height="100%"
          width="100%"
          language={selectedLanguage.monacolanguage}
          value={content}
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
