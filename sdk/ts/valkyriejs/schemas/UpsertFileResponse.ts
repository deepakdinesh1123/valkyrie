import * as z from "zod";


export const UpsertFileResponseSchema = z.object({
    "fileName": z.string(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "path": z.string(),
    "success": z.boolean(),
});
export type UpsertFileResponse = z.infer<typeof UpsertFileResponseSchema>;
