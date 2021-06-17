package api

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"syscall"
)

type response struct {
	Code   int
	Data   interface{} `json:"data,omitempty"`
	Errors apiErrors   `json:"errors,omitempty"`
}

func (r response) serveJSON(w http.ResponseWriter) {
	if r.Code == 0 {
		panic(errors.New("response status code not defined"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	if err := json.NewEncoder(w).Encode(r); err != nil {
		panic(err)
	}
}

func HandleAPIError(w http.ResponseWriter, err interface{}) {
	switch err := err.(type) {
	case apiErrors:
		st := http.StatusOK
		for _, e := range err {
			if e.Status > st {
				st = e.Status
			}
		}

		res := response{
			Code:   st,
			Errors: err,
		}
		res.serveJSON(w)
		logFatalErrors(err)

	case *apiError:
		res := response{
			Code:   err.Status,
			Errors: apiErrors{*err},
		}

		res.serveJSON(w)
		logFatalError(err)

	case *net.OpError:
		if err.Err == syscall.EPIPE || err.Err == syscall.ECONNRESET {
			break
		}

		res := response{
			Code: http.StatusInternalServerError,
			Errors: apiErrors{
				*newAPIError("Internal Server Error", errInternalServerNetPipe, err),
			},
		}
		res.serveJSON(w)
		logFatalError(&res.Errors[0])

	case error:
		res := response{
			Code: http.StatusInternalServerError,
			Errors: apiErrors{
				*newAPIError("Internal Server Error", errInternalServer, err),
			},
		}
		res.serveJSON(w)
		logFatalError(&res.Errors[0])

	case string:
		res := response{
			Code: http.StatusInternalServerError,
			Errors: apiErrors{
				*newAPIError("Internal Server Error", errInternalServer, errors.New(err)),
			},
		}
		res.serveJSON(w)
		logFatalError(&res.Errors[0])
	default:
		panic(err)
	}
}
