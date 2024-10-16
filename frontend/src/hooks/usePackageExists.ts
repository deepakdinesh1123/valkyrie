import { useCallback, useEffect, useState } from 'react';
import { api } from '@/utils/api'; 

export const usePackagesExist = (language: string, packages: string[]) => {
  const [existsResponse, setExistsResponse] = useState<{ exists: boolean; nonExistingPackages: string[] }>({ exists: false, nonExistingPackages: [] });
  const [error, setError] = useState<string | null>(null);

  const checkPackagesExist = useCallback(async () => {
    setError(null);
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
    } finally {
    }
  }, [language, packages]);

  useEffect(() => {
    if (language && packages.length > 0) {
      checkPackagesExist();
    }
  }, [language, packages, checkPackagesExist]);

  return { existsResponse, error };
};
