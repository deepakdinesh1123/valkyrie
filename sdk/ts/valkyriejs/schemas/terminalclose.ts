import * as z from "zod";


export const MsgtypeSchema = z.enum([
    "TerminalClose",
]);
export type Msgtype = z.infer<typeof MsgtypeSchema>;

export const TerminalcloseSchema = z.object({
    "msgType": MsgtypeSchema.optional(),
});
export type Terminalclose = z.infer<typeof TerminalcloseSchema>;
