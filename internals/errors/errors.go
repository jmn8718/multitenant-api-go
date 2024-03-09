package errors

import (
	"encoding/json"
	goError "errors"
	"fmt"
)

type HttpError struct {
	StatusCode     int    `json:"-"`
	Reason         string `json:"reason,omitempty"`
	Code           string `json:"code"`
	EscalateReason bool   `json:"-"`
}

// creating custom marshal to not escalate reason
func (e HttpError) MarshalJSON() ([]byte, error) {
	if e.EscalateReason {
		return json.Marshal(&struct {
			Reason string `json:"reason"`
			Code   string `json:"code"`
		}{
			Code:   e.Code,
			Reason: e.Reason,
		})
	} else {
		return json.Marshal(&struct {
			Code string `json:"code"`
		}{
			Code: e.Code,
		})
	}
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%d :: code: %s,  reason: %s", e.StatusCode, e.Code, e.Reason)
}

func NewHttpError(statusCode int, error_code string) HttpError {
	return HttpError{
		StatusCode: statusCode,
		Code:       error_code,
	}
}

func NewHttpErrorWithReason(statusCode int, error_code string, reason string, escalate bool) HttpError {
	return HttpError{
		StatusCode:     statusCode,
		Reason:         reason,
		Code:           error_code,
		EscalateReason: escalate,
	}
}

func Is(baseError error, targetError error) bool {
	return goError.Is(baseError, targetError)
}
