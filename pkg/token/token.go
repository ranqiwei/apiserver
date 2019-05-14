package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

var (
	//ErrMissingHeader means the `Authorization` header was empty.
	ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")
)

//Context is  the context of Json web token.  payload
type Context struct {
	Id       uint64
	Username string
}

//-----------------------------------------------------

//secretFunc validates the secret format
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

//Parse validates the token with the specified secret,
//and returns the context if the token was valid.  Parse the payload
func Parse(tokenString string, secret string) (*Context, error) {
	ctx := &Context{}

	token, err := jwt.Parse(tokenString, secretFunc(secret))

	if err != nil {
		return ctx, err
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.Id = uint64(claims["id"].(float64))
		ctx.Username = claims["username"].(string)
		return ctx, nil
	} else {
		return ctx, nil
	}
}

//ParseRequest gets the token from the header and
//pass it to the Parse function to parses the token.
func ParseRequest(c *gin.Context) (*Context, error) { //call Parse
	header := c.Request.Header.Get("Authorization")

	secret := viper.GetString("jwt_secret")
	if len(secret) == 0 {
		return &Context{}, ErrMissingHeader
	}

	var t string
	//Parse the header to get the token part.
	_, _ = fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, secret)
}

//Sign signs the context with the specified secret.
func Sign(c Context, secret string) (tokenString string, err error) {
	if secret == "" {
		secret = viper.GetString("jwt_secret")
	}

	//The token content
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       c.Id,
		"username": c.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err = token.SignedString([]byte(secret))
	return
}
