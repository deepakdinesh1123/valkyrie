import * as z from "zod";


export const UninstallNixPackageResponseSchema = z.object({
    "msg": z.string(),
    "msgType": z.string().optional(),
    "success": z.boolean(),
});
export type UninstallNixPackageResponse = z.infer<typeof UninstallNixPackageResponseSchema>;
