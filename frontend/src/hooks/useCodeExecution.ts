import { useState, useEffect } from 'react';
import { api } from '@/utils/api';

export const useCodeExecution = () => {
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);
  const [eventSource, setEventSource] = useState<EventSource | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const eventPath = "http://localhost:8080/";

  const executeCode = async (runData: {
    language: string;
    code: string;
    environment: {
      systemDependencies: string[];
      languageDependencies: string[];
      args: string;
    }
    
  }) => {
    try {
      console.log(runData);
      setIsLoading(true); 
      const response = await api.execute(runData);
      setTerminalOutput(['Processing...']);

      const source = new EventSource(`${eventPath}${response.data.events}`);
      setEventSource(source);

      source.onmessage = (event) => {
        const data = JSON.parse(event.data);

        if (data.status === 'completed') {
          setTerminalOutput((prev) => [
            ...prev.slice(0, prev.length - 1),
            data.logs || 'No logs available.',
          ]);
          setIsLoading(false); 
          source.close();
        }
      };

      source.onerror = (error) => {
        console.error('EventSource error:', error);
        setTerminalOutput((prev) => [...prev, 'EventSource connection error.']);
        setIsLoading(false); 
        source.close();
      };
    } catch (error) {
      console.error('Execution failed:', error);
      setTerminalOutput((prev) => [...prev, 'Execution failed.']);
      setIsLoading(false); 
    }
  };

  useEffect(() => {
    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  }, [eventSource]);

  return { terminalOutput, executeCode, isLoading };
};
