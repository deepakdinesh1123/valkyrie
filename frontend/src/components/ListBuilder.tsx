import React, { useState, useCallback, useEffect, useMemo } from "react";
import { Search, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import debounce from "lodash/debounce";

interface SearchableListBuilderProps {
  items: { name: string; version: string }[];
  onSelectionChange: (selectedItems: string[]) => void;
  onSearchChange: (searchTerm: string) => void;
  resetTrigger?: any;
  nonExistingPackages?: string[];
}

const ListBuilder: React.FC<SearchableListBuilderProps> = ({
  items,
  onSelectionChange,
  onSearchChange,
  resetTrigger,
  nonExistingPackages = [],
}) => {
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [selectedItems, setSelectedItems] = useState<string[]>([]);

  const debouncedSearchChange = useCallback(
    debounce((value: string) => onSearchChange(value), 300),
    [onSearchChange]
  );

  useEffect(() => {
    debouncedSearchChange(searchTerm);
    return () => debouncedSearchChange.cancel();
  }, [searchTerm, debouncedSearchChange]);

  useEffect(() => {
    if (resetTrigger !== undefined) {
      setSelectedItems([]);
      setSearchTerm("");
      onSelectionChange([]);
      onSearchChange("");
    }
  }, [resetTrigger, onSelectionChange, onSearchChange]);

  useEffect(() => {
    if (nonExistingPackages.length > 0) {
      const updatedSelectedItems = selectedItems.filter(
        item => !nonExistingPackages.includes(item)
      );
      if (updatedSelectedItems.length !== selectedItems.length) {
        setSelectedItems(updatedSelectedItems);
        onSelectionChange(updatedSelectedItems);
        setSearchTerm("")
      }
    }
  }, [nonExistingPackages, selectedItems, onSelectionChange]);

  const filteredItems = useMemo(() => {
    return items.filter(
      (item) =>
        item.name.toLowerCase().includes(searchTerm.toLowerCase()) &&
        !selectedItems.includes(item.name)
    );
  }, [items, searchTerm, selectedItems]);

  const handleSelect = (item: { name: string; version: string }) => {
    const newSelectedItems = [...selectedItems, item.name];
    setSelectedItems(newSelectedItems);
    onSelectionChange(newSelectedItems);
  };

  const handleRemove = (itemName: string) => {
    const newSelectedItems = selectedItems.filter((i) => i !== itemName);
    setSelectedItems(newSelectedItems);
    onSelectionChange(newSelectedItems);
  };

  const handleSearchInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  return (
    <div className="w-full max-w-md">
      <div className="relative mb-4 ">
        <Input
          type="text"
          placeholder="Search items..."
          value={searchTerm}
          onChange={handleSearchInput}
          className="pl-10 pr-4 py-2 w-full outline-1 bg-neutral-900 text-white border-none"
        />
        <Search
          className="absolute outline-1 left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
          size={20}
        />
      </div>
      <div className="mb-4 flex flex-wrap gap-2">
        {selectedItems.map((itemName) => (
          <Badge
            key={itemName}
            variant="secondary"
            className="py-1 px-2 text-sm flex items-center gap-1 bg-neutral-900 text-white hover:bg-neutral-900"
          >
            {itemName}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleRemove(itemName)}
              className="h-5 w-5 p-0 hover:bg-red-100 hover:text-black rounded-full bg-neutral-800"
            >
              <X size={14} className="text-gray-500 hover:text-black-500" />
            </Button>
          </Badge>
        ))}
      </div>
      <div className="bg-neutral-900 text-white shadow-md rounded-md overflow-hidden">
        <ul className="max-h-40 overflow-y-auto">
          {filteredItems.map((item) => (
            <li
              key={item.name}
              className="px-4 py-2 hover:bg-gray-200 hover:text-black cursor-pointer"
              onClick={() => handleSelect(item)}
            >
              {`${item.name}`}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default ListBuilder;