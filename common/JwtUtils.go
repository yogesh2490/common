package common

import (
	"crypto/rsa"
	"errors"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"strings"
	"utils/models"
)

var EXPIRED_TOKEN_ERR = errors.New("expired")
var MISSING_TOKEN_ERR = errors.New("missing")
var WITHHELD_TOKEN_ERR = errors.New("withheld")

const RSA_PUB_BEG_COMMENT = "-----BEGIN PUBLIC KEY-----"
const RSA_PUB_END_COMMENT = "-----END PUBLIC KEY-----"

type IJwtUtil interface {
	DecodeApplicationJwt(fullTokenString string) (models.AppJwtClaims, error)
}

/*
	Implements IJwtUtil interface
*/
type JwtUtil struct {
	config_getter IConfigGetter
}

func GetJwtUtil(config_getter IConfigGetter) IJwtUtil {
	return &JwtUtil{config_getter}
}

/*
	Decode JwtToken and return struct containing the jwt
*/
func (this JwtUtil) DecodeApplicationJwt(fullTokenString string) (models.AppJwtClaims, error) {
	if fullTokenString == "" {
		fmt.Println("Missing JWT token")
		return models.AppJwtClaims{}, MISSING_TOKEN_ERR
	}

	jwt := strings.TrimPrefix(fullTokenString, "Bearer ")

	claims := models.AppJwtClaims{}
	token, err := jwtgo.ParseWithClaims(
		jwt,
		&claims,
		func(token *jwtgo.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
				fmt.Println("Unexpected token signing method %s", ok)
				return "", WITHHELD_TOKEN_ERR
			}
			iss := token.Header["iss"]
			switch iss {
			case models.IDM_SIGNED:
				publicKeyValue := this.config_getter.MustGetConfigVar("IDM_PUBLIC_KEY")
				return parsePublicKey(publicKeyValue)
			default:
				fmt.Println("Unknown Token Issuer. iss: %s", iss)
				return "", WITHHELD_TOKEN_ERR
			}
		},
	)

	if err == nil && token.Valid {
		return claims, nil
	} else if ve, ok := err.(*jwtgo.ValidationError); ok {
		if ve.Errors&jwtgo.ValidationErrorMalformed != 0 {
			fmt.Println("Malformed Token: %v", err)
			return models.AppJwtClaims{}, WITHHELD_TOKEN_ERR
		} else if ve.Errors&(jwtgo.ValidationErrorExpired) != 0 {
			fmt.Println("Expired Token: %v", err)
			return models.AppJwtClaims{}, EXPIRED_TOKEN_ERR
		} else if ve.Errors&(jwtgo.ValidationErrorIssuer) != 0 {
			fmt.Println("Token signed by wrong issuer: %v", err)
		} else {
			fmt.Println("Ivalid Token: %v", err)
			return models.AppJwtClaims{}, WITHHELD_TOKEN_ERR
		}
	} else {
		fmt.Println("Token not returned from Parse!")
		return models.AppJwtClaims{}, WITHHELD_TOKEN_ERR
	}
}

/*
	Parse the public key
*/
func parsePublicKey(publicKeyValue string) (*rsa.PublicKey, error) {
	if !strings.HasPrefix(publicKeyValue, RSA_PUB_BEG_COMMENT) {
		publicKeyValue = RSA_PUB_BEG_COMMENT + "\n" + publicKeyValue
	}
	if !strings.HasSuffix(publicKeyValue, RSA_PUB_END_COMMENT) {
		publicKeyValue += "\n" + RSA_PUB_END_COMMENT
	}

	publicKey, parse_err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(publicKeyValue))
	if parse_err != nil {
		fmt.Println("Error parsing public_key: %v", parse_err)
		return nil, errors.New("Could not parse public key from PEM")
	}

	return publicKey, nil
}
