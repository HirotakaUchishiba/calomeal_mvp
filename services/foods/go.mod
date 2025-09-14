module github.com/HirotakaUchishiba/calomeal_mvp/services/foods

go 1.25.0

require (
	github.com/HirotakaUchishiba/calomeal_mvp/proto/foods/v1 v0.0.0
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9
)

require (
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)

replace github.com/HirotakaUchishiba/calomeal_mvp/proto/foods/v1 => ../../proto/foods/v1
