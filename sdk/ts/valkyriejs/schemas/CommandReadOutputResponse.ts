import * as z from "zod";


export const CommandReadOutputResponseSchema = z.object({
    "commandId": z.string(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "stdout": z.string(),
    "success": z.boolean(),
});
export type CommandReadOutputResponse = z.infer<typeof CommandReadOutputResponseSchema>;
