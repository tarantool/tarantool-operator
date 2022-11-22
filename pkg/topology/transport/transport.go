package transport

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type Transport interface {
	Exec(ctx context.Context, pod *v1.Pod, res any, lua string, args ...any) error
}
