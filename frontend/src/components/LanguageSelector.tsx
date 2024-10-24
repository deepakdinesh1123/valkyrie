import React, { useEffect, useRef } from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Loader2 } from "lucide-react";
import { Language } from '@/api-client';

interface LanguageSelectorProps {
    languages: Language[];
    selectedLanguage: Language | null;
    onLanguageChange: (language: Language) => void;
    isLoading?: boolean;
}

export const LanguageSelector: React.FC<LanguageSelectorProps> = ({
    languages,
    selectedLanguage,
    onLanguageChange,
    isLoading = false
}) => {
    const isInitialLoad = useRef(true);

    const sortedLanguages = [...languages].sort((a, b) =>
        a.name.localeCompare(b.name)
    );

    useEffect(() => {
        if (isInitialLoad.current && sortedLanguages.length > 0) {
            onLanguageChange(sortedLanguages[0]);
            isInitialLoad.current = false;
        }
    }, [sortedLanguages, onLanguageChange]);

    if (isLoading) {
        return (
            <div className="w-[180px] h-9 flex items-center justify-center bg-neutral-900 rounded-md">
                <Loader2 className="h-4 w-4 text-white animate-spin" />
            </div>
        );
    }

    return (
        <div className="language-selector focus:outline-none">
            <Select
                value={selectedLanguage?.name}
                onValueChange={(value) => {
                    const language = sortedLanguages.find((lang) => lang.name === value);
                    if (language) {
                        onLanguageChange(language);
                    }
                }}
            >
                <SelectTrigger className="w-[180px] outline-none bg-neutral-900 text-white rounded-md px-2 py-1 transition-colors duration-200 ease-in-out">
                    <SelectValue placeholder="Select a language" />
                </SelectTrigger>
                <SelectContent className="bg-black text-white">
                    {sortedLanguages.map((lang) => (
                        <SelectItem key={lang.name} value={lang.name}>
                            {lang.name}
                        </SelectItem>
                    ))}
                </SelectContent>
            </Select>
        </div>
    );
};