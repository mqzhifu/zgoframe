package global

type Config struct {
	Mysql ConfigMysql
	Jwt Jwt
	Zap Zap
	Captcha Captcha
	Redis Redis
	System System
	Gin Gin
	Casbin Casbin
}

type ConfigMysql struct {
	Ip				string
	Port	string
	Config string
	DbName string
	Username string
	Password string
	MaxIdleConns int
	MaxOpenConns int
	LogMode  bool
	LogZap bool
}

type Gin struct {
	Ip string
	Port string
	StaticPath string
	ReqLimitTimes int
}

type System struct {
	DbType string
	ENV string
}

type Jwt struct {
	Key string
	ExpiresTime	int64
	BufferTime int64
}

type Zap struct {
	Level string
	Dir string
	ShowLine bool
	LinkName string
	LogInConsole bool
	Format string
	StacktraceKey string
	Prefix string
	EncodeLevel string
}

type Captcha struct {
	NumberLength	int
	ImgWidth		int
	ImgHeight		int
}

type Redis struct {
	Ip 		string
	Port 	string
	DbNumber int
	Password string
}

type Casbin struct {
	ModelPath string
}