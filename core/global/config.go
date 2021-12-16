package global

type Config struct {
	Mysql   ConfigMysql
	Jwt     Jwt
	Zap     Zap
	Captcha Captcha
	Redis   Redis
	System  System
	Http     Http
	Casbin  Casbin
	Viper 	Viper
	Etcd 	Etcd
	ServiceDiscovery ServiceDiscovery
	Alert Alert
	Metrics Metrics
	Websocket Websocket
	Grpc Grpc
	Email Email
}

type ConfigMysql struct {
	Status string
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

type Http struct {
	Status string
	Ip string
	Port string
	StaticPath string
	ReqLimitTimes int
}

type System struct {
	Status string
	AppId int
	ServiceId int
	DbType string
	ENV string
	ErrorMsgFile string
}

type Jwt struct {
	Status string
	Key string
	ExpiresTime	int64
	BufferTime int64
}

type Viper struct {
	Status string
	Watch string
}

type Zap struct {
	Status string
	Level string
	Dir string
	ShowLine bool
	LinkName string
	LogInConsole bool
	Format string
	StacktraceKey string
	Prefix string
	EncodeLevel string
	AutoAlert bool
}

type Captcha struct {
	Status string
	NumberLength	int
	ImgWidth		int
	ImgHeight		int
}

type Redis struct {
	Status string
	Ip 		string
	Port 	string
	DbNumber int
	Password string
}

type Casbin struct {
	Status string
	ModelPath string
}

type Etcd struct{
	Status string
	Ip string
	Port string
	Username string
	Password string
	Url string
}

type ServiceDiscovery struct {
	Status string
}

type Metrics struct {
	Status string
}

type Alert struct {
	Status string
	Ip string
	Port string
	Uri string
}

type Websocket struct {
	Status string
	Uri string
}

type Grpc struct {
	Status string
	Ip string
	Port string
}

type Email struct {
	Status string
	Host string
	Ps string
	Port int
	From string
}

