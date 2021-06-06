package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kdefombelle/go-sample/common"
	"github.com/kdefombelle/go-sample/logger"
)

func invalidRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	apiError := common.APIError{
		Error:            "invalid_request",
		ErrorDescription: fmt.Sprintf("%s", err),
	}
	json, _ := json.Marshal(apiError)
	_, _ = w.Write(json)
}

func notFound(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	apiError := common.APIError{
		Error:            "not_found",
		ErrorDescription: fmt.Sprintf("%s", err),
	}
	json, _ := json.Marshal(apiError)
	_, _ = w.Write(json)
}

func conflict(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusConflict)
	apiError := common.APIError{
		Error:            "conflict",
		ErrorDescription: fmt.Sprintf("%s", err),
	}
	json, _ := json.Marshal(apiError)
	_, _ = w.Write(json)
}

func internalServer(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	apiError := common.APIError{
		Error:            "internal",
		ErrorDescription: fmt.Sprintf("%s", err),
	}
	json, _ := json.Marshal(apiError)
	_, _ = w.Write(json)
}

// Check mandatory parameter p which has label l is present,
// otherwise return http.StatusBadRequest (HTTP error 400)
func checkMandatoryParameter(w http.ResponseWriter, p string, l string) {
	if p == "" {
		badRequest(w, errors.New("the "+l+" should be provided"))
	}
}

func badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	apiError := common.APIError{
		Error:            "bad request",
		ErrorDescription: fmt.Sprintf("%s", err),
	}
	json, _ := json.Marshal(apiError)
	_, _ = w.Write(json)
}

// validateRequest reads and parses the JSON as a byte array and stores the result
// in the value pointed to by 'request'. (request is a pointer)
func validateRequest(data io.Reader, request interface{}) error {
	bytesRequest, err := ioutil.ReadAll(data)
	if err != nil {
		logger.Logger.Warnf("Cannot read request [%v]", err)
		return err
	}
	logger.Logger.Debugf(string(bytesRequest))

	err = json.Unmarshal(bytesRequest, request)
	if err != nil {
		logger.Logger.Warnw("Cannot deserialise", "request", request)
		return err
	}

	err = validate.Struct(request)
	if err != nil {
		logger.Logger.Warnw("Request not validated", "request", request)
		return err
	}

	return nil
}
