module jrplugin

go 1.23

replace github.com/jrnd-io/jrv2 => ../..

require (
	github.com/hashicorp/go-plugin v1.6.1
	//	github.com/jrnd-io/jrv2 v0.0.0-20240824134657-26a112d020c0
	google.golang.org/protobuf v1.34.2
)

require github.com/jrnd-io/jrv2 v0.0.0-00010101000000-000000000000

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/oklog/run v1.0.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/grpc v1.65.0 // indirect
)
