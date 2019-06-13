module kubesphere.io/alert

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/emicklei/go-restful v2.9.3+incompatible
	github.com/emicklei/go-restful-openapi v1.0.0
	github.com/fatih/camelcase v1.0.0
	github.com/fatih/structs v1.1.0
	github.com/gin-contrib/sse v0.0.0-20190301062529-5545eab6dad3 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-openapi/spec v0.19.0 // indirect
	github.com/golang/protobuf v1.3.1
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/google/gops v0.3.6
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.8.5
	github.com/jinzhu/gorm v1.9.4
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/stretchr/testify v1.3.0
	github.com/ugorji/go v1.1.4 // indirect
	golang.org/x/net v0.0.0-20190311183353-d8887717615a
	golang.org/x/tools v0.0.0-20190312170243-e65039ee4138
	google.golang.org/genproto v0.0.0-20190307195333-5fe7a883aa19
	google.golang.org/grpc v1.20.1
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	k8s.io/api v0.0.0-20181213150558-05914d821849 // indirect
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	k8s.io/klog v0.3.0 // indirect
	openpitrix.io/libqueue v0.3.1
	openpitrix.io/logger v0.1.0
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace openpitrix.io/libqueue v0.3.1 => github.com/openpitrix/libqueue v0.3.1
