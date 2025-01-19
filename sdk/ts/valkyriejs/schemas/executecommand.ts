import * as z from "zod";


export const MsgtypeSchema = z.enum([
    "ExecuteCommand",
]);
export type Msgtype = z.infer<typeof MsgtypeSchema>;

export const EnvironmentVariableSchema = z.object({
    "key": z.string(),
    "value": z.string(),
});
export type EnvironmentVariable = z.infer<typeof EnvironmentVariableSchema>;

export const ExecutecommandSchema = z.object({
    "command": z.string(),
    "env": z.array(EnvironmentVariableSchema).optional(),
    "msgType": MsgtypeSchema.optional(),
    "sandboxId": z.number(),
    "stderr": z.boolean().optional(),
    "stdin": z.boolean().optional(),
    "stdout": z.boolean().optional(),
    "workDir": z.string().optional(),
});
export type Executecommand = z.infer<typeof ExecutecommandSchema>;
