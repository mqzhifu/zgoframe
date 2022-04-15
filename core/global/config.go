package global

type Config struct {
	Mysql            ConfigMysql
	Jwt              Jwt
	Zap              Zap
	Captcha          Captcha
	Redis            Redis
	System           System
	Http             Http
	Casbin           Casbin
	Viper            Viper
	Etcd             Etcd
	Alert            Alert
	Metrics          Metrics
	Websocket        Websocket
	Grpc             Grpc
	Email            Email
	Protobuf         Protobuf
	ServiceDiscovery ServiceDiscovery
	PushGateway      PushGateway
	Gateway          Gateway
	ConfigCenter 	 ConfigCenter
}

type Protobuf struct {
	BasePath      string
	PbServicePath string
	ProtoPath     string
	IdMapFileName string
}

type ConfigMysql struct {
	Status       string
	Ip           string
	Port         string
	Config       string
	DbName       string
	Username     string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
	LogMode      bool
	LogZap       bool
}

type Http struct {
	Status        string
	Ip            string
	Port          string
	StaticPath    string
	ReqLimitTimes int
}

type PushGateway struct {
	Status string
	Ip     string
	Port   string
}

type System struct {
	Status    string
	ProjectId int
	//AppId int
	//ServiceId int
	DbType       string
	ENV          int
	ErrorMsgFile string
	OpDirName	string
}

type Jwt struct {
	Status      string
	Key         string
	ExpiresTime int64
	//BufferTime  int64
}

type Viper struct {
	Status string
	Watch  string
}

type Zap struct {
	Status           string
	Level            string
	LevelInt8        int8
	BaseDir          string
	ShowLine         bool
	SoftLinkFileName string
	Format           string

	Prefix        string
	EncodeLevel   string
	AutoAlert     bool
	StacktraceKey string
	LogInConsole  bool
	ModuleName    string //使用设定，不在配置文件中
	FileName      string //使用设定，不在配置文件中
}

type Captcha struct {
	Status       string
	NumberLength int
	ImgWidth     int
	ImgHeight    int
}

type Redis struct {
	Status   string
	Ip       string
	Port     string
	DbNumber int
	Password string
}

type Casbin struct {
	Status    string
	ModelPath string
}

type Etcd struct {
	Status   string
	Ip       string
	Port     string
	Username string
	Password string
	Url      string
}

type ServiceDiscovery struct {
	Status string
	Prefix string
}

type Metrics struct {
	Status string
}

type Alert struct {
	Status string
	Host   string
	Port   string
	Uri    string
}

type Websocket struct {
	Status string
	Uri    string
}

type Grpc struct {
	Status               string
	Ip                   string
	Port                 string
	ServicePackagePrefix string
}

type Email struct {
	Status string
	Host   string
	Ps     string
	Port   string
	From   string
}

type Gateway struct {
	Status   string
	ListenIp string
	OutIp    string
	WsPort   string
	TcpPort  string
	WsUri    string
}

type ConfigCenter struct {
	Status   string
	PersistenceType int
	DataPath    string
}
