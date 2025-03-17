import * as z from "zod";


export const TerminalCloseSchema = z.object({
    "msgType": z.string().optional(),
    "terminalId": z.string(),
});
export type TerminalClose = z.infer<typeof TerminalCloseSchema>;
