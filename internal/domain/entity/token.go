package entity

import "time"

type Token struct {
	Token      *string
	TokenUUID  string
	UserUUID   string
	FamilyUUID string
	ExpiresIn  *int64
}

type TokenSession struct {
	UserUUID         string
	FamilyUUID       string
	AccessTokenUUID  string
	RefreshTokenUUID string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}
