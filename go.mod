module github.com/tjalfe/jndispatcher

go 1.24.3

replace github.com/tjalfe/pcrypt => ../pcrypt

replace github.com/tjalfe/randomstring => ../randomstring

replace github.com/tjalfe/psign => ../psign

require (
	github.com/spf13/pflag v1.0.10
	github.com/tjalfe/pcrypt v0.0.0-00010101000000-000000000000
	github.com/tjalfe/psign v0.0.0-00010101000000-000000000000
	github.com/twmb/franz-go v1.19.5
	gopkg.in/yaml.v3 v3.0.1
	software.sslmate.com/src/go-pkcs12 v0.6.0
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.11.2 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/term v0.35.0 // indirect
)
