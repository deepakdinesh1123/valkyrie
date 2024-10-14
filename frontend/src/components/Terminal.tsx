import React from "react";

interface TerminalProps {
  output: string[];
}

const Terminal: React.FC<TerminalProps> = ({ output }) => {
  return (
    <div className="terminal-container bg-black text-green-500 p-4 h-64 overflow-y-auto font-mono">
      {output.length === 0 ? (
        <div className="text-gray-500"></div>
      ) : (
        output.map((line, index) => <div key={index}>{line}</div>)
      )}
    </div>
  );
};

export default Terminal;
