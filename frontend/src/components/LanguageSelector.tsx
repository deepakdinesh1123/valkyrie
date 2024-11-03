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
import { useEffect } from "react";
import { CheckIcon, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { LanguageResponse, LanguageVersion } from "@/api-client";

interface LanguageSelectorProps {
    selectedLanguage: LanguageResponse;
    selectedLanguageVersion: LanguageVersion;
    onLanguageChange: (language: LanguageResponse) => void;
    onVersionChange: (version: LanguageVersion) => void;
    isLoading?: boolean;
}

export const LanguageSelector: React.FC<LanguageSelectorProps> = ({
    selectedLanguage,
    selectedLanguageVersion,
    onLanguageChange,
    onVersionChange,
}) => {
    const { languages, setSelectedLanguage } = useLanguages();
    const { languageVersions, refetch } = useLanguageVersions(selectedLanguage?.id);
    const [languageOpen, setLanguageOpen] = React.useState(false);
    const [versionOpen, setVersionOpen] = React.useState(false);

    const sortedLanguages = [...languages].sort((a, b) =>
        a.name.localeCompare(b.name)
    );

    const sortedVersions = [...languageVersions].sort((a, b) =>
        b.version.localeCompare(a.version, undefined, { numeric: true, sensitivity: 'base' })
    );

    useEffect(() => {
        if (sortedLanguages.length && !selectedLanguage) {
            const initialLanguage = sortedLanguages[0];
            setSelectedLanguage(initialLanguage);
            onLanguageChange(initialLanguage);
        }
    }, [sortedLanguages, selectedLanguage, setSelectedLanguage, onLanguageChange]);

    useEffect(() => {
        if (selectedLanguage?.id) {
            refetch().then(() => {
                if (sortedVersions.length) {
                    onVersionChange(sortedVersions[0]);
                }
            });
        }
    }, [selectedLanguage, refetch, sortedVersions, onVersionChange]);

    const handleLanguageChange = (language: LanguageResponse) => {
        setSelectedLanguage(language);
        setLanguageOpen(false);
        onLanguageChange(language);
    };

    return (
        <div className="h-fit space-x-2">
            <Popover open={languageOpen} onOpenChange={setLanguageOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={languageOpen}
                        className="w-[180px] justify-between bg-neutral-900 text-white hover:bg-neutral-700 hover:text-white"
                    >
                        {selectedLanguage?.name || "Select Language..."}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[180px] p-0 text-white bg-neutral-900">
                    <Command className="text-white bg-neutral-900">
                        <CommandInput placeholder="Search Language" className="h-9" />
                        <CommandList>
                            <CommandEmpty>No Language found.</CommandEmpty>
                            <CommandGroup className="text-white bg-neutral-900">
                                {sortedLanguages.map((language) => (
                                    <CommandItem
                                        key={language.name}
                                        value={language.name}
                                        onSelect={() => handleLanguageChange(language)}
                                    >
                                        {language.name}
                                        <CheckIcon
                                            className={cn(
                                                "ml-auto h-4 w-4",
                                                selectedLanguage?.name === language.name ? "opacity-100" : "opacity-0"
                                            )}
                                        />
                                    </CommandItem>
                                ))}
                            </CommandGroup>
                        </CommandList>
                    </Command>
                </PopoverContent>
            </Popover>

            <Popover open={versionOpen} onOpenChange={setVersionOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={versionOpen}
                        className="w-[150px] justify-between bg-neutral-900 text-white hover:bg-neutral-700 hover:text-white"
                    >
                        {selectedLanguageVersion?.version || "Select Version..."}
                        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
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
                                            setVersionOpen(false);
                                            onVersionChange(version);
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