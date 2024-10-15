import { useState } from 'react';
import { api } from '@/utils/api';

export const useCodeExecution = () => {
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);

  const executeCode = async (runData: {
    language: string;
    code: string;
    systemDependencies: string[];
    languageDependencies: string[];
  }) => {
    try {
      const response = await api.execute(runData);
      const jobOutput = `Job ID: ${response.data.jobId}\nEvents URL: ${response.data.events}`;
      setTerminalOutput((prev) => [...prev, jobOutput]);
    } catch (error) {
      console.error('Execution failed:', error);
      setTerminalOutput((prev) => [...prev, 'Execution failed.']);
    }
  };

  return { terminalOutput, executeCode };
};