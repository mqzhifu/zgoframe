package model

import (
	"zgoframe/global"
)

type JwtBlacklist struct {
	global.MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
