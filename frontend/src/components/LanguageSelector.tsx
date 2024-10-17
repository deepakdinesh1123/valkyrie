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
            <SelectTrigger className="w-[180px] outline-none bg-neutral-900 text-white rounded-md px-2 py-1 transition-colors duration-200 ease-in-out">
                <SelectValue placeholder="Select a language" />
            </SelectTrigger>
            <SelectContent className="bg-black text-white">
                {languages.map((lang) => (
                    <SelectItem key={lang.name} value={lang.name}>
                        {lang.name}
                    </SelectItem>
                ))}
            </SelectContent>
        </Select>
    </div>
);