import { useState, useEffect, useCallback } from 'react';
import { api } from '@/utils/api';
import { LanguageVersion } from '@/api-client';

const initLanguageVersion: LanguageVersion = {
  language_id: 1,
  version: "3.10",
  nix_package_name: "python310",
  flake_template: "python/python.flake.tmpl",
  script_template: "python/python.script.tmpl",
  search_query: "python310Packages",
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
        flake_template: lang.flake_template,
        script_template: lang.script_template,
        search_query: lang.search_query,
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