import * as z from "zod";


export const DeleteFileSchema = z.object({
    "msgType": z.string().optional(),
    "path": z.string(),
});
export type DeleteFile = z.infer<typeof DeleteFileSchema>;
