import * as z from "zod";


export const DeleteFileResponseSchema = z.object({
    "msg": z.string(),
    "msgType": z.string().optional(),
    "path": z.string(),
    "success": z.boolean(),
});
export type DeleteFileResponse = z.infer<typeof DeleteFileResponseSchema>;
