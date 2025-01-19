import * as z from "zod";


export const ShellSchema = z.enum([
    "bash",
    "nix",
    "nix-shell",
    "sh",
]);
export type Shell = z.infer<typeof ShellSchema>;

export const TerminalSchema = z.object({
    "msgType": z.string().optional(),
    "nix_flake": z.string().optional(),
    "nix_shell": z.string().optional(),
    "packages": z.array(z.string()).optional(),
    "shell": ShellSchema,
});
export type Terminal = z.infer<typeof TerminalSchema>;
