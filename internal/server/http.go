package server

import (
	"github.com/gin-gonic/gin"
	apiV1 "github.com/ljinf/im_server_standalone/api/v1"
	"github.com/ljinf/im_server_standalone/docs"
	"github.com/ljinf/im_server_standalone/internal/handler"
	"github.com/ljinf/im_server_standalone/internal/middleware"
	"github.com/ljinf/im_server_standalone/pkg/jwt"
	"github.com/ljinf/im_server_standalone/pkg/log"
	"github.com/ljinf/im_server_standalone/pkg/server/http"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewHTTPServer(
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	userHandler *handler.UserHandler,
	wsHandler handler.WebSocketHandler,
	relationHandler *handler.RelationshipHandler,
	chatHandler *handler.ChatHandler,
) *http.Server {
	gin.SetMode(gin.DebugMode)
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	// swagger doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		//ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	s.Use(
		middleware.CORSMiddleware(),
		//middleware.ResponseLogMiddleware(logger),
		//middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)
	s.GET("/", func(ctx *gin.Context) {
		logger.WithContext(ctx).Info("hello")
		apiV1.HandleSuccess(ctx, map[string]interface{}{
			":)": "Thank you for using nunu!",
		})
	})

	s.GET("/ws", wsHandler.AcceptConn)

	v1 := s.Group("/v1")
	{
		// No route group has permission
		noAuthRouter := v1.Group("/")
		{
			noAuthRouter.POST("/register", userHandler.Register)
			noAuthRouter.POST("/login", userHandler.Login)
		}
		// Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
		}

		relationGroup := v1.Group("/relationship").Use(middleware.StrictAuth(jwt, logger))
		{
			//好友关系申请
			relationGroup.GET("/apply/add", relationHandler.AddApplyFriendship)
			relationGroup.GET("/apply/list", relationHandler.GetApplyFriendshipList)
			relationGroup.PUT("/apply/edit", relationHandler.UpdateApplyFriendshipInfo)
			relationGroup.DELETE("/apply/del", relationHandler.DelApplyFriendshipInfo)

			//关系相关
			relationGroup.GET("/relation/list", relationHandler.GetRelationshipList)
			relationGroup.PUT("/relation/edit", relationHandler.UpdateRelationship)
			relationGroup.DELETE("/relation/del", relationHandler.DelRelationship)
			relationGroup.POST("/relation/add/follow", relationHandler.AddRelationshipFollow)
		}

		chatGroup := v1.Group("/chat").Use(middleware.StrictAuth(jwt, logger))
		{
			chatGroup.POST("/send", chatHandler.SendChatMessage)
			chatGroup.GET("/conversation/list", chatHandler.GetUserConversationList)
			chatGroup.POST("/msg/history/list", chatHandler.GetUserMsgList)
			chatGroup.POST("/report/msg/read", chatHandler.ReportReadMsgSeq)
		}
	}

	return s
}
