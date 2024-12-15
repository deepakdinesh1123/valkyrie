import React, { useRef, useEffect, useState } from "react";

interface TerminalProps {
  output: string[];
  onInputChange: (input: string) => void;
  tabName?: string;
}

const Terminal: React.FC<TerminalProps> = ({ output, onInputChange, tabName }) => {
  const [activeTab, setActiveTab] = useState<'output' | 'input'>('output');
  const terminalRef = useRef<HTMLDivElement>(null);
  const [inputValue, setInputValue] = useState('');

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [output]);

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newValue = e.target.value;
    setInputValue(newValue);
    onInputChange(newValue);
  };

  return (
    <div className="terminal-container border rounded-lg bg-neutral-900 ">
      <div className="flex border-b border-neutral-700">
        <button
          onClick={() => setActiveTab('output')}
          className={`px-4 py-2 ${
            activeTab === 'output'
              ? 'bg-neutral-800 text-white border-b-2 border-blue-500'
              : 'text-white hover:text-white bg-neutral-800'
          }`}
        >
          {tabName || 'Output'}
        </button>
        <button
          onClick={() => setActiveTab('input')}
          className={`px-4 py-2 ${
            activeTab === 'input'
              ? 'bg-neutral-800 text-white border-b-2 border-blue-500'
              : 'text-white hover:text-white bg-neutral-800'
          }`}
        >
          Input
        </button>
      </div>

      {activeTab === 'output' ? (
        <div
          ref={terminalRef}
          className="p-4  overflow-y-auto font-mono text-sm"
        >
          {output.length === 0 ? (
            <div className="text-white">No output yet...</div>
          ) : (
            output.map((line, index) => (
              <div key={index} className="whitespace-pre-wrap text-white">
                {line}
              </div>
            ))
          )}
        </div>
      ) : (
        <div className="p-4 h-full">
          <textarea
            className="h-full w-full bg-neutral-800 text-white p-2 rounded border border-neutral-700 font-mono text-sm resize-none focus:outline-none focus:border-blue-500"
            placeholder="Enter your input here..."
            value={inputValue}
            onChange={handleInputChange}
          />
        </div>
      )}
    </div>
  );
};

export default Terminal;