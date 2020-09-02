package rule

//go:generate go run -tags codegen ../../internal/rule/gen.go -config ../../internal/rule/default.yaml -outfile ../rule/default.go
//go:generate gofmt -s -w ../rule/default.go
