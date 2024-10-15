export const getLanguagePrefix = (languageName: string) => {
  const [prefix] = languageName.split("-"); 
  return prefix;
};

