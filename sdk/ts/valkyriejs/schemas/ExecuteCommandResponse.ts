import * as z from "zod";


export const StateSchema = z.enum([
    "exited",
    "running",
    "starting",
    "stopped",
]);
export type State = z.infer<typeof StateSchema>;

export const ExecuteCommandResponseSchema = z.object({
    "commandId": z.string(),
    "msg": z.string(),
    "msgType": z.string().optional(),
    "state": StateSchema.optional(),
    "stdout": z.string(),
    "success": z.boolean(),
});
export type ExecuteCommandResponse = z.infer<typeof ExecuteCommandResponseSchema>;
