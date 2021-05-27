module github.com/c3sr/evaluation

go 1.15

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	github.com/jaegertracing/jaeger => github.com/uber/jaeger v1.22.0
	github.com/uber/jaeger => github.com/jaegertracing/jaeger v1.22.0
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/GeertJohan/go-sourcepath v0.0.0-20150925135350-83e8b8723a9b
	github.com/c3sr/config v1.0.1
	github.com/c3sr/database v1.0.0
	github.com/c3sr/dlframework v1.3.2
	github.com/c3sr/go-echarts v1.0.0
	github.com/c3sr/logger v1.0.1
	github.com/c3sr/machine v1.0.0
	github.com/c3sr/nvidia-smi v1.0.0
	github.com/c3sr/parallel v1.0.1
	github.com/c3sr/tracer v1.0.1
	github.com/c3sr/utils v1.0.0
	github.com/chewxy/math32 v1.0.8
	github.com/fatih/structs v1.1.0
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/golang/snappy v0.0.3
	github.com/ianlancetaylor/demangle v0.0.0-20210406231658-61c622dd7d50
	github.com/k0kubun/pp/v3 v3.0.7
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/mailru/easyjson v0.7.7
	github.com/mattn/go-zglob v0.0.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/errors v0.9.1
	github.com/ready-steady/assert v0.0.0-20171126095531-4075406641e2 // indirect
	github.com/ready-steady/sort v0.0.0-20151130154609-c3763d4578b8
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.7.0
	github.com/thoas/go-funk v0.8.0
	github.com/uber/jaeger v1.22.0
	github.com/unknwon/com v1.0.1
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	upper.io/db.v3 v3.8.0+incompatible
)
