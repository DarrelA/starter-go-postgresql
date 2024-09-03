package error

const (
	ErrTypeError            = "type error"
	ErrUUIDError            = "uuid error"
	ErrMsgRedisError        = "redis error"
	ErrMsgPostgresError     = "postgres error"
	ErrMsgGoogleOAuth2Error = "google oauth2 error"

	ErrMsgSomethingWentWrong  = "something went wrong"
	ErrMsgPleaseLoginAgain    = "please login again"
	ErrMsgInvalidToken        = "invalid token"
	ErrMsgInvalidCredentials  = "invalid credentials"
	ErrMsgEmailIsAlreadyTaken = "email is already taken"
)
