import * as z from "zod";


export const MsgtypeSchema = z.enum([
    "TerminalWrite",
]);
export type Msgtype = z.infer<typeof MsgtypeSchema>;

export const TerminalwriteSchema = z.object({
    "content": z.string(),
    "msgType": MsgtypeSchema.optional(),
});
export type Terminalwrite = z.infer<typeof TerminalwriteSchema>;
