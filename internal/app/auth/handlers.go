package auth

import (
	"database/sql"
	"huangc28/go-ios-iap-vendor/config"
	"huangc28/go-ios-iap-vendor/db"
	"huangc28/go-ios-iap-vendor/internal/apperrors"
	jwtactor "huangc28/go-ios-iap-vendor/internal/pkg/jwtactor"
	"huangc28/go-ios-iap-vendor/internal/pkg/requestbinder"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandlerBody struct {
	Username string `form:"username" json:"username" binding:"required,gt=0"`
	Password string `form:"password" json:"password" binding:"required,gt=0"`
}

func LoginHandler(c *gin.Context) {
	body := LoginHandlerBody{}

	if err := requestbinder.Bind(c, &body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToBindAPIBody,
				err.Error(),
			),
		)

		return
	}

	// Generate retrieve password by username.
	// Password should be hashed
	dao := NewAuthDao(db.GetDB())
	user, err := dao.GetUserByUsername(body.Username)

	if err == sql.ErrNoRows {
		c.JSON(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.UserNotFound,
				err.Error(),
			),
		)

		return
	}

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(
				apperrors.FailedToGetUserByUsername,
				err.Error(),
			),
		)

		return
	}

	// TODO add bcrypt hash on password
	if body.Password != user.Password {
		c.JSON(
			http.StatusBadRequest,
			apperrors.NewErr(
				apperrors.FailedToGetUserByUsername,
				err.Error(),
			),
		)

		return
	}

	// Generate jwt
	jwt, err := jwtactor.CreateToken(
		user.Uuid,
		config.GetAppConf().JWTSecret,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			apperrors.NewErr(apperrors.FailedToGenJWT),
		)

		return
	}

	c.JSON(http.StatusOK, struct {
		Jwt string `json:"jwt"`
	}{
		Jwt: jwt,
	})
}
