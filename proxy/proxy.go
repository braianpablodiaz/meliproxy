package proxy

import (
	"github.com/braianpablodiaz/meli-proxy/environment"
	"github.com/gin-gonic/gin"
	"github.com/braianpablodiaz/meli-proxy/controller"
	"github.com/braianpablodiaz/meli-proxy/repository"
	"github.com/braianpablodiaz/meli-proxy/middleware"
	"github.com/go-redis/redis/v7"
)


type Proxy struct {
	Router *gin.Engine
	Repository *redis.Client
	ApiClient string
}

func NewProxy() *Proxy{
	return &Proxy{}
}

func (proxy *Proxy) StartProxy() { 
	proxy.configRouter()
	proxy.configEndPoint()
	proxy.configApiclient()
	proxy.configRedis()
	proxy.Router.Run( environment.GetEnv( "PROXY_PORT" , ":8080") )
}

func (proxy *Proxy) configRouter() {
	router := gin.Default()
	//router.Use(proxy.rateLimitMiddleware())
	router.Use(proxy.rateLimitIpMiddleware())
	router.Use(proxy.rateLimitIpPathMiddleware())
	router.Use(proxy.rateLimitPathMiddleware())
	router.Use(proxy.configHeaders())
	proxy.Router = router

	middleware.Init()
}

func (proxy *Proxy) configEndPoint() {
	proxy.Router.Any("/*proxyPath", proxy.handleRequest(controller.Proxy))
}

func (proxy *Proxy) configApiclient() {
	proxy.ApiClient = environment.GetEnv("API_MERCADO_LIBRE", "")
}

func (proxy *Proxy) configRedis() {
	proxy.Repository = repository.InitialConnection()
}

func (proxy *Proxy) configHeaders() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func (proxy *Proxy) rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.RateLimitMiddleware(proxy.Repository, c )
	}
}

func (proxy *Proxy) rateLimitIpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.RateLimitIpMiddleware(proxy.Repository, c )
	}
}

func (proxy *Proxy) rateLimitIpPathMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.RateLimitIpPathMiddleware(proxy.Repository, c )
	}
}

func (proxy *Proxy) rateLimitPathMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware.RateLimitPathMiddleware(proxy.Repository, c )
	}
}

func (proxy *Proxy) handleRequest(handler controller.RequestHandlerFunction) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(proxy.Repository, c , proxy.ApiClient)
	}
}
