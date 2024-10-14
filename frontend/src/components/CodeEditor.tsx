import React, { useState, useEffect } from "react";
import Editor, { EditorProps } from "@monaco-editor/react";

interface Language {
  name: string;
  extension: string;
  monacoLanguage: string;
  defaultCode: string;
}

interface CodeEditorProps {
  languages: Language[];
  selectedLanguage: Language;
  onChange?: (tabName: string, content: string) => void;
  editorOptions?: EditorProps["options"];
}

const CodeEditor: React.FC<CodeEditorProps> = ({
  selectedLanguage,
  onChange,
  editorOptions,
}) => {
  const [content, setContent] = useState(selectedLanguage.defaultCode);

  useEffect(() => {
    setContent(selectedLanguage.defaultCode);
  }, [selectedLanguage]);

  const handleEditorChange = (newValue: string | undefined) => {
    const newContent = newValue || "";
    setContent(newContent);
    onChange?.(`main.${selectedLanguage.extension}`, newContent);
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E1E] text-white">
      <div className="flex justify-between items-center bg-[#252526] p-2">
        <div className="flex">
          <button className="px-4 py-2 cursor-pointer text-sm bg-[#1E1E1E] text-white border-t-2 border-blue-500">
            {`main.${selectedLanguage.extension}`}
          </button>
        </div>
      </div>
      <div className="flex-grow">
        <Editor
          height="100%"
          width="100%"
          language={selectedLanguage.monacoLanguage}
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
