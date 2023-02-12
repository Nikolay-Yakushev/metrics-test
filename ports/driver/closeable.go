package ports

import "context"

type Closeable interface {
	Stop(ctx context.Context)error
	GetDescription()string
	GetCloseOrder()int
}