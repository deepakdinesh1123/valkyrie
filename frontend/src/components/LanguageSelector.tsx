import * as React from "react";
import { Button } from "@/components/ui/button";
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import { useLanguages } from "@/hooks/useLanguages";
import { useLanguageVersions } from "@/hooks/useLanguageVersions";
import { useEffect, useState } from "react";
import { CheckIcon, ChevronsUpDown, Loader } from "lucide-react";
import { cn } from "@/lib/utils";
import { LanguageResponse, LanguageVersion } from "@/api-client";

interface LanguageSelectorProps {
    selectedLanguage: LanguageResponse | null;
    selectedLanguageVersion: LanguageVersion | null;
    onLanguageChange: (language: LanguageResponse) => void;
    onVersionChange: (version: LanguageVersion) => void;
}

export const LanguageSelector: React.FC<LanguageSelectorProps> = ({
    selectedLanguage,
    selectedLanguageVersion,
    onLanguageChange,
    onVersionChange,
}) => {
    const { languages, setSelectedLanguage, loading: loadingLanguages } = useLanguages();
    const { languageVersions, setSelectedLanguageVersion } = useLanguageVersions(selectedLanguage?.id || 0);

    const [languageOpen, setLanguageOpen] = useState(false);
    const [versionOpen, setVersionOpen] = useState(false);

    useEffect(() => {
        if (!selectedLanguage && languages.length > 0) {
            const initialLanguage = languages[0];
            setSelectedLanguage(initialLanguage);
            onLanguageChange(initialLanguage);
        }
    }, [languages, selectedLanguage, setSelectedLanguage, onLanguageChange]);

    const sortedLanguages = [...languages].sort((a, b) => a.name.localeCompare(b.name));

    const sortedVersions = [...languageVersions].sort((a, b) =>
        b.version.localeCompare(a.version, undefined, { numeric: true, sensitivity: "base" })
    );

    const handleLanguageChange = (language: LanguageResponse) => {
        setSelectedLanguage(language);
        setLanguageOpen(false);
        onLanguageChange(language);
    };

    const handleVersionChange = (version: LanguageVersion) => {
        setSelectedLanguageVersion(version);
        setVersionOpen(false);
        onVersionChange(version);
    };

    return (
        <div className="h-fit space-x-2">
            {/* Language Selector */}
            <Popover open={languageOpen} onOpenChange={setLanguageOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={languageOpen}
                        className="w-[180px] justify-between bg-neutral-900 text-white hover:bg-neutral-700 hover:text-white"
                    >
                        {loadingLanguages ? (
                            <Loader className="animate-spin h-4 w-4" />
                        ) : (
                            selectedLanguage?.name || "Select Language..."
                        )}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[180px] p-0 text-white bg-neutral-900">
                    <Command className="text-white bg-neutral-900">
                        {loadingLanguages ? (
                            <div className="flex items-center justify-center p-4">
                                <Loader className="animate-spin h-6 w-6" />
                            </div>
                        ) : (
                            <>
                                <CommandInput placeholder="Search Language" className="h-9" />
                                <CommandList>
                                    <CommandEmpty>No Language found.</CommandEmpty>
                                    <CommandGroup className="text-white bg-neutral-900">
                                        {sortedLanguages.map((language) => (
                                            <CommandItem
                                                key={language.id}
                                                value={language.name}
                                                onSelect={() => handleLanguageChange(language)}
                                            >
                                                {language.name}
                                                <CheckIcon
                                                    className={cn(
                                                        "ml-auto h-4 w-4",
                                                        selectedLanguage?.id === language.id ? "opacity-100" : "opacity-0"
                                                    )}
                                                />
                                            </CommandItem>
                                        ))}
                                    </CommandGroup>
                                </CommandList>
                            </>
                        )}
                    </Command>
                </PopoverContent>
            </Popover>

            {/* Version Selector */}
            <Popover open={versionOpen} onOpenChange={setVersionOpen}>
                <PopoverTrigger asChild>
                    <PopoverTrigger asChild>
                        <Button
                            variant="outline"
                            role="combobox"
                            aria-expanded={versionOpen}
                            className="w-[150px] justify-between bg-neutral-900 text-white hover:bg-neutral-700 hover:text-white"
                            title={selectedLanguageVersion?.version || "Select Version..."}
                        >
                                <span className="truncate">
                                    {selectedLanguageVersion?.version || "Select Version..."}    
                                </span>
                            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                        </Button>
                    </PopoverTrigger>

                </PopoverTrigger>
                <PopoverContent className="w-[150px] p-0 text-white bg-neutral-900">
                    <Command className="text-white bg-neutral-900">
                        <CommandInput placeholder="Search Version" className="h-9" />
                        <CommandList>
                            <CommandEmpty>No Version found.</CommandEmpty>
                            <CommandGroup className="text-white bg-neutral-900">
                                {sortedVersions.map((version) => (
                                    <CommandItem
                                        key={version.version}
                                        value={version.version}
                                        onSelect={() => {
                                            handleVersionChange(version);
                                        }}
                                    >
                                        {version.version}
                                        <CheckIcon
                                            className={cn(
                                                "ml-auto h-4 w-4",
                                                selectedLanguageVersion?.version === version.version ? "opacity-100" : "opacity-0"
                                            )}
                                        />
                                    </CommandItem>
                                ))}
                            </CommandGroup>
                        </CommandList>
                    </Command>
                </PopoverContent>
            </Popover>
        </div>
    );
};

export default LanguageSelector;