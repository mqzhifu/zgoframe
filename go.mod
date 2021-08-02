module zgoframe

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/casbin/casbin v1.9.1
	github.com/casbin/casbin/v2 v2.11.0
	github.com/casbin/gorm-adapter/v3 v3.0.2
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/mojocn/base64Captcha v1.3.1
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/viper v1.8.1
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.6.7
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.17.0
	gorm.io/driver/mysql v0.3.0
	gorm.io/gorm v1.20.9
	sigs.k8s.io/yaml v1.2.0 // indirect

)

replace github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4

//replace go.etcd.io/bbolt v1.3.4 => github.com/coreos/bbolt v1.3.4

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
