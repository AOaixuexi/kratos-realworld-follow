package server

import (
	"context"
	v1 "helloworld/api/web/v1"
	"helloworld/internal/conf"
	"helloworld/internal/pkg/middleware"
	"helloworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// 白名单
func NewSkipRoutersMatcher() selector.MatchFunc {

	skipRouters := map[string]struct{}{
		"/web.v1.Web/Login":    {},
		"/web.v1.Web/Register": {},
	}

	return func(ctx context.Context, operation string) bool {
		log.Infof("Checking operation: %s", operation)
		if _, ok := skipRouters[operation]; ok {
			log.Infof("Operation %s is in skip list", operation)
			return false
		}
		log.Infof("Operation %s is not in skip list", operation)
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, jwtc *conf.JWT, web *service.WebService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server(middleware.JWTAuth(jwtc.Secret)).Match(NewSkipRoutersMatcher()).Build(),
			logging.Server(logger),
		),
		http.Filter(
			handlers.CORS(
				// 指定允许的请求头
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				// 指定允许的请求方法
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
				// 指定允许的源
				handlers.AllowedOrigins([]string{"*"}),
				//postman请求时，会先发送OPTIONS请求，这里设置为true，表示允许OPTIONS请求
			),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterWebHTTPServer(srv, web)
	return srv
}
