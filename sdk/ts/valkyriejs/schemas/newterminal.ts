import * as z from "zod";


export const MsgtypeSchema = z.enum([
    "NewTerminal",
]);
export type Msgtype = z.infer<typeof MsgtypeSchema>;

// Shell type to use

export const ShellSchema = z.enum([
    "bash",
    "nix",
    "nix-shell",
    "sh",
]);
export type Shell = z.infer<typeof ShellSchema>;

export const NewterminalSchema = z.object({
    "msgType": MsgtypeSchema.optional(),
    "nix_flake": z.string().optional(),
    "nix_shell": z.string().optional(),
    "packages": z.array(z.string()).optional(),
    "shell": ShellSchema,
});
export type Newterminal = z.infer<typeof NewterminalSchema>;
