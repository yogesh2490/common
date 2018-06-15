package common

import (
	"net/http"
)

/*
	This class helps to determine service health
*/
type IHealthController interface {
	GetHealth(w http.ResponseWriter, r *http.Request)
}

// Implements IHealthController
type healthController struct{}

/*
	Return an implementation of IHealthController
*/

func GetHealthController() IHealthController {
	return &healthController{}
}

/*
	This returns a UTF-8 charset response: "Alive"
*/

func (this *healthController) GetHealth(w http.ResponseWriter, r *http.Request) {
	result := []byte("Alive")

	w.Header().Set("Content-Type", "charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
