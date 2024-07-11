package entity

type Token struct {
	Token     *string
	TokenUUID string
	UserUUID  string
	ExpiresIn *int64
}
