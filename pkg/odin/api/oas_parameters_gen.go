// Code generated by ogen, DO NOT EDIT.

package api

import (
	"net/http"
	"net/url"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// GetExecutionResultParams is parameters of getExecutionResult operation.
type GetExecutionResultParams struct {
	ExecutionId int64
}

func unpackGetExecutionResultParams(packed middleware.Parameters) (params GetExecutionResultParams) {
	{
		key := middleware.ParameterKey{
			Name: "executionId",
			In:   "path",
		}
		params.ExecutionId = packed[key].(int64)
	}
	return params
}

func decodeGetExecutionResultParams(args [1]string, argsEscaped bool, r *http.Request) (params GetExecutionResultParams, _ error) {
	// Decode path: executionId.
	if err := func() error {
		param := args[0]
		if argsEscaped {
			unescaped, err := url.PathUnescape(args[0])
			if err != nil {
				return errors.Wrap(err, "unescape path")
			}
			param = unescaped
		}
		if len(param) > 0 {
			d := uri.NewPathDecoder(uri.PathDecoderConfig{
				Param:   "executionId",
				Value:   param,
				Style:   uri.PathStyleSimple,
				Explode: false,
			})

			if err := func() error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToInt64(val)
				if err != nil {
					return err
				}

				params.ExecutionId = c
				return nil
			}(); err != nil {
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "executionId",
			In:   "path",
			Err:  err,
		}
	}
	return params, nil
}
