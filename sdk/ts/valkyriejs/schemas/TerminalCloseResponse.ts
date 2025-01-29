import * as z from "zod";


export const TerminalCloseResponseSchema = z.object({
    "msg": z.string(),
    "success": z.boolean(),
    "terminalId": z.string(),
});
export type TerminalCloseResponse = z.infer<typeof TerminalCloseResponseSchema>;
