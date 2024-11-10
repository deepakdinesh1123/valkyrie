import { useState, useEffect } from 'react';
import { api } from '@/utils/api';

export const useCodeExecution = () => {
  const [terminalOutput, setTerminalOutput] = useState<string[]>([]);
  const [eventSource, setEventSource] = useState<EventSource | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const eventPath = import.meta.env.VITE_BASE_PATH;

  const executeCode = async (runData: {
    language: string;
    version: string;
    code: string;
    environment: {
      systemDependencies: string[];
      languageDependencies: string[];
    },
    cmdLineArgs: string;
    
  }) => {
    try {
      console.log(runData);
      setIsLoading(true); 
      const response = await api.execute(runData);

      const source = new EventSource(`${eventPath}${response.data.events}`);
      setEventSource(source);

      source.onmessage = (event) => {
        const data = JSON.parse(event.data);

        switch (data.status) {
          case 'pending':
            setTerminalOutput([ 'Waiting for worker']);
            break;
          case 'scheduled':
            setTerminalOutput([ 'Processing...']);
            break;
          case 'completed':
            setTerminalOutput([data.logs || 'No logs available.',]);
            setIsLoading(false); 
            source.close();
            break;
          default:
            break;
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
