import * as z from "zod";


export const MsgtypeSchema = z.enum([
    "AddFile",
]);
export type Msgtype = z.infer<typeof MsgtypeSchema>;

export const AddfileSchema = z.object({
    "content": z.string(),
    "file_name": z.string(),
    "msgType": MsgtypeSchema.optional(),
    "path": z.string(),
    "sandboxId": z.number(),
});
export type Addfile = z.infer<typeof AddfileSchema>;
