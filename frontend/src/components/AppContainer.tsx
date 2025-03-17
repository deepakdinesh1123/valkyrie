import React from "react";
import { useServerStatus } from "@/components/ServerStatusProvider";
import { AlertTriangle, RefreshCw } from 'lucide-react';

interface AppContainerProps {
  children: React.ReactNode;
}

const AppContainer: React.FC<AppContainerProps> = ({ children }) => {
  const { isServerDown } = useServerStatus();

  if (isServerDown) {
    return (
        <div className="flex items-center justify-center min-h-screen bg-neutral-950 text-center overflow-hidden">
        <div className="relative px-6 py-16 max-w-lg mx-auto">
          {/* Animated error icon with pulsing effect */}
          <div className="animate-pulse-slow mb-8">
            <AlertTriangle 
              className="mx-auto text-neutral-500 opacity-80" 
              size={180} 
              strokeWidth={1.5}
            />
          </div>
  
          {/* Error message with dynamic styling */}
          <div className="space-y-4 relative z-10">
            <h2 className="text-4xl font-bold text-white drop-shadow-lg animate-fade-in">
              Server Temporarily Unavailable
            </h2>
            <p className="text-neutral-300 max-w-md mx-auto animate-slide-up">
              We're experiencing some technical difficulties. Our team is working 
              hard to bring the service back online as quickly as possible.
            </p>
  
            {/* Retry button with hover effect */}
            <div className="mt-8">
              <button 
                onClick={() => window.location.reload()}
                className="group flex items-center justify-center mx-auto space-x-2 px-6 py-3 
                           bg-neutral-800 hover:bg-neutral-700 text-white 
                           rounded-full transition-all duration-300 
                           transform hover:scale-105 hover:shadow-lg"
              >
                <RefreshCw 
                  className="mr-2 group-hover:animate-spin" 
                  size={20} 
                />
                Retry Connection
              </button>
            </div>
          </div>
  
          {/* Subtle background decorative elements */}
          <div className="absolute inset-0 pointer-events-none">
            <div className="absolute top-0 right-0 w-72 h-72 bg-neutral-900/20 rounded-full blur-3xl"></div>
            <div className="absolute bottom-0 left-0 w-64 h-64 bg-neutral-900/20 rounded-full blur-3xl"></div>
          </div>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};

export default AppContainer;
