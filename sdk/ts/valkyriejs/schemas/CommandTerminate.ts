import * as z from "zod";


export const CommandTerminateSchema = z.object({
    "commandId": z.string(),
    "msgType": z.string().optional(),
});
export type CommandTerminate = z.infer<typeof CommandTerminateSchema>;
