import * as z from "zod";


export const UpsertFileSchema = z.object({
    "content": z.string(),
    "fileName": z.string(),
    "msgType": z.string().optional(),
    "patch": z.string(),
    "path": z.string(),
});
export type UpsertFile = z.infer<typeof UpsertFileSchema>;
