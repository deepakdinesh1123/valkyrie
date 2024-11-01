"use client";

import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
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
}) => {
    const [open, setOpen] = React.useState(false);
    const [query, setQuery] = React.useState("");

    const sortedLanguages = [...languages].sort((a, b) =>
        a.name.localeCompare(b.name)
    );

    const filteredLanguages =
        query === ""
            ? sortedLanguages
            : sortedLanguages.filter((lang) =>
                  lang.name.toLowerCase().includes(query.toLowerCase())
              );

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className="w-[200px] justify-between bg-neutral-900 text-white hover:bg-neutral-900 hover:text-white"
                >
                    {selectedLanguage ? selectedLanguage.name : "Select language..."}
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50 " />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[200px] p-0">
                <Command className="bg-neutral-900 text-white">
                    <CommandInput 
                        placeholder="Search language..."
                        onValueChange={setQuery}
                    />
                    <CommandList>
                        <CommandEmpty>No languages found.</CommandEmpty>
                        <CommandGroup className="bg-neutral-900 text-white">
                            {filteredLanguages.map((lang) => (
                                <CommandItem
                                    key={lang.name}
                                    value={lang.name}
                                    onSelect={(currentValue:string) => {
                                        const language = filteredLanguages.find(l => l.name === currentValue);
                                        if (language) {
                                            onLanguageChange(language);
                                        }
                                        setOpen(false);
                                    }}
                                >
                                    <Check
                                        className={cn(
                                            "mr-2 h-4 w-4",
                                            selectedLanguage?.name === lang.name ? "opacity-100" : "opacity-0"
                                        )}
                                    />
                                    {lang.name}
                                </CommandItem>
                            ))}
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
    );
};
