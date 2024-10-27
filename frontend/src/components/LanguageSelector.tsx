import * as React from "react"
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from "@/components/ui/command"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import { useLanguages } from "@/hooks/useLanguages";
import { useLanguageVersions } from "@/hooks/useLanguageVersions";
import { useEffect } from "react"
import { CheckIcon } from "lucide-react"
import { cn } from "@/lib/utils"


const LanguageSelector = () => {
    const { languages, selectedLanguage, setSelectedLanguage } = useLanguages();
    const { languageVersions, refetch } = useLanguageVersions(selectedLanguage?.id);
    const [open, setOpen] = React.useState(false)
    const [value, setValue] = React.useState("")
    console.log(selectedLanguage);

    useEffect(() => {
        if (selectedLanguage?.id) {
            refetch();
        }
    }, [selectedLanguage, refetch]);

    return (
        <div>
            <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant="outline"
                        role="combobox"
                        aria-expanded={open}
                        className="w-[200px] justify-between"
                    >
                        {value
                            ? languages.find((language) => language.name === value)?.name
                            : "Select Language..."}
                    </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[200px] p-0">
                    <Command>
                        <CommandInput placeholder="Search Language" className="h-9" />
                        <CommandList>
                            <CommandEmpty>No Language found.</CommandEmpty>
                            <CommandGroup>
                                {languages.map((language) => (
                                    <CommandItem
                                        key={language.name}
                                        value={language.name}
                                        onSelect={(currentValue) => {
                                            setValue(currentValue === value ? "" : currentValue)
                                            setOpen(false)
                                            setSelectedLanguage(language)
                                        }}
                                    >
                                        {language.name}
                                        <CheckIcon
                                            className={cn(
                                                "ml-auto h-4 w-4",
                                                value === language.name ? "opacity-100" : "opacity-0"
                                            )}
                                        />
                                    </CommandItem>
                                ))}
                            </CommandGroup>
                        </CommandList>
                    </Command>
                </PopoverContent>
            </Popover>
            <Select>
                <SelectTrigger className="w-[180px]">
                    <SelectValue placeholder="Select version" />
                </SelectTrigger>
                <SelectContent>
                    <SelectGroup>
                        {languageVersions.map((ver) => (
                            <SelectItem
                                key={ver.version}
                                value={ver.version}>
                                {ver.version}
                            </SelectItem>
                        ))}
                    </SelectGroup>
                </SelectContent>
            </Select>
        </div>
    )
}

export default LanguageSelector;