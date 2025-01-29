import * as z from "zod";


export const TerminalReadSchema = z.object({
    "msgType": z.string().optional(),
    "terminalId": z.string(),
    "timeout": z.union([z.number(), z.null()]).optional(),
});
export type TerminalRead = z.infer<typeof TerminalReadSchema>;
