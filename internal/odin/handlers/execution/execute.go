package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/go-chi/chi/v5"

	"github.com/deepakdinesh1123/valkyrie/internal/constants"
	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/models/execution"
	VTasks "github.com/deepakdinesh1123/valkyrie/internal/tasks"
	"github.com/rs/zerolog"
)

func Execute(w http.ResponseWriter, r *http.Request) {
	environment := chi.URLParam(r, "environment")
	body, _ := io.ReadAll(r.Body)
	log := r.Context().Value(constants.ContextKey("logger")).(zerolog.Logger)
	task := tasks.Signature{
		Name: fmt.Sprintf("execute_%s", environment),
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: string(body),
			},
		},
		Headers: map[string]interface{}{
			"foo":    "bar",
			"map":    map[string]string{"cool": "nested header"},
			"struct": struct{ Foo string }{Foo: "foo struct"},
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	asyncResult, err := VTasks.MachineryServer.SendTaskWithContext(ctx, &task)
	if err != nil {
		log.Debug().Msg("")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task_uuid := asyncResult.GetState().TaskUUID
	result, _ := json.Marshal(&execution.ExecutionResult{
		ExecutionID: task_uuid,
	})
	w.Write([]byte(result))
}

func GetTaskState(w http.ResponseWriter, r *http.Request) {
	executionId := chi.URLParam(r, "executionId")
	taskState, err := VTasks.MachineryServer.GetBackend().GetState(executionId)
	if err != nil {
		logs.Logger.Debug().Msg("")
	}
	taskStateResponse, _ := json.Marshal(taskState)
	w.Write([]byte(taskStateResponse))
}

func GetTaskResult(w http.ResponseWriter, r *http.Request) {

}
