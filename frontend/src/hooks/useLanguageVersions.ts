import { useState, useEffect, useCallback } from 'react';
import { api } from '@/utils/api';
import { LanguageVersion } from '@/api-client';

const initLanguageVersion: LanguageVersion = {
  language_id: 1,
  version: "3.11",
  nix_package_name: "python311",
  template: "python/python.script.tmpl",
  search_query: "python311Packages",
  default_version: true,
};

export const useLanguageVersions = (languageId: number) => {
  const [languageVersions, setLanguageVersions] = useState<LanguageVersion[]>([]);
  const [selectedLanguageVersion, setSelectedLanguageVersion] = useState<LanguageVersion>(initLanguageVersion);

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