module github.com/tjalfe/jndispatcher

go 1.24.3

replace github.com/tjalfe/pcrypt => ../pcrypt

replace github.com/tjalfe/randomstring => ../randomstring

replace github.com/tjalfe/psign/pkg/psign => ../psign/pkg/psign

require (
	github.com/spf13/pflag v1.0.10
	github.com/tjalfe/pcrypt v0.0.0-00010101000000-000000000000
)

require (
	github.com/tjalfe/randomstring v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/term v0.26.0 // indirect
)
