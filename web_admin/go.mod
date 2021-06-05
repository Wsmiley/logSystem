module web_admin

go 1.15

//require github.com/astaxie/beego v1.12.1

require (
	github.com/astaxie/beego v1.12.3
	github.com/beego/admin v0.0.0-20210305083807-6b74f2e7468f
	github.com/beego/beego/v2 v2.0.1
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/jmoiron/sqlx v1.3.3
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/smartystreets/goconvey v1.6.4
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/genproto v0.0.0-20210415145412-64678f1ae2d5 // indirect
	google.golang.org/grpc v1.37.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
