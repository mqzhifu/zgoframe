module zgoframe

go 1.16

//replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4
replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.6

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/abrander/go-supervisord v0.0.0-20210517172913-a5469a4c50e2
	github.com/alibabacloud-go/darabonba-openapi/v2 v2.0.4
	github.com/alibabacloud-go/dysmsapi-20170525/v3 v3.0.5
	github.com/alibabacloud-go/tea v1.1.20
	github.com/aliyun/aliyun-oss-go-sdk v2.2.7+incompatible
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.27+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20230327231512-ba87abf18a23 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.6.0
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3
	github.com/gomodule/redigo v1.8.9
	github.com/gorilla/websocket v1.5.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mojocn/base64Captcha v1.3.5
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/common v0.42.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.3.2
	github.com/swaggo/gin-swagger v1.0.0
	github.com/swaggo/swag v1.7.9
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/ugorji/go/codec v0.0.0-20181209151446-772ced7fd4c2 // indirect
	github.com/xiang90/probing v0.0.0-20221125231312-a49e3df8f510 // indirect
	go.etcd.io/etcd v3.3.27+incompatible
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/net v0.8.0
	golang.org/x/tools v0.6.0
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.54.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gorm.io/driver/mysql v1.5.0
	gorm.io/gorm v1.25.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)
