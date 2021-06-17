package api

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/getsentry/raven-go"
	"github.com/gofrs/uuid"

	"github.com/sajib-hassan/warden/pkg/logging"
)

type errorCode struct {
	Code   string
	Status int
}

type validationError map[string][]string

func (v validationError) Error() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v validationError) Add(key string, val string) {
	v[key] = append(v[key], val)
}

type apiError struct {
	ID     string          `json:"id"`
	Code   string          `json:"code"`
	Detail json.RawMessage `json:"detail,omitempty"`
	Status int             `json:"status"`
	Title  string          `json:"title"`
	Source error
	Tags   map[string]string
}

func (e apiError) buildTags() map[string]string {
	tags := e.Tags
	if tags == nil {
		tags = map[string]string{}
	}

	tags["error_id"] = e.ID
	tags["title"] = e.Title
	tags["code"] = e.Code
	return tags
}

var (
	//errBadRequest = &errorCode{Code: "400001", Status: http.StatusBadRequest}

	errUnAuthorized = &errorCode{Code: "401001", Status: http.StatusUnauthorized}

	//errPasswordMismatched = &errorCode{Code: "403001", Status: http.StatusForbidden}

	errURINotFound = &errorCode{Code: "404001", Status: http.StatusNotFound}
	//errUserNotFound = &errorCode{Code: "404002", Status: http.StatusNotFound}

	errInvalidMethod = &errorCode{Code: "405001", Status: http.StatusMethodNotAllowed}

	//errEntityNotUnique = &errorCode{Code: "409001", Status: http.StatusConflict}

	//errInvalidData = &errorCode{Code: "422001", Status: http.StatusUnprocessableEntity}

	errInternalServer        = &errorCode{Code: "500001", Status: http.StatusInternalServerError}
	errInternalServerNetPipe = &errorCode{Code: "500002", Status: http.StatusInternalServerError}
)

type apiErrors []apiError

func newAPIError(title string, erC *errorCode, src error) *apiError {
	err := &apiError{
		ID:     uuid.Must(uuid.NewV4()).String(),
		Code:   erC.Code,
		Status: erC.Status,
		Title:  title,
		Source: src,
	}

	if src != nil && erC.Status < 500 {
		if _, valid := src.(*validationError); valid {
			err.Detail = json.RawMessage(src.Error())
		} else {
			err.Detail, _ = json.Marshal(src.Error())
		}
	}

	return err
}

func logFatalErrors(err apiErrors) {
	for _, er := range err {
		logFatalError(&er)
	}
}

func logFatalError(er *apiError) {
	if er.Source == nil || er.Status < 500 {
		return
	}

	raven.CaptureError(er.Source, er.buildTags())
	logging.Logger.Errorf("panic for %+v\n", er.Source)
	logging.Logger.Println(string(debug.Stack()))
}
