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
	fileHandler handler.FileHandler,
	communityHandler handler.CommunityHandler,
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
			noAuthRouter.POST("/verificationCode", userHandler.VerificationCode)
			//可通过 http://localhost:8080/static/css/style.css 访问 ./xx/css/style.css 文件。
			noAuthRouter.Static("/static", conf.GetString("assets.dir")) // 将 ./xx 目录映射到 /static 路径
		}
		// Non-strict permission routing group
		noStrictAuthRouter := v1.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
		{
			noStrictAuthRouter.GET("/user", userHandler.GetProfile)
			noStrictAuthRouter.POST("/search", userHandler.SearchProfile)
		}

		// Strict permission routing group
		strictAuthRouter := v1.Group("/").Use(middleware.StrictAuth(jwt, logger))
		{
			strictAuthRouter.PUT("/user", userHandler.UpdateProfile)
			strictAuthRouter.POST("/upload", fileHandler.UploadImage)
		}

		relationGroup := v1.Group("/relationship").Use(middleware.StrictAuth(jwt, logger))
		{
			//好友关系申请
			relationGroup.POST("/apply/add", relationHandler.AddApplyFriendship)
			relationGroup.GET("/apply/list", relationHandler.GetApplyFriendshipList)
			relationGroup.PUT("/apply/edit", relationHandler.UpdateApplyFriendshipInfo)
			relationGroup.DELETE("/apply/del", relationHandler.DelApplyFriendshipInfo)

			//关系相关
			relationGroup.POST("/relation/list", relationHandler.GetRelationshipList)
			relationGroup.PUT("/relation/edit", relationHandler.UpdateRelationship)
			relationGroup.DELETE("/relation/del", relationHandler.DelRelationship)
			relationGroup.POST("/relation/add/follow", relationHandler.AddRelationshipFollow) // 添加关注
		}

		chatGroup := v1.Group("/chat").Use(middleware.StrictAuth(jwt, logger))
		{
			chatGroup.POST("/send", chatHandler.SendChatMessage)
			//同步会话
			chatGroup.POST("/conversation/list", chatHandler.GetUserConversationList)
			//会话的所有用户
			chatGroup.POST("/conversation/users", chatHandler.GetConversationUserList)
			//同步历史消息
			chatGroup.POST("/msg/history/list", chatHandler.GetUserMsgList) //用户消息链
			chatGroup.POST("/msg/list", chatHandler.GetConversationMsgList) //会话消息链

			chatGroup.POST("/report/msg/read", chatHandler.ReportReadMsgSeq)
		}

		//社区相关

		communityGroup := v1.Group("/community")
		{
			// Non-strict permission routing group
			noStrictCommunityGroup := communityGroup.Group("/").Use(middleware.NoStrictAuth(jwt, logger))
			{
				noStrictCommunityGroup.POST("/moment/list", communityHandler.GetMomentList)
				noStrictCommunityGroup.POST("/comment/list", communityHandler.GetMomentCommentList)
			}

			// Strict permission routing group
			strictCommunityGroup := communityGroup.Group("/").Use(middleware.StrictAuth(jwt, logger))
			{
				strictCommunityGroup.POST("/moment/add", communityHandler.AddMoment)
				strictCommunityGroup.POST("/moment/edit", communityHandler.EditMoment)
				strictCommunityGroup.POST("/moment/like", communityHandler.AddMomentLike)

				strictCommunityGroup.POST("/comment/add", communityHandler.AddMomentComment)
				strictCommunityGroup.POST("/comment/like", communityHandler.LikeMomentComment)
			}
		}
	}

	return s
}
