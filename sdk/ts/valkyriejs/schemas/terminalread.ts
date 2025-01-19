import * as z from "zod";


export const MsgtypeSchema = z.enum([
    "TerminalRead",
]);
export type Msgtype = z.infer<typeof MsgtypeSchema>;

export const TerminalreadSchema = z.object({
    "msgType": MsgtypeSchema.optional(),
    "timeout": z.union([z.number(), z.null()]).optional(),
});
export type Terminalread = z.infer<typeof TerminalreadSchema>;
