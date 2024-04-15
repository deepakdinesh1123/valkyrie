package execution

import (
	"fmt"
	"net/http"

	"github.com/deepakdinesh1123/valkyrie/executor/constants"
)

func Execute(w http.ResponseWriter, r *http.Request) {
	machinery_server := r.Context().Value(constants.ContextKey("machinery_server"))
	fmt.Printf("the server is machinery_server: %v\n", machinery_server)
	w.Write([]byte("Hello World"))
}
