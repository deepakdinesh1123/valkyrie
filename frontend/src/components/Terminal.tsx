import React, { useRef, useEffect } from "react";

interface TerminalProps {
  output: string[];
}

const Terminal: React.FC<TerminalProps> = ({ output }) => {
  const terminalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (terminalRef.current) {
      terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
    }
  }, [output]);

  return (
    <div
      ref={terminalRef}
      className="output-terminal font-mono text-green-500 bg-black"
    >
      {output.length === 0 ? (
        <div className="text-green-500">No output yet...</div>
      ) : (
        output.map((line, index) => <div key={index}>{line}</div>)
      )}
    </div>
  );
};

export default Terminal;
