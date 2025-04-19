package pool

import (
	"context"

	"github.com/jackc/puddle/v2"
)

func K8sSandboxConstructor(ctx context.Context) (Pod, error) {
	return Pod{}, nil
}

func K8sSandboxDestructor(pod Pod) {

}

func NewK8sSandboxPool(ctx context.Context, initPoolSize int32, maxPoolSize int32) (*puddle.Pool[Pod], error) {
	pool, err := puddle.NewPool(&puddle.Config[Pod]{Constructor: K8sExecConstructor, Destructor: K8sExecDestructor, MaxSize: maxPoolSize})
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
