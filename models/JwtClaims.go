package models

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"time"
)

const IDM_SIGNED = "IDM"

/*
	Token from application
*/
type AppJwtClaims struct {
	OrgID        int      `json:orgid`
	UserID       string   `json:idmGuid`
	Premissions  []string `json:permissions`
	Capabilities []string `json:capabilities`
	jwtgo.StandardClaims
}

type AuthorizationString struct {
	Authorization string `json:authorization`
}

/*
	Validate the claims
*/
func (c AppJwtClaims) Valid() error {
	vErr := new(jwtgo.ValidationError)
	now := time.Now().Unix()

	if c.VerifyIssuer(IDM_SIGNED, true) == false {
		vErr.Errors |= jwtgo.ValidationErrorIssuer
		return vErr
	}

	if c.VerifyExpiresAt(now, true) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwtgo.ValidationErrorExpired
		return vErr
	}

	if c.UserID == "" || c.OrgID == 0 {
		vErr.Errors |= jwtgo.ValidationErrorClaimsInvalid
		return vErr
	}
	if vErr.Errors == 0 {
		return c.StandardClaims.Valid()
	}

	return vErr
}
