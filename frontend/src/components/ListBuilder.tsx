import React, { useState, useCallback, useEffect, useMemo } from "react";
import { Search, X } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import debounce from "lodash/debounce";

interface SearchableListBuilderProps {
  items: string[];
  onSelectionChange: (selectedItems: string[]) => void;
  onSearchChange: (searchTerm: string) => void;
  resetTrigger?: any;
}

const ListBuilder: React.FC<SearchableListBuilderProps> = ({
  items,
  onSelectionChange,
  onSearchChange,
  resetTrigger,
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

  // Reset effect
  useEffect(() => {
    if (resetTrigger !== undefined) {
      // Reset selected items and search term
      setSelectedItems([]);
      setSearchTerm("");
      // Notify parent of the changes
      onSelectionChange([]);
      onSearchChange("");
    }
  }, [resetTrigger, onSelectionChange, onSearchChange]);

  const filteredItems = useMemo(() => {
    return items.filter(
      (item) =>
        item.toLowerCase().includes(searchTerm.toLowerCase()) &&
        !selectedItems.includes(item)
    );
  }, [items, searchTerm, selectedItems]);

  const handleSelect = (item: string) => {
    const newSelectedItems = [...selectedItems, item];
    setSelectedItems(newSelectedItems);
    onSelectionChange(newSelectedItems);
  };

  const handleRemove = (item: string) => {
    const newSelectedItems = selectedItems.filter((i) => i !== item);
    setSelectedItems(newSelectedItems);
    onSelectionChange(newSelectedItems);
  };

  const handleSearchInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  return (
    <div className="w-full max-w-md">
      <div className="relative mb-4">
        <Input
          type="text"
          placeholder="Search items..."
          value={searchTerm}
          onChange={handleSearchInput}
          className="pl-10 pr-4 py-2 w-full"
        />
        <Search
          className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
          size={20}
        />
      </div>
      <div className="mb-4 flex flex-wrap gap-2">
        {selectedItems.map((item) => (
          <Badge
            key={item}
            variant="secondary"
            className="py-1 px-2 text-sm flex items-center gap-1"
          >
            {item}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleRemove(item)}
              className="h-5 w-5 p-0 hover:bg-red-100 rounded-full"
            >
              <X size={14} className="text-gray-500 hover:text-red-500" />
            </Button>
          </Badge>
        ))}
      </div>
      <div className="bg-white shadow-md rounded-md overflow-hidden">
        <ul className="max-h-40 overflow-y-auto">
          {filteredItems.map((item) => (
            <li
              key={item}
              className="px-4 py-2 hover:bg-gray-100 cursor-pointer"
              onClick={() => handleSelect(item)}
            >
              {item}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default ListBuilder;
