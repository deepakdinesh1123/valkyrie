import * as z from "zod";


export const AddFileSchema = z.object({
    "content": z.string(),
    "file_name": z.string(),
    "msgType": z.string().optional(),
    "path": z.string(),
    "sandboxId": z.number(),
});
export type AddFile = z.infer<typeof AddFileSchema>;
