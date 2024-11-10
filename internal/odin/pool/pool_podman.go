//go:build podman && !darwin

package pool

import (
	"context"

	"github.com/jackc/puddle/v2"
)

func NewContainerPool(ctx context.Context, initPoolSize int32, maxPoolSize int32, engine string) (*puddle.Pool[Container], error) {

	pool, err := puddle.NewPool(&puddle.Config[Container]{Constructor: PodConstructor, Destructor: Poddestructor, MaxSize: maxPoolSize})
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