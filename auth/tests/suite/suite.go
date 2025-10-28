package suite

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	GRPC: scfg.GRPCSrvConfig{
		Port: 50051,
	},
	Database: scfg.DatabaseConfig{
		Host:       "postgres",
		Port:       5432,
		SSLMode:    false,
		Name:       "postgres",
		User:       "admin",
		Password:   "pass",
		Migrations: "./migrations",
	},
	HealthPort: 3000,
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

	dbRes, err := setupPostgres(testName, pool, cfg, networkName)
	if err != nil {
		fmt.Printf("Failed setup PostgreSQL: %v", err)
		panic(err)
	}
	resources = append(resources, dbRes)

	cfg.Database.Host = strings.TrimPrefix(dbRes.Container.Name, "/")
	authRes, port := mustSetupAuth(testName, pool, cfg, networkName, time.Second*20)
	resources = append(resources, authRes)

	cc, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed create grpc client: %v", err)
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

func setupPostgres(testName string, pool *dockertest.Pool, cfg *config.Config, network string) (*dockertest.Resource, error) {
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
		return nil, err
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
		return nil, err
	}

	return res, nil
}

func mustSetupAuth(testName string, pool *dockertest.Pool, cfg *config.Config, network string, setupTimeout time.Duration) (*dockertest.Resource, uint16) {
	exposedPort := fmt.Sprintf("%d/tcp", cfg.GRPC.Port)
	exposedHealthPort := fmt.Sprintf("%d/tcp", cfg.HealthPort)
	res, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "got-auth",
		Tag:        "latest",
		Name:       "auth-" + testName,
		NetworkID:  network,
		Env: []string{
			fmt.Sprintf("ENV=%s", cfg.Env),
			fmt.Sprintf("TOKEN_TTL=%v", cfg.TokenTTL),
			fmt.Sprintf("TOKEN_SECRET=%s", cfg.TokenSecret),
			fmt.Sprintf("GRPC_PORT=%d", cfg.GRPC.Port),
			fmt.Sprintf("DB_HOST=%s", cfg.Database.Host),
			fmt.Sprintf("DB_PORT=%d", cfg.Database.Port),
			fmt.Sprintf("DB_NAME=%s", cfg.Database.Name),
			fmt.Sprintf("DB_USER=%s", cfg.Database.User),
			fmt.Sprintf("DB_PASSWORD=%s", cfg.Database.Password),
			fmt.Sprintf("DB_MIGRATIONS=%s", cfg.Database.Migrations),
			fmt.Sprintf("HEALTH_PORT=%d", cfg.HealthPort),
		},
		ExposedPorts: []string{exposedPort, exposedHealthPort},
	}, func(config *docker.HostConfig) {
		config.PortBindings = map[docker.Port][]docker.PortBinding{
			docker.Port(exposedPort):       {{HostIP: "0.0.0.0", HostPort: ""}},
			docker.Port(exposedHealthPort): {{HostIP: "0.0.0.0", HostPort: ""}},
		}
	})
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		if err := pool.Purge(res); err != nil {
			log.Printf("Could not purge container: %s", err)
		}
	}

	hostPort, err := strconv.Atoi(res.GetPort(exposedPort))
	if err != nil {
		cleanup()
		panic(err)
	}

	healthHostPort, err := strconv.Atoi(res.GetPort(exposedHealthPort))
	if err != nil {
		cleanup()
		panic(err)
	}

	err = waitHealthCheck(fmt.Sprintf("http://localhost:%d/readyz", healthHostPort), setupTimeout)
	if err != nil {
		cleanup()
		panic(err)
	}

	return res, uint16(hostPort)
}

func waitHealthCheck(address string, timeout time.Duration) error {
	startTime := time.Now()
	for {
		resp, err := http.Get(address)
		if err != nil {
			continue
		}
		if time.Since(startTime) > timeout {
			return fmt.Errorf("%s healthcheck timeout", address)
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}
