import React, { createContext, useContext, useEffect, useState } from "react";
import { checkServerHealth } from "@/utils/checkServer";

const ServerStatusContext = createContext<{ isServerDown: boolean }>({
  isServerDown: false,
});

export const useServerStatus = () => useContext(ServerStatusContext);

export const ServerStatusProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [isServerDown, setIsServerDown] = useState(false);

  useEffect(() => {
    const checkStatus = async () => {
      const isUp = await checkServerHealth();
      setIsServerDown(!isUp);
    };

    checkStatus();
    const interval = setInterval(checkStatus, 10000);

    return () => clearInterval(interval); 
  }, []);

  return (
    <ServerStatusContext.Provider value={{ isServerDown }}>
      {children}
    </ServerStatusContext.Provider>
  );
};
