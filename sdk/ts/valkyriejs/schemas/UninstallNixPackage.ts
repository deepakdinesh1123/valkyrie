import * as z from "zod";


export const UninstallNixPackageSchema = z.object({
    "msgType": z.string().optional(),
    "pkgName": z.string(),
});
export type UninstallNixPackage = z.infer<typeof UninstallNixPackageSchema>;
