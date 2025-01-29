import * as z from "zod";

// Shell type to use

export const ShellSchema = z.enum([
    "bash",
    "nix",
    "nix-shell",
    "sh",
]);
export type Shell = z.infer<typeof ShellSchema>;

export const NewTerminalSchema = z.object({
    "msgType": z.string().optional(),
    "nixFlake": z.string().optional(),
    "nixShell": z.string().optional(),
    "packages": z.array(z.string()).optional(),
    "shell": ShellSchema,
});
export type NewTerminal = z.infer<typeof NewTerminalSchema>;
