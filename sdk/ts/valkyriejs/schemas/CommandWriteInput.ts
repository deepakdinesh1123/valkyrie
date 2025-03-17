import * as z from "zod";


export const CommandWriteInputSchema = z.object({
    "commandId": z.string(),
    "input": z.string().optional(),
    "msgType": z.string().optional(),
});
export type CommandWriteInput = z.infer<typeof CommandWriteInputSchema>;
