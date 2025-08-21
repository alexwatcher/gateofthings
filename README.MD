# Gate Of Things

### Tools

project use Task for automation scripts

```
brew install protobuf

go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go get github.com/go-task/task/v3/cmd/task@latest
go get github.com/envoyproxy/protoc-gen-validate@latest
go get github.com/pressly/goose/v3/cmd/goose@latest
```

### Structure
All services organized in shared Go workspace.
**auth** - authorization service
**shared** - library of common code and data structures
**protos** - protobuf and grps definitions
**configs** - configuration files for dependencies (OTL collector, Loki, Tempo, Prometheus, Grafana)

### Local setup
Local development and testing are powered by Docker Compose, which orchestrates all necessary services and dependencies.

```
docker compose up -d
docker compose down -v
```

### Tasks
The project uses Task as a task runner to automate common development workflows, such as generating protobuf code.
Root direcotry has Task file which include all subprojects.

To work with protobufs use following commands:
```
task protos:gen
task protos:clr
```

### Functional tests
Functional tests use dockertest to spin up test containers.
By default, dockertest connects to Docker via:
```
/var/run/docker.sock
```

However, on macOS, Docker Desktop may expose the socket at a different path:
```
$HOME/.docker/run/docker.sock
```

Solution: set path into DOCKER_HOST
```
export DOCKER_HOST=unix://${HOME}/.docker/run/docker.sock
```
