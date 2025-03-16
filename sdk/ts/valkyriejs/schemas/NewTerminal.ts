import * as z from "zod";


export const NewTerminalSchema = z.object({
    "env": z.union([z.record(z.string(), z.string()), z.null()]).optional(),
    "msgType": z.string().optional(),
    "nixFlake": z.union([z.null(), z.string()]).optional(),
    "nixShell": z.union([z.null(), z.string()]).optional(),
    "packages": z.union([z.array(z.string()), z.null()]).optional(),
});
export type NewTerminal = z.infer<typeof NewTerminalSchema>;
