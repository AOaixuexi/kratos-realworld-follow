package service

import (
	"github.com/google/wire"

	v1 "helloworld/api/web/v1"
	"helloworld/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewWebService)

type WebService struct {
	v1.UnimplementedWebServer

	sc *biz.SocialUsecase
	uc *biz.UserUsecase
	log *log.Helper
}

func NewWebService(uc *biz.UserUsecase, sc *biz.SocialUsecase, logger log.Logger) *WebService {
	return &WebService{uc: uc, sc:sc, log: log.NewHelper(logger)}
}
