// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"helloworld/internal/biz"
	"helloworld/internal/conf"
	"helloworld/internal/data"
	"helloworld/internal/server"
	"helloworld/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, jwt *conf.JWT, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(dataData, logger)
	profileRepo := data.NewProfileRepo(dataData, logger)
	userUsecase := biz.NewUserUsecase(userRepo, profileRepo, logger, jwt)
	articleRepo := data.NewArticleRepo(dataData, logger)
	commentRepo := data.NewCommentRepo(dataData, logger)
	socialUsecase := biz.NewSocialUsecase(articleRepo, profileRepo, commentRepo, logger)
	webService := service.NewWebService(userUsecase, socialUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, webService, logger)
	httpServer := server.NewHTTPServer(confServer, jwt, webService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
