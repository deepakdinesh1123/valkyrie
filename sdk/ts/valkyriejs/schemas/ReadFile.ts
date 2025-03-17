import * as z from "zod";


export const ReadFileSchema = z.object({
    "msgType": z.string().optional(),
    "path": z.string(),
});
export type ReadFile = z.infer<typeof ReadFileSchema>;
