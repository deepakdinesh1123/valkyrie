import { useState, useEffect } from 'react';
import { api } from '@/utils/api';
import { Language } from '@/api-client';

const initlanguage: Language = 
  {
    name: "python-3.10",
    extension: "py",
    monacolanguage: "python",
    defaultcode:
      '# Type your Python code here\n\ndef main():\n    pass\n\nif __name__ == "__main__":\n    main()',
    searchquery: "python310Packages"
  }


export const useLanguages = () => {
  const [languages, setLanguages] = useState<Language[]>([]);
  const [selectedLanguage, setSelectedLanguage] = useState<Language>(initlanguage);

  useEffect(() => {
    const fetchLanguages = async () => {
      try {
        const response = await api.getAllLanguages();
        const languageList = response.data.languages.map((lang) => ({
          name: lang.name,
          extension: lang.extension,
          monacolanguage: lang.monacolanguage,
          defaultcode: lang.defaultcode,
          searchquery: lang.searchquery,
        }));
        setLanguages(languageList);
        setSelectedLanguage(languageList[0]);

      } catch (error) {
        console.error('Failed to fetch languages:', error);
      }
    };

    fetchLanguages();
  }, []);

  return { languages, selectedLanguage, setSelectedLanguage };
};