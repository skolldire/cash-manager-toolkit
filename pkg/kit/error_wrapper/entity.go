package error_wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CommonApiError struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	Err      error  `json:"-"`
	HttpCode int    `json:"-"`
}

func (e *CommonApiError) Error() string {
	return fmt.Sprintf("Error %s: %s \ntrace: %s", e.Code, e.Msg, e.Err.Error())
}

func (e *CommonApiError) Unwrap() error {
	return e.Err
}

func NewCommonApiError(code, msg string, err error, httpCode int) error {
	return &CommonApiError{
		Code:     code,
		Msg:      msg,
		Err:      err,
		HttpCode: httpCode,
	}
}

func HandleApiErrorResponse(err error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")

	var errType *CommonApiError
	if errors.As(err, &errType) {
		if errType.Err == nil {
			err = fmt.Errorf("[error_wrapper]HandleApiErrorResponse: The err attribute is null")
		}
		fmt.Println("CommonApiError: %w", err)
		w.WriteHeader(errType.HttpCode)
		b, _ := json.Marshal(&errType)
		_, _ = w.Write(b)
		return nil
	}

	fmt.Println("Error: %w", err)
	w.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(CommonApiError{Code: "GE-001", Msg: "Internal Error"})
	_, _ = w.Write(b)
	return nil
}
