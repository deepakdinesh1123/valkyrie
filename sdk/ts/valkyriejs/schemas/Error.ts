import * as z from "zod";


export const PurpleErrorSchema = z.object({
    "message": z.string(),
});
export type PurpleError = z.infer<typeof PurpleErrorSchema>;
