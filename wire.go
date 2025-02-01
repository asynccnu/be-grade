//go:generate wire
//go:build wireinject

package main

import (
	"github.com/asynccnu/be-grade/cron"
	"github.com/asynccnu/be-grade/grpc"
	"github.com/asynccnu/be-grade/ioc"
	"github.com/asynccnu/be-grade/repository/dao"
	"github.com/asynccnu/be-grade/service"
	"github.com/google/wire"
)

func InitApp() App {
	wire.Build(
		grpc.NewGradeGrpcService,
		service.NewGradeService,
		dao.NewGradeDAO,
		// 第三方
		ioc.InitEtcdClient,
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitGRPCxKratosServer,
		ioc.InitUserClient,
		ioc.InitCounterClient,
		ioc.InitFeedClient,
		ioc.InitClasslistClient,
		cron.NewGradeController,
		cron.NewCron,
		NewApp,
	)
	return App{}
}
