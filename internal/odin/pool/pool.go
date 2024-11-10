//go:build all && !darwin

package pool

import (
	"context"
	"fmt"

	"github.com/jackc/puddle/v2"
)

func NewContainerPool(ctx context.Context, initPoolSize int32, maxPoolSize int32, engine string) (*puddle.Pool[Container], error) {
	var constructor puddle.Constructor[Container]
	var destructor puddle.Destructor[Container]

	switch engine {
	case "docker":
		constructor = DockerConstructor
		destructor = DockerDestructor
	case "podman":
		constructor = PodConstructor
		destructor = Poddestructor
	default:
		return nil, fmt.Errorf("container engine not supported")
	}

	pool, err := puddle.NewPool(&puddle.Config[Container]{Constructor: constructor, Destructor: destructor, MaxSize: maxPoolSize})
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
}
