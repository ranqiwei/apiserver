package router

import (
	_ "apiserver/docs"
	"apiserver/handler/sd"
	"apiserver/handler/user"
	"apiserver/router/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
)

/*路由加载的函数*/

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	//pprof router
	pprof.Register(g)

	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	//404 Handler
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect api route")
	})

	//swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//login router
	g.POST("/login", user.Login)

	//user api
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware())
	{
		//u.POST("/:username", user.Create)
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	//The healthy check handler
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/memory", sd.RamCheck)
		svcd.GET("/cpu", sd.CPUCheck)
	}

	return g
}
