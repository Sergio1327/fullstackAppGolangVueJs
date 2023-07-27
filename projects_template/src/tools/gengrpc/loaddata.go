package gengrpc

import (
	"projects_template/internal/entity/global"

	"google.golang.org/grpc/status"
)

type LoadDataFunc[T, K any] func() (T, error)
type PrepareDataFunc[T, K any] func(t T) K

func LoadData[T, K any](loadFunc LoadDataFunc[T, K], prepareFunc PrepareDataFunc[T, K]) (k K, err error) {
	var t T

	if t, err = loadFunc(); err != nil {
		if errStatus, ok := status.FromError(err); ok {
			if errStatus.Message() == global.ErrNoData.Error() {
				err = global.ErrNoData
				return
			}
		}
		return
	}

	k = prepareFunc(t)
	return
}
