package request

import (
	"github.com/dgrijalva/jwt-go"
)

// Custom claims structure
type CustomClaims struct {
	//AuthorityId string
	//UUID        uuid.UUID
	ProjectId 		int		`json:"project_id"`
	SourceType 	int		`json:"source_type"`
	Id          int		`json:"id"`
	Username    string	`json:"username"`
	NickName    string	`json:"nick_name"`

	BufferTime  int64	`json:"buffer_time"`
	jwt.StandardClaims
}
