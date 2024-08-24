// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

func encodeCancelJobResponse(response CancelJobRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *CancelJobOK:
		w.WriteHeader(200)

		return nil

	case *CancelJobBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *CancelJobInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeDeleteJobResponse(response DeleteJobRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *DeleteJobOK:
		w.WriteHeader(200)

		return nil

	case *DeleteJobBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *DeleteJobNotFound:
		w.WriteHeader(404)

		return nil

	case *DeleteJobInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeExecuteResponse(response ExecuteRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *ExecuteOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ExecuteBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *ExecuteInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetAllExecutionResultsResponse(response GetAllExecutionResultsRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetAllExecutionResultsOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetAllExecutionResultsBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetAllExecutionResultsInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetAllExecutionsResponse(response GetAllExecutionsRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetAllExecutionsOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetAllExecutionsBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetAllExecutionsInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetExecutionConfigResponse(response GetExecutionConfigRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *ExecutionConfig:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *Error:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetExecutionResultsByIdResponse(response GetExecutionResultsByIdRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetExecutionResultsByIdOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetExecutionResultsByIdBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetExecutionResultsByIdNotFound:
		w.WriteHeader(404)

		return nil

	case *GetExecutionResultsByIdInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetExecutionWorkersResponse(response GetExecutionWorkersRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetExecutionWorkersOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetExecutionWorkersBadRequest:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(400)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *GetExecutionWorkersInternalServerError:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}

func encodeGetVersionResponse(response GetVersionRes, w http.ResponseWriter) error {
	switch response := response.(type) {
	case *GetVersionOK:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	case *Error:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(500)

		e := new(jx.Encoder)
		response.Encode(e)
		if _, err := e.WriteTo(w); err != nil {
			return errors.Wrap(err, "write")
		}

		return nil

	default:
		return errors.Errorf("unexpected response type: %T", response)
	}
}
