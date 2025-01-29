import * as z from "zod";


export const TerminalWriteSchema = z.object({
    "input": z.string(),
    "msgType": z.string().optional(),
    "terminalId": z.string(),
});
export type TerminalWrite = z.infer<typeof TerminalWriteSchema>;
