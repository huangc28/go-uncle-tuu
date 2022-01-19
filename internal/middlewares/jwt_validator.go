package middlewares

import (
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claim struct {
	Uuid       string `json:"uuid"`
	Authorized bool   `json:"authorized"`
	jwt.StandardClaims
}

type JwtMiddlewareOptions struct {
	Secret string
}

type JWTToken struct {
	AuthJWT string `header:"Authorization" form:"jwt"`
}

func extractTokenFromRequest(c *gin.Context) (string, error) {
	headerJWT := JWTToken{}

	if err := c.ShouldBindHeader(&headerJWT); err != nil {
		return "", err
	}

	if len(headerJWT.AuthJWT) > 0 {
		strArr := strings.Split(headerJWT.AuthJWT, " ")

		if len(strArr) >= 2 {
			return strArr[1], nil
		}

		return headerJWT.AuthJWT, nil
	}

	if err := c.ShouldBindQuery(&headerJWT); err != nil {
		return "", err
	}

	return headerJWT.AuthJWT, nil
}

func JWTValidator(opt JwtMiddlewareOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractTokenFromRequest(c)

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				apperrors.NewErr(
					apperrors.FailedToBindJWTInHeader,
					err.Error(),
				),
			)

			return
		}

		if len(token) <= 0 {
			c.JSON(
				http.StatusInternalServerError,
				apperrors.NewErr(
					apperrors.MissingJWTToken,
					err.Error(),
				),
			)

			return
		}

		claims := &Claim{}
		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(opt.Secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(
					http.StatusUnauthorized,
					apperrors.NewErr(apperrors.MissingJWTToken),
				)

				return
			}

			c.JSON(
				http.StatusBadRequest,
				apperrors.NewErr(

					apperrors.FailedToParseSignature,
					err.Error(),
				),
			)

			return

		}

		if !tkn.Valid {
			c.JSON(
				http.StatusUnauthorized,
				apperrors.NewErr(
					apperrors.InvalidSigature,
					err.Error(),
				),
			)

			return
		}

		// @TODO check redis makesure incoming jwt token is still valid.
		c.Set("user_uuid", claims.Uuid)
		c.Set("jwt", token)

		c.Next()
	}
}
