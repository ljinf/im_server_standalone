//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/ljinf/im_server_standalone/internal/server"
	"github.com/ljinf/im_server_standalone/pkg/app"
	"github.com/ljinf/im_server_standalone/pkg/log"
	"github.com/spf13/viper"
)

var serverSet = wire.NewSet(
	server.NewTask,
)

// build App
func newApp(
	task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		serverSet,
		newApp,
	))
}
