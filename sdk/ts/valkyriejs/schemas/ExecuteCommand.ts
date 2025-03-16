import * as z from "zod";


export const ExecuteCommandSchema = z.object({
    "command": z.string(),
    "env": z.array(z.record(z.string(), z.string())).optional(),
    "msgType": z.string().optional(),
    "stderr": z.boolean().optional(),
    "stdin": z.boolean().optional(),
    "stdout": z.boolean().optional(),
    "workDir": z.string().optional(),
});
export type ExecuteCommand = z.infer<typeof ExecuteCommandSchema>;
