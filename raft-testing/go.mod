module github.com/netrixframework/raft-testing

go 1.17

replace (
	github.com/netrixframework/netrix v0.3.5 => /go/src/github.com/netrixframework/netrix
)


require go.etcd.io/etcd/raft/v3 v3.5.1

require (
	github.com/cockroachdb/datadriven v0.0.0-20200714090401-bf6692d28da5
	github.com/netrixframework/go-clientlibrary v0.1.7
	github.com/netrixframework/netrix v0.3.5
	github.com/spf13/cobra v1.4.0
	go.etcd.io/etcd/client/pkg/v3 v3.5.2
	go.etcd.io/etcd/server/v3 v3.5.1
	go.uber.org/zap v1.17.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cockroachdb/errors v1.2.4 // indirect
	github.com/cockroachdb/logtags v0.0.0-20190617123548-eb05cc24525f // indirect
	github.com/getsentry/raven-go v0.2.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.7.7 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/exp v0.0.0-20220428152302-39d4317da171 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
