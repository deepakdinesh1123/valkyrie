import * as z from "zod";


export const CommandReadOutputSchema = z.object({
    "commandId": z.string(),
    "msgType": z.string().optional(),
});
export type CommandReadOutput = z.infer<typeof CommandReadOutputSchema>;
