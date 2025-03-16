import * as z from "zod";


export const TerminalReadResponseSchema = z.object({
    "eof": z.boolean().optional(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "output": z.string(),
    "success": z.boolean(),
    "terminalId": z.string(),
});
export type TerminalReadResponse = z.infer<typeof TerminalReadResponseSchema>;
