package controllers

import (
	"encoding/json"
	"net/http"
)

type IController interface {
	writeErrorStatusWithMessage(response http.ResponseWriter, httpStatus int, message string)
}

type controller struct{}

/*
	Writes the provided http status out and returns with the provided message in json with the key the "error"
	Example: {"error": "unknown_error"}
*/
func (this *controller) writeErrorStatusWithMessage(response http.ResponseWriter, httpStatus int, message string) {
	returnData := map[string]interface{}{
		"error": message,
	}
	result, _ := json.Marshal(returnData)

	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.WriteHeader(httpStatus)
	response.Write(result)
}
