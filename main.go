package main

import (
	"apiserver/config"
	"apiserver/router"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"apiserver/model"
	"apiserver/router/middleware"
	"github.com/lexkong/log"
	"net/http"
	"time"
)

var (
	//配置文件
	cfg = pflag.StringP("config", "c", "", "api config file path")
)

func main() {
	//viper初始化调用
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	//int db
	model.DB.Init()
	defer model.DB.Close()

	//Set gin mode使用
	gin.SetMode(viper.GetString("runmode"))
	//Set the Gin engine
	g := gin.New()
	//外部加入的中间件
	middlewares := []gin.HandlerFunc{middleware.RequestId(), middleware.Logging()}
	//路由加载
	router.Load(
		g,
		middlewares...,
	)

	//router健康检查
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Infof("The router has been deployed successfully.")
	}()

	//start https listen if certificate exists
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if key != "" && cert != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Infof(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Infof(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

//pingServer pings the http server to make sure the routers is working
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		/*runtime.Breakpoint()
		debug.PrintStack()*/
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Infof("waiting for the router, retry in 1 second.")
		time.Sleep(5 * time.Second)
	}
	return errors.New("can't connect to the route")
}
