package main

import (
	httpServer "github.com/asim/go-micro/plugins/server/http/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/gin-gonic/gin"

	r "sTest/router"
	/*定时任务初始化*/
	_ "sTest/pkg/timing"
)

var (
	ServerName = "sTest"
)

func main() {
	// create service
	srv := httpServer.NewServer(
		server.Name(ServerName),
		server.Address(":8080"),
	)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// register routers
	r.InitRouter(router)

	hd := srv.NewHandler(router)
	if err := srv.Handle(hd); err != nil {
		logger.Fatal(err)
	}

	// Create service
	service := micro.NewService(
		micro.Server(srv),
		micro.Registry(registry.NewRegistry()),
	)
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
