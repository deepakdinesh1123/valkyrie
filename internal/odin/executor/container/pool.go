package container

import (
	"context"

	"github.com/jackc/puddle/v2"
)

func constructor(ctx context.Context) (Container, error) {
	return Container{}, nil
}

func destructor(value Container) {

}

func NewContainerPool(ctx context.Context, poolSize int32) (*puddle.Pool[Container], error) {
	pool, err := puddle.NewPool(&puddle.Config[Container]{Constructor: constructor, Destructor: destructor, MaxSize: poolSize})
	if err != nil {
		return nil, err
	}
	return pool, nil
}
