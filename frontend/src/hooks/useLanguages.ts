import { useState, useEffect } from 'react';
import { api } from '@/utils/api';
import { LanguageResponse } from '@/api-client';

const initlanguage: LanguageResponse = 
  {
    id: 1,
    name: "python",
    extension: "py",
    monaco_language: "python",
  }

export const useLanguages = () => {
  const [languages, setLanguages] = useState<LanguageResponse[]>([]);
  const [selectedLanguage, setSelectedLanguage] = useState<LanguageResponse>(initlanguage);

  useEffect(() => {
    const fetchLanguages = async () => {
      try {
        const response = await api.getAllLanguages();
        const languageList = response.data.languages.map((lang) => ({
          id: lang.id,
          name: lang.name,
          extension: lang.extension,
          monaco_language: lang.monaco_language,
        }));
        setLanguages(languageList);
      } catch (error) {
        console.error('Failed to fetch languages:', error);
      }
    };

    fetchLanguages();
  }, []);

  return { languages, selectedLanguage, setSelectedLanguage };
};