package pool

import (
	"context"
	"fmt"

	"github.com/jackc/puddle/v2"
)

func NewSandboxPool(ctx context.Context, initPoolSize int32, maxPoolSize int32, engine string) (*puddle.Pool[Container], error) {
	switch engine {
	case "docker":
		pool, err := puddle.NewPool(&puddle.Config[Container]{Constructor: DockerSandboxConstructor, Destructor: DockerSandboxDestructor, MaxSize: maxPoolSize})
		if err != nil {
			return nil, err
		}
		for i := 0; i < int(initPoolSize); i += 1 {
			err := pool.CreateResource(ctx)
			if err != nil {
				return nil, err
			}
		}
		return pool, nil
	default:
		return nil, fmt.Errorf("container engine %s not supported", engine)
	}
}
