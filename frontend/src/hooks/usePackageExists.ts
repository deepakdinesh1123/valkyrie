import { useCallback, useEffect, useState } from 'react';
import { api } from '@/utils/api'; 

export const usePackagesExist = (language: string, packages: string[]) => {
  const [existsResponse, setExistsResponse] = useState<{ exists: boolean; nonExistingPackages: string[] }>({
    exists: false,
    nonExistingPackages: [],
  });
  const [error, setError] = useState<string | null>(null);

  const checkPackagesExist = useCallback(async () => {
    setError(null);
    setExistsResponse({ exists: false, nonExistingPackages: [] }); 

    if (packages.length === 0) return; 

    try {
      const response = await api.packagesExist({
        language,
        packages,
      });

      setExistsResponse({
        exists: response.data.exists,
        nonExistingPackages: response.data.nonExistingPackages || [],
      });
    } catch (err) {
      console.error('Error checking package existence:', err);
      setError('Failed to check package existence.');
      setExistsResponse({ exists: false, nonExistingPackages: [] });
    }
  }, [language, packages]);

  useEffect(() => {
    if (language) {
      checkPackagesExist();
    }
  }, [language, checkPackagesExist]); 

  return { existsResponse, error, setExistsResponse };
};
