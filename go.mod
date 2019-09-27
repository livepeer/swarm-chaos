module github.com/livepeer/swarm-chaos

go 1.13

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	golang.org/x/net v0.0.0-20190916140828-c8589233b77d // indirect
	google.golang.org/grpc v1.23.1 // indirect
)

replace github.com/docker/docker v1.13.1 => github.com/docker/engine v1.4.2-0.20190822205725-ed20165a37b4
