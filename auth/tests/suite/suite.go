package suite

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alexwatcher/gateofthings/auth/internal/config"
	authv1 "github.com/alexwatcher/gateofthings/protos/gen/go/auth/v1"
	scfg "github.com/alexwatcher/gateofthings/shared/pkg/config"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const testsDuration = 120 * time.Second

var cfg = &config.Config{
	Env:         "test",
	TokenTTL:    time.Minute * 10,
	TokenSecret: "test",
	Telemetry:   scfg.TelemetryConfig{},
	GRPC:        scfg.GRPCSrvConfig{},
	Database: scfg.DatabaseConfig{
		Host:       "postgres",
		Port:       5432,
		SSLMode:    false,
		Name:       "postgres",
		User:       "admin",
		Password:   "pass",
		Migrations: "./migrations",
	},
}

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient authv1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite, func()) {
	t.Helper()
	t.Parallel()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// create network
	testName := t.Name()
	networkName := "test-network-" + testName
	network, err := pool.CreateNetwork(networkName)
	if err != nil {
		log.Fatalf("Could not create network: %s", err)
	}

	var resources []*dockertest.Resource
	defer func() {
		if r := recover(); r != nil {
			for _, res := range resources {
				if err := pool.Purge(res); err != nil {
					log.Printf("Could not purge container: %s", err)
				}
			}
			if err := pool.RemoveNetwork(network); err != nil {
				log.Printf("Could not remove network: %s", err)
			}
			panic(r)
		}
	}()

	dbRes := mustSetupPostgres(testName, pool, cfg, networkName)
	resources = append(resources, dbRes)

	cfg.Database.Host = strings.TrimPrefix(dbRes.Container.Name, "/")
	authRes, port := mustSetupAuth(testName, pool, cfg, networkName)
	resources = append(resources, authRes)

	cc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), testsDuration)
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
		cc.Close()
	})

	return ctx, &Suite{t, cfg, authv1.NewAuthClient(cc)}, func() {
		for _, res := range resources {
			if err := pool.Purge(res); err != nil {
				log.Printf("Could not purge container: %s", err)
			}
		}
		if err := pool.RemoveNetwork(network); err != nil {
			log.Printf("Could not remove network: %s", err)
		}
	}
}

func mustSetupPostgres(testName string, pool *dockertest.Pool, cfg *config.Config, network string) *dockertest.Resource {
	res, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17.5",
		Name:       "postgres-" + testName,
		NetworkID:  network,
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", cfg.Database.User),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cfg.Database.Password),
			fmt.Sprintf("POSTGRES_DB=%s", cfg.Database.Name),
		},
	})
	if err != nil {
		panic(err)
	}

	err = pool.Retry(func() error {
		out, err := exec.Command("docker", "exec", res.Container.ID, "pg_isready", "-U", "root").CombinedOutput()
		if err != nil {
			return fmt.Errorf("pg_isready failed: %v - output: %s", err, string(out))
		}
		if !strings.Contains(string(out), "accepting connections") {
			return fmt.Errorf("postgres not ready: %s", out)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return res
}

func mustSetupAuth(testName string, pool *dockertest.Pool, cfg *config.Config, network string) (*dockertest.Resource, uint16) {
	port := 3000
	exposedPort := fmt.Sprintf("%d/tcp", port)
	res, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "got-auth",
		Tag:        "latest",
		Name:       "auth-" + testName,
		NetworkID:  network,
		Env: []string{
			fmt.Sprintf("ENV=%s", cfg.Env),
			fmt.Sprintf("TOKEN_TTL=%v", cfg.TokenTTL),
			fmt.Sprintf("TOKEN_SECRET=%s", cfg.TokenSecret),
			fmt.Sprintf("GRPC_PORT=%d", port),
			fmt.Sprintf("DB_HOST=%s", cfg.Database.Host),
			fmt.Sprintf("DB_PORT=%d", cfg.Database.Port),
			fmt.Sprintf("DB_NAME=%s", cfg.Database.Name),
			fmt.Sprintf("DB_USER=%s", cfg.Database.User),
			fmt.Sprintf("DB_PASSWORD=%s", cfg.Database.Password),
			fmt.Sprintf("DB_MIGRATIONS=%s", cfg.Database.Migrations),
		},
		ExposedPorts: []string{exposedPort},
	}, func(config *docker.HostConfig) {
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			docker.Port(exposedPort): {{HostIP: "0.0.0.0", HostPort: ""}},
		}
	})
	if err != nil {
		panic(err)
	}

	hostPort, err := strconv.Atoi(res.GetPort(exposedPort))
	if err != nil {
		if err := pool.Purge(res); err != nil {
			log.Printf("Could not purge container: %s", err)
		}
		panic(err)
	}

	// TODO: implement and then use healthcheck to wait availability of auth service
	time.Sleep(time.Second * 5)

	return res, uint16(hostPort)
}
