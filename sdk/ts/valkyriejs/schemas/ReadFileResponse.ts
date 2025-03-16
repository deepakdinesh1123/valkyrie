import * as z from "zod";


export const ReadFileResponseSchema = z.object({
    "content": z.string(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "path": z.string(),
    "success": z.boolean(),
});
export type ReadFileResponse = z.infer<typeof ReadFileResponseSchema>;
