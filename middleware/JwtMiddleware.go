package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"utils/common"
	"utils/models"
)

type IJwtMiddleware interface {
	ApplicationJwtMiddleware(h http.HandlerFunc) http.HandlerFunc
}

/*
	Implements IJwtMiddleware Interface
*/
type JwtMiddleware struct {
	jwtUtil common.IJwtUtil
}

func GetJwtMiddleware(jwt common.IjwtUtil) IJwtMiddleware {
	return &JwtMiddleware{jwt}
}

/*
	ApplicationJwtMiddleware is HandlerFunc wrapper that looks for an authentication token in the header, query
	or body as a parameter called "authorization" to run through a jwt validator. If it's not valid, it returns an
	error meesage and 401 status. If it's valid, it builds out the context for the request to pass into the final
	endpoint. Context data can be accessed by calling: request.Context().Value("UserID")
*/
func (this JwtMiddleware) ApplicationJwtMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			authorization = r.URL.Query().Get("token")
			if authorization == "" {
				defer r.Body.Close()
				b := bytes.NewBuffer(make([]byte, 0))
				reader := io.TeeReader(r.Body, b)

				authorizationBody := models.AuthorizationString{}
				_ = json.NewDecoder(reader).Decode(&authorizationBody)
				authorization = authorizationBody.Authorization

				// Resets the body with an untouched copy of the body
				r.Body = ioutil.NopCloser(b)
			}
		}

		claims, jwtErr := this.jwtUtil.DecodeApplicationJwt(authorization)
		if jwtErr != nil {
			writeMiddlewareErrorStatusWithMessage(w, http.StatusUnauthorized, jwtErr.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "org_id", claims.OrgID)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "permissions", claims.Permissions)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

/*
	writeMiddlewareErrorStatusWithMessage is a helper function for JwtMiddleware that takes a statusCode and error
	message and writes them to a given responseWriter
*/
func writeMiddlewareErrorStatusWithMessage(response http.ResponseWriter, httpStatus int, message string) {
	returnData := map[string]interface{}{"error": message}
	result, _ := json.Marshal(returnData)

	response.Header().Set("Context-Type", "application/json")
	response.WriteHeader(httpStatus)
	response.Write(result)
}
