import { useCallback, useEffect, useState } from 'react';
import { api } from '@/utils/api';

export const useDefaultPackages = (selectedLanguage: string) => {
    const [defaultSystemPackages, setDefaultSystemPackages] = useState<{ name: string; version: string }[]>([]);
    const [defaultLanguagePackages, setDefaultLanguagePackages] = useState<{ name: string; version: string }[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchSystemPackages = async () => {
            setLoading(true);
            setError(null);
            try {
                const response = await api.fetchSystemPackages();
                setDefaultSystemPackages(response.data.packages);
            } catch (err) {
                console.error('Error fetching system packages:', err);
                setError('Failed to fetch system packages.');
            } finally {
                setLoading(false);
            }
        };
        fetchSystemPackages();
    }, []);
    
    const fetchLanguagePackages = useCallback(async () => {
        if (!selectedLanguage) return;
        try {
            const response = await api.fetchLanguagePackages(selectedLanguage);
            setDefaultLanguagePackages(response.data.packages);

        } catch (error) {
            console.error('Failed to fetch LanguagePackages:', error);
        }
    }, [selectedLanguage]);
    fetchLanguagePackages();

    useEffect(() => {
        fetchLanguagePackages();
    }, [fetchLanguagePackages]);

    return { defaultSystemPackages, defaultLanguagePackages, loading, error }
};
