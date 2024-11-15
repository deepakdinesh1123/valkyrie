// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"

	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
)

func recordError(string, error) {}

// handleCancelJobRequest handles cancelJob operation.
//
// Cancel Job.
//
// PUT /executions/{JobId}/
func (s *Server) handleCancelJobRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "CancelJob",
			ID:   "cancelJob",
		}
	)
	params, err := decodeCancelJobParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response CancelJobRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "CancelJob",
			OperationSummary: "Cancel Job",
			OperationID:      "cancelJob",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "JobId",
					In:   "path",
				}: params.JobId,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = CancelJobParams
			Response = CancelJobRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackCancelJobParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.CancelJob(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.CancelJob(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeCancelJobResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleDeleteJobRequest handles deleteJob operation.
//
// Delete job.
//
// DELETE /executions/{JobId}/
func (s *Server) handleDeleteJobRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "DeleteJob",
			ID:   "deleteJob",
		}
	)
	params, err := decodeDeleteJobParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response DeleteJobRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "DeleteJob",
			OperationSummary: "Delete job",
			OperationID:      "deleteJob",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "JobId",
					In:   "path",
				}: params.JobId,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = DeleteJobParams
			Response = DeleteJobRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackDeleteJobParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.DeleteJob(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.DeleteJob(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeDeleteJobResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleExecuteRequest handles execute operation.
//
// Execute a script.
//
// POST /executions/execute/
func (s *Server) handleExecuteRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "Execute",
			ID:   "execute",
		}
	)
	request, close, err := s.decodeExecuteRequest(r)
	if err != nil {
		err = &ogenerrors.DecodeRequestError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeRequest", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}
	defer func() {
		if err := close(); err != nil {
			recordError("CloseRequest", err)
		}
	}()

	var response ExecuteRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "Execute",
			OperationSummary: "Execute a script",
			OperationID:      "execute",
			Body:             request,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = *ExecutionRequest
			Params   = struct{}
			Response = ExecuteRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.Execute(ctx, request)
				return response, err
			},
		)
	} else {
		response, err = s.h.Execute(ctx, request)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeExecuteResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetAllExecutionResultsRequest handles getAllExecutionResults operation.
//
// Get all execution results.
//
// GET /executions/results/
func (s *Server) handleGetAllExecutionResultsRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetAllExecutionResults",
			ID:   "getAllExecutionResults",
		}
	)
	params, err := decodeGetAllExecutionResultsParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetAllExecutionResultsRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetAllExecutionResults",
			OperationSummary: "Get all execution results",
			OperationID:      "getAllExecutionResults",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "page",
					In:   "query",
				}: params.Page,
				{
					Name: "pageSize",
					In:   "query",
				}: params.PageSize,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetAllExecutionResultsParams
			Response = GetAllExecutionResultsRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetAllExecutionResultsParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetAllExecutionResults(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetAllExecutionResults(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetAllExecutionResultsResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetAllExecutionsRequest handles getAllExecutions operation.
//
// Get all executions.
//
// GET /executions/
func (s *Server) handleGetAllExecutionsRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetAllExecutions",
			ID:   "getAllExecutions",
		}
	)
	params, err := decodeGetAllExecutionsParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetAllExecutionsRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetAllExecutions",
			OperationSummary: "Get all executions",
			OperationID:      "getAllExecutions",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "page",
					In:   "query",
				}: params.Page,
				{
					Name: "pageSize",
					In:   "query",
				}: params.PageSize,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetAllExecutionsParams
			Response = GetAllExecutionsRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetAllExecutionsParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetAllExecutions(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetAllExecutions(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetAllExecutionsResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetExecutionConfigRequest handles getExecutionConfig operation.
//
// Get execution config.
//
// GET /execution/config/
func (s *Server) handleGetExecutionConfigRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err error
	)

	var response GetExecutionConfigRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetExecutionConfig",
			OperationSummary: "Get execution config",
			OperationID:      "getExecutionConfig",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = GetExecutionConfigRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetExecutionConfig(ctx)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetExecutionConfig(ctx)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetExecutionConfigResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetExecutionResultsByIdRequest handles getExecutionResultsById operation.
//
// Get execution result.
//
// GET /executions/{JobId}/
func (s *Server) handleGetExecutionResultsByIdRequest(args [1]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetExecutionResultsById",
			ID:   "getExecutionResultsById",
		}
	)
	params, err := decodeGetExecutionResultsByIdParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetExecutionResultsByIdRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetExecutionResultsById",
			OperationSummary: "Get execution result",
			OperationID:      "getExecutionResultsById",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "JobId",
					In:   "path",
				}: params.JobId,
				{
					Name: "page",
					In:   "query",
				}: params.Page,
				{
					Name: "pageSize",
					In:   "query",
				}: params.PageSize,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetExecutionResultsByIdParams
			Response = GetExecutionResultsByIdRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetExecutionResultsByIdParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetExecutionResultsById(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetExecutionResultsById(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetExecutionResultsByIdResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetExecutionWorkersRequest handles getExecutionWorkers operation.
//
// Get all execution workers.
//
// GET /executions/workers
func (s *Server) handleGetExecutionWorkersRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err          error
		opErrContext = ogenerrors.OperationContext{
			Name: "GetExecutionWorkers",
			ID:   "getExecutionWorkers",
		}
	)
	params, err := decodeGetExecutionWorkersParams(args, argsEscaped, r)
	if err != nil {
		err = &ogenerrors.DecodeParamsError{
			OperationContext: opErrContext,
			Err:              err,
		}
		defer recordError("DecodeParams", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	var response GetExecutionWorkersRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetExecutionWorkers",
			OperationSummary: "Get all execution workers",
			OperationID:      "getExecutionWorkers",
			Body:             nil,
			Params: middleware.Parameters{
				{
					Name: "page",
					In:   "query",
				}: params.Page,
				{
					Name: "pageSize",
					In:   "query",
				}: params.PageSize,
			},
			Raw: r,
		}

		type (
			Request  = struct{}
			Params   = GetExecutionWorkersParams
			Response = GetExecutionWorkersRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			unpackGetExecutionWorkersParams,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetExecutionWorkers(ctx, params)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetExecutionWorkers(ctx, params)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetExecutionWorkersResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}

// handleGetVersionRequest handles getVersion operation.
//
// Get version.
//
// GET /version/
func (s *Server) handleGetVersionRequest(args [0]string, argsEscaped bool, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var (
		err error
	)

	var response GetVersionRes
	if m := s.cfg.Middleware; m != nil {
		mreq := middleware.Request{
			Context:          ctx,
			OperationName:    "GetVersion",
			OperationSummary: "Get version",
			OperationID:      "getVersion",
			Body:             nil,
			Params:           middleware.Parameters{},
			Raw:              r,
		}

		type (
			Request  = struct{}
			Params   = struct{}
			Response = GetVersionRes
		)
		response, err = middleware.HookMiddleware[
			Request,
			Params,
			Response,
		](
			m,
			mreq,
			nil,
			func(ctx context.Context, request Request, params Params) (response Response, err error) {
				response, err = s.h.GetVersion(ctx)
				return response, err
			},
		)
	} else {
		response, err = s.h.GetVersion(ctx)
	}
	if err != nil {
		defer recordError("Internal", err)
		s.cfg.ErrorHandler(ctx, w, r, err)
		return
	}

	if err := encodeGetVersionResponse(response, w); err != nil {
		defer recordError("EncodeResponse", err)
		if !errors.Is(err, ht.ErrInternalServerErrorResponse) {
			s.cfg.ErrorHandler(ctx, w, r, err)
		}
		return
	}
}
