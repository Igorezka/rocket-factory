module github.com/Igorezka/rocket-factory/payment

go 1.24.3

replace github.com/Igorezka/rocket-factory/shared => ../shared

require (
	github.com/Igorezka/rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.73.0
)

require (
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.36.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
