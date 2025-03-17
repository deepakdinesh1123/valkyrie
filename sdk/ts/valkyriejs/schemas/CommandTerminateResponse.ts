import * as z from "zod";


export const CommandTerminateResponseSchema = z.object({
    "commandId": z.string(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "success": z.boolean(),
});
export type CommandTerminateResponse = z.infer<typeof CommandTerminateResponseSchema>;
