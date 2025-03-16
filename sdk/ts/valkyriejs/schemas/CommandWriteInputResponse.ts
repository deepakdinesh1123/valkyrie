import * as z from "zod";


export const CommandWriteInputResponseSchema = z.object({
    "commandId": z.string(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "success": z.boolean(),
});
export type CommandWriteInputResponse = z.infer<typeof CommandWriteInputResponseSchema>;
