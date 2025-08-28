# Gate Of Things

### Tools

project use Task for automation scripts

use command to install task
```
go get github.com/go-task/task/v3/cmd/task@latest
```

also project use protobuf, use following command to install it on mac os
```
brew install protobuf
```

to install other dependencies, use following Task command
```
task install_tools
```

### Structure
All services organized in shared Go workspace.
**auth** - authorization service
**gateway** - http-grpc reverse proxy
**shared** - library of common code and data structures
**protos** - protobuf and grps definitions
**configs** - configuration files for dependencies (OTL collector, Loki, Tempo, Prometheus, Grafana)

### Local setup
Local development and testing are powered by Docker Compose, which orchestrates all necessary services and dependencies.
To build all services run following command:
```
task build
```
it include build of migrator, gateway, auth services

To bring up the services run following command:
```
task up
```

To shutdown environment:
```
task down
```

### gRPC
To generate protobugs and grpc use command
```
task protos:gen
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
