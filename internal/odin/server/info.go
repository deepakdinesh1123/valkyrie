package server

import (
	"context"

	"github.com/deepakdinesh1123/valkyrie/internal/odin/config"
	"github.com/deepakdinesh1123/valkyrie/pkg/odin/api"
)

func (s *OdinServer) GetVersion(ctx context.Context, params api.GetVersionParams) (api.GetVersionRes, error) {
	return &api.GetVersionOK{
		Version: config.ODIN_VERSION,
	}, nil
}

func (s *OdinServer) GetExecutionConfig(ctx context.Context, params api.GetExecutionConfigParams) (api.GetExecutionConfigRes, error) {
	user := ctx.Value(config.UserKey).(string)

	if user != "admin" {
		return &api.GetExecutionConfigForbidden{}, nil
	}

	return &api.ExecutionConfig{
		ODINWORKERPROVIDER:    s.envConfig.ODIN_WORKER_PROVIDER,
		ODINWORKERCONCURRENCY: int32(s.envConfig.ODIN_WORKER_CONCURRENCY),
		ODINWORKERBUFFERSIZE:  int32(s.envConfig.ODIN_WORKER_BUFFER_SIZE),
		ODINWORKERTASKTIMEOUT: s.envConfig.ODIN_WORKER_TASK_TIMEOUT,
		ODINWORKERPOLLFREQ:    s.envConfig.ODIN_WORKER_POLL_FREQ,
		ODINWORKERRUNTIME:     s.envConfig.ODIN_WORKER_RUNTIME,
		ODINLOGLEVEL:          s.envConfig.ODIN_LOG_LEVEL,
	}, nil
}

func (s *OdinServer) GetAllLanguages(ctx context.Context, params api.GetAllLanguagesParams) (api.GetAllLanguagesRes, error) {
	var languages []string
	for lang := range config.Languages {
		languages = append(languages, lang)
	}
	return &api.GetAllLanguagesOK{
		Languages: languages,
	}, nil
}
