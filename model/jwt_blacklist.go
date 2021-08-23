package model

type JwtBlacklist struct {
	MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
