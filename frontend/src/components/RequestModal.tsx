import React from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogClose } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface RequestPackageModalProps {
    isOpen: boolean;
    onClose: () => void;
}

const RequestPackageModal: React.FC<RequestPackageModalProps> = ({ isOpen, onClose }) => {
    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="bg-black">
                <DialogHeader>
                    <DialogTitle className="text-white">Request a Package</DialogTitle>
                    <DialogDescription>
                        If you need a package that's not available, please submit a request to our team. We'll review it and add it to our system if possible. Currently, we are supporting only packages that are already available as{' '}
                        <a className="underline text-blue-500" href="https://nixos.org/manual/nixpkgs/stable/#overview-of-nixpkgs" target='_blank'>
                            nixpkgs
                        </a>
                        , so if you are not sure if your package is available, head over to{' '}
                        <a className="underline text-blue-500" href="https://search.nixos.org/packages" target='_blank'>
                            NixOS search
                        </a>{' '}
                        and search in the 24.11 channel.
                    </DialogDescription>
                </DialogHeader>
                <p className="text-white">
                    To request a package, please fill out this{' '}
                    <a className="underline text-blue-500" href="https://forms.gle/XpSVTpf3ix4rAjrr9" target='_blank'>
                        form
                    </a>{' '}
                    with the following information:
                </p>
                <ul className="list-disc pl-5 text-white">
                    <li>Package Name</li>
                    <li>Package Version</li>
                    <li>Type</li>
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

export default RequestPackageModal;