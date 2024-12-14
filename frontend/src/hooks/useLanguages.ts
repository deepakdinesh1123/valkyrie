import { useState, useEffect } from 'react';
import { api } from '@/utils/api';
import { LanguageResponse } from '@/api-client';

export const useLanguages = () => {
  const [languages, setLanguages] = useState<LanguageResponse[]>([]);
  const [selectedLanguage, setSelectedLanguage] = useState<LanguageResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchLanguages = async () => {
      try {
        const response = await api.getAllLanguages();
        const languageList = response.data.languages.map((lang) => ({
          id: lang.id,
          name: lang.name,
          extension: lang.extension,
          monaco_language: lang.monaco_language,
          default_code: lang.default_code,
        }));
        setLanguages(languageList);

        if (languageList.length > 0) {
          setSelectedLanguage(languageList[0]);
        }
      } catch (error) {
        console.error('Failed to fetch languages:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchLanguages();
  }, []);

  return { languages, selectedLanguage, setSelectedLanguage, loading };
};
