import React, { useCallback, useEffect, useState } from "react";
import LanguageSelector from "./components/LanguageSelector";


const App: React.FC = () => {


  return (
    <div className="flex h-screen overflow-hidden relative">
      <LanguageSelector>

      </LanguageSelector>
    </div>
  );
};

export default App;