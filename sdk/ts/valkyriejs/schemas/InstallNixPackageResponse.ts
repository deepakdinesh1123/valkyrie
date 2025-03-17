import * as z from "zod";


export const InstallNixPackageResponseSchema = z.object({
    "msg": z.string(),
    "msgType": z.string().optional(),
    "success": z.boolean(),
});
export type InstallNixPackageResponse = z.infer<typeof InstallNixPackageResponseSchema>;
