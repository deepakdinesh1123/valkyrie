import React from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Language } from '@/api-client';

interface LanguageSelectorProps {
    languages: Language[];
    selectedLanguage: Language | null;
    onLanguageChange: (language: Language) => void;
}

export const LanguageSelector: React.FC<LanguageSelectorProps> = ({
    languages,
    selectedLanguage,
    onLanguageChange,
}) => (
    <div className="language-selector">
        <Select
            value={selectedLanguage?.name}
            onValueChange={(value) => {
                const language = languages.find((lang) => lang.name === value);
                if (language) {
                    onLanguageChange(language);
                }
            }}
        >
            <SelectTrigger className="w-[180px] bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 rounded-md px-2 py-1 transition-colors duration-200 ease-in-out">
                <SelectValue placeholder="Select a language" />
            </SelectTrigger>
            <SelectContent className="dark:bg-gray-800">
                {languages.map((lang) => (
                    <SelectItem key={lang.name} value={lang.name}>
                        {lang.name}
                    </SelectItem>
                ))}
            </SelectContent>
        </Select>
    </div>
);