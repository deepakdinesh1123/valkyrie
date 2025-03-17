import * as z from "zod";


export const TerminalReadSchema = z.object({
    "msgType": z.string().optional(),
    "terminalId": z.string(),
});
export type TerminalRead = z.infer<typeof TerminalReadSchema>;
