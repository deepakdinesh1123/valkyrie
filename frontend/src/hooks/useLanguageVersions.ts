import { useState, useEffect, useCallback } from 'react';
import { api } from '@/utils/api';
import { LanguageVersion } from '@/api-client';

export const useLanguageVersions = (languageId: number) => {
  const [languageVersions, setLanguageVersions] = useState<LanguageVersion[]>([]);
  const [selectedLanguageVersion, setSelectedLanguageVersion] = useState<LanguageVersion | null>(null);

  const fetchLanguageVersions = useCallback(async () => {
    if (!languageId) return; 

    try {
      const response = await api.getAllVersions(languageId);
      const LanguageVersionList = response.data.languageVersions.map((lang) => ({
        language_id: lang.language_id,
        version: lang.version,
        nix_package_name: lang.nix_package_name,
        template: lang.template,
        search_query: lang.search_query,
        default_version: lang.default_version,
      }));
      setLanguageVersions(LanguageVersionList);
      setSelectedLanguageVersion(LanguageVersionList[0]);

    } catch (error) {
      console.error('Failed to fetch LanguageVersions:', error);
    }
  }, [languageId]); 

  useEffect(() => {
    fetchLanguageVersions();
  }, [fetchLanguageVersions]);

  return { languageVersions, selectedLanguageVersion, setSelectedLanguageVersion, refetch: fetchLanguageVersions };
};