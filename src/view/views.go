package view

import "context"

type View interface {
	Serve()
	Shutdown(context.Context)
}
