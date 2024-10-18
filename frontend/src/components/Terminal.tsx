import React, { useRef, useEffect } from "react";

interface TerminalProps {
  output: string[];
  tabName?: string;
}

const Terminal: React.FC<TerminalProps> = ({ output, tabName }) => {
  const terminalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [output]);

  return (
    <div className="terminal-container">
      <div className="terminal-tab">
        <span className="font-medium">{tabName}</span>
      </div>
      <div
        ref={terminalRef}
        className="output-terminal"
      >
        {output.length === 0 ? (
          <div>No output yet...</div>
        ) : (
          output.map((line, index) => (
            <div key={index}>{line}</div>
          ))
        )}

      </div>
    </div>
  );
};

export default Terminal;