import React from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogClose } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface HelpModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const HelpModal: React.FC<HelpModalProps> = ({ isOpen, onClose }) => {
    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="bg-black p-6">
                <DialogHeader>
                    <DialogTitle className="text-white">Help</DialogTitle>
                    <DialogDescription className="text-white">
                        Odin is a code executor with ability to run various languages and supports addition of numerous packages or dependencies.
                    </DialogDescription>
                </DialogHeader>
                <p className="text-white mt-2">
                    To use this application, please follow these instructions:
                </p>
                <ul className="list-disc pl-5 text-white">
                    <li>Select a language from the dropdown menu.</li>
                    <li>Write your code in the editor.</li>
                    <li>Add any args to pass if necessary.</li>
                    <li>Add system and language dependencies if needed.</li>
                    <li>Click "Run code" to execute your code.</li>
                    <li>View the output in the terminal below.</li>
                </ul>
                <DialogClose asChild>
                    <Button className="mt-4 border border-transparent hover:border-white transition-colors">
                        Close
                    </Button>
                </DialogClose>
            </DialogContent>
        </Dialog>
    );
};

export default HelpModal;
