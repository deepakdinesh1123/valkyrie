import { useEffect, useState } from 'react';
import { api } from '@/utils/api';

export const useSystemPackages = (searchString: string) => {
  const [systemPackages, setSystemPackages] = useState<{ name: string; version: string }[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSystemPackages = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await api.searchSystemPackages(searchString);
        setSystemPackages(response.data.packages);
      } catch (err) {
        console.error('Error fetching system packages:', err);
        setError('Failed to fetch system packages.');
      } finally {
        setLoading(false);
      }
    };

    if (searchString) {
      fetchSystemPackages();
    }
  }, [searchString]);

  return { systemPackages, loading, error };
};
