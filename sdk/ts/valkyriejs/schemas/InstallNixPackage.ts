import * as z from "zod";


export const InstallNixPackageSchema = z.object({
    "channel": z.string().optional(),
    "msgType": z.string().optional(),
    "pkgName": z.string(),
});
export type InstallNixPackage = z.infer<typeof InstallNixPackageSchema>;
