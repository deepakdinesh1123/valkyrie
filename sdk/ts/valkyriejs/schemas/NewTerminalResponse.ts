import * as z from "zod";


export const NewTerminalResponseSchema = z.object({
    "msg": z.string(),
    "success": z.boolean(),
    "terminalID": z.string(),
});
export type NewTerminalResponse = z.infer<typeof NewTerminalResponseSchema>;
