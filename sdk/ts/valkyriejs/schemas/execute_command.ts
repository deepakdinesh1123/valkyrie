import * as z from "zod";


export const EnvironmentVariableSchema = z.object({
    "key": z.string(),
    "value": z.string(),
});
export type EnvironmentVariable = z.infer<typeof EnvironmentVariableSchema>;

export const ExecuteCommandSchema = z.object({
    "command": z.string(),
    "env": z.array(EnvironmentVariableSchema).optional(),
    "msgType": z.string().optional(),
    "sandboxId": z.number(),
    "stderr": z.boolean().optional(),
    "stdin": z.boolean().optional(),
    "stdout": z.boolean().optional(),
    "workDir": z.string().optional(),
});
export type ExecuteCommand = z.infer<typeof ExecuteCommandSchema>;
