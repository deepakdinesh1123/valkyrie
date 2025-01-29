import * as z from "zod";


export const TerminalWriteResponseSchema = z.object({
    "msg": z.string(),
    "success": z.boolean(),
    "terminalId": z.string(),
});
export type TerminalWriteResponse = z.infer<typeof TerminalWriteResponseSchema>;
