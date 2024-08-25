package server

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/deepakdinesh1123/valkyrie/internal/logs"
	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
)

type ExecutionSuite struct {
	suite.Suite
	envConfig  *config.EnvConfig
	logger     *zerolog.Logger
	odinServer *OdinServer
	server     *httptest.Server
	workers    []*testcontainers.Container
}

func (s *ExecutionSuite) SetupSuite() {
	envConfig, err := config.GetEnvConfig(
		config.WithPostgresDB("test"),
		config.WithPostgresUser("test"),
		config.WithPostgresPassword("test"),
		config.WithPostgresPort(5400),
	)
	if err != nil {
		s.T().Errorf("failed to get environment config -> %s", err)
	}

	s.envConfig = envConfig
	s.logger = logs.GetLogger("debug")
	srv, err := NewServer(context.Background(), s.envConfig, false, true, s.logger)
	if err != nil {
		s.T().Errorf("failed to create server -> %s", err)
	}
	s.odinServer = srv
	s.server = httptest.NewServer(srv.server)

	workerContainer, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../../../.",
				Dockerfile: "build/package/dockerfiles/odin.worker.dockerfile",
			},
		},
	})
	if err != nil {
		s.T().Errorf("failed to create worker container -> %s", err)
	}

	s.workers = append(s.workers, &workerContainer)
	err = workerContainer.Start(context.TODO())
	if err != nil {
		s.T().Errorf("failed to start worker container -> %s", err)
	}
}

func (s *ExecutionSuite) TearDownSuite() {
	s.server.Close()
}

func (s *ExecutionSuite) TestExecute() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := api.NewClient(s.server.URL)
	if err != nil {
		s.T().Errorf("failed to create client -> %s", err)
	}
	res, err := client.Execute(ctx, &api.ExecutionRequest{
		Code:     "print('hello')",
		Language: "python",
	})
	switch res := res.(type) {
	case *api.ExecuteBadRequest:
		s.T().Errorf("bad request -> %s", res.Message)
	case *api.ExecuteInternalServerError:
		s.T().Errorf("internal server error -> %s", res.Message)
	}
	assert.Nil(s.T(), err)
}
func TestExecutionSuite(t *testing.T) {
	suite.Run(t, new(ExecutionSuite))
}
