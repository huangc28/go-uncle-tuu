package apperrors

const (
	FailedToGetUserByUsername = "2000001"
	PasswordNotMatch          = "2000002"
	UserNotFound              = "2000003"
	FailedToGenJWT            = "2000004"
)

var authErrMap = map[string]string{
	FailedToGetUserByUsername: "failed to get user by username",
	PasswordNotMatch:          "password not match",
	UserNotFound:              "user not found",
	FailedToGenJWT:            "failed to generate jwt token",
}
