package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jllopis/try5/keys"
	logger "github.com/jllopis/try5/log"
	"github.com/jllopis/try5/store/manager"
	"github.com/jllopis/try5/tryerr"
)

// GenerateToken genera un JSON Web Token para el account identificado con 'uid'.
// La cuenta debe estar creada antes de poder generar el JWT y debe contar con las
// claves RSA.
//
// Devuelve un 'string' con el token firmado con la clave privada correspondiente a la
// cuenta.
//
// El tiempo de expiración del token es de 24 horas desde su expedición ('exp' claim).
func GenerateToken(uid string, privKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Header["kid"] = uid
	token.Claims["iat"] = time.Now().UTC().Unix()
	token.Claims["exp"] = time.Now().UTC().Add(time.Duration(24 * time.Hour)).Unix()
	token.Claims["sub"] = uid
	// jwt-go utiliza las claves codificadas en formato PEM
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		logger.LogE("GenerateToken error", "error", err.Error())
		return "", err
	}
	return tokenString, nil
}

// Validate realiza la verificación del JWT. Si es valido devuelve 'nil' y se devolverá el error
// obtenido en caso contrario.
//
// Los parámetros necesario son un context.Context que contenga la referencia al 'store.Storer' de
// la aplicación para localizar la clave pública correspondiente a la cuenta con que se ha firmado
// el JWT. El segundo parámetro es el JWT que se va a verificar.
//
// La función comprueba que el método de firma utilizado (RSA) corresponde con el especificado en el token.
func Validate(m *manager.Manager, t string) error {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// validate the expected algo
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, tryerr.ErrJWTWrongSigningMethod
		}
		var pkey *keys.Key
		var err error
		if pkey, err = m.GetKeyByAccountID(token.Header["kid"].(string)); err != nil {
			logger.LogE("error getting private key", "error", err.Error())
			return nil, err
		}
		logger.LogD("jwt.Parse", "Found PrivKey", pkey.PubKey)
		return pkey.PubKey, nil
	})

	if err == nil && token.Valid {
		return nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			logger.LogE("that's not even a token", "pkg", "jwt", "func", "Validate()", "wannabe token", ve)
			return ve
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			logger.LogE("token expired or not active yet", "pkg", "jwt", "func", "Validate()", "wannabe token", ve)
			return ve
		} else {
			logger.LogE("(ValidationError) Couldn't handle this token", "pkg", "jwt", "func", "Validate()", "token", ve)
			return tryerr.ErrUnauthorized
		}
	} else {
		logger.LogE("(ValidationError) Couldn't handle this token", "pkg", "jwt", "func", "Validate()", "token", ve)
		return tryerr.ErrUnauthorized
	}
}
