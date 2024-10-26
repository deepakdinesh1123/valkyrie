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
                    <li>Add system and language dependencies if needed. (The versions of the packages are what nixos-24.05 provides.)</li>
                    <li>Click "Run code" to execute your code.</li>
                    <li>View the output in the terminal below.</li>
                </ul>
                <span className='text-white'>Community links:</span>
                <a
                    href="https://discord.gg/3cJpQNgT"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="hover:opacity-80 transition-opacity"
                >

                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                        className="w-5 h-5 text-white"
                    >
                        <path
                            d="M20.317 4.369a19.791 19.791 0 00-4.885-1.527.074.074 0 00-.079.037c-.21.375-.444.864-.608 1.248a18.292 18.292 0 00-5.487 0 12.327 12.327 0 00-.617-1.248.079.079 0 00-.079-.037A19.425 19.425 0 003.68 4.369a.07.07 0 00-.032.027C.533 9.39-.32 14.313.099 19.163a.082.082 0 00.031.058 19.875 19.875 0 005.996 3.03.079.079 0 00.084-.027c.464-.637.873-1.312 1.226-2.016a.074.074 0 00-.041-.105 13.12 13.12 0 01-1.872-.9.076.076 0 01-.008-.126c.125-.094.25-.191.371-.292a.073.073 0 01.077-.01c3.927 1.793 8.18 1.793 12.061 0a.073.073 0 01.079.009c.122.1.247.198.372.292a.076.076 0 01-.007.125 12.663 12.663 0 01-1.873.901.075.075 0 00-.04.105c.366.704.776 1.379 1.224 2.016a.079.079 0 00.084.028 19.875 19.875 0 005.997-3.03.08.08 0 00.031-.058c.5-5.192-.83-10.058-3.575-14.767a.061.061 0 00-.03-.028zM8.02 15.331c-1.182 0-2.158-1.085-2.158-2.419 0-1.333.953-2.418 2.158-2.418 1.21 0 2.174 1.09 2.158 2.418 0 1.334-.953 2.419-2.158 2.419zm7.974 0c-1.182 0-2.158-1.085-2.158-2.419 0-1.333.953-2.418 2.158-2.418 1.21 0 2.174 1.09 2.158 2.418 0 1.334-.953 2.419-2.158 2.419z"
                        />
                    </svg>
                </a>
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
