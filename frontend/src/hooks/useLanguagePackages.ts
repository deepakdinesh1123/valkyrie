import { useCallback, useEffect, useState } from 'react';
import { api } from '@/utils/api';

export const useLanguagePackages = ( selectedLanguage: string, searchString: string,) => {
  const [languagePackages, setLanguagePackages] = useState<{ name: string; version: string }[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchLanguagePackages = async () => {
      setLoading(true);
      setError(null);
      try {
        const languageParam = selectedLanguage.replace('-', '');
        const response = await api.searchLanguagePackages(searchString, languageParam);
        setLanguagePackages(response.data.packages);
      } catch (err) {
        console.error('Error fetching language packages:', err);
        setError('Failed to fetch language packages.');
      } finally {
        setLoading(false);
      }
    };

    if (searchString && selectedLanguage) {
      fetchLanguagePackages();
    }
  }, [searchString, selectedLanguage]);

  const resetLanguagePackages = useCallback(() => {
  setLanguagePackages([]);
}, []);

  return { languagePackages, loading, error, resetLanguagePackages };
};
