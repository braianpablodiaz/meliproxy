package middleware

import (
	"github.com/braianpablodiaz/meli-proxy/environment"
	"github.com/gin-gonic/gin"
	"github.com/braianpablodiaz/meli-proxy/controller"
	"github.com/go-redis/redis/v7"
	"time"
	"strconv"
	"log"
	"errors"
)

var envIpRateLimit int = 0
var envIpRateLimitPerSecond int = 0
var envIpPathRateLimit int = 0
var envIpPathRateLimitPerSecond int = 0
var envPathRateLimit int = 0
var envPathRateLimitPerSecod int = 0
var lockRetry int 


func Init(){
	setEnviroment(&envIpRateLimit, "IP_RATE_LIMIT", "10")
	setEnviroment(&envIpRateLimitPerSecond, "IP_RATE_LIMIT_PER_SECOND", "1")
	setEnviroment(&envIpPathRateLimit, "IP_PATH_RATE_LIMIT", "10")
	setEnviroment(&envIpPathRateLimitPerSecond, "IP_PATH_RATE_LIMIT_PER_SECOND", "1")
	setEnviroment(&envPathRateLimit, "PATH_RATE_LIMIT", "10")
	setEnviroment(&envPathRateLimitPerSecod, "PATH_RATE_LIMIT_PER_SECOND", "1")
	setEnviroment(&lockRetry, "LOCK_RETRY", "100")
}

func setEnviroment(value *int , key string , defaultValue string){
	_value, err := strconv.Atoi(environment.GetEnv(key, defaultValue))
	if err != nil {
        log.Fatal(err.Error())
	} else {
		*value = _value
	}
}

func RateLimitMiddleware(rdb *redis.Client, c *gin.Context) {
	requestModel := &controller.Request{ 
		Ip: c.ClientIP(),
		Path: c.Request.URL.EscapedPath(),
		Method: c.Request.Method,
	}

	allow := make(chan bool)
	go ipRateLimit(rdb, c , requestModel, allow)
	go ipPathRateLimit(rdb, c , requestModel, allow)
	go pathRateLimit(rdb, c , requestModel, allow)
	ipAllow, ipPathAllow, pathAllow := <-allow, <-allow, <-allow

	if ipAllow == false || ipPathAllow == false || pathAllow == false{
		c.AbortWithStatusJSON(429, gin.H{"error": "A lot of request per second"})
	}

	c.Next()
}


func RateLimitIpMiddleware(rdb *redis.Client, c *gin.Context) {
	
	requestModel := &controller.Request{ 
		Ip: c.ClientIP(),
		Path: c.Request.URL.EscapedPath(),
		Method: c.Request.Method,
	}

	allow := make(chan bool)
	go ipRateLimit(rdb, c , requestModel, allow)
	ipAllow := <-allow

	if ipAllow == false{
		c.AbortWithStatusJSON(429, gin.H{"error": "A lot of request per second"})
	}

	c.Next()
}

func RateLimitIpPathMiddleware(rdb *redis.Client, c *gin.Context) {
	
	requestModel := &controller.Request{ 
		Ip: c.ClientIP(),
		Path: c.Request.URL.EscapedPath(),
		Method: c.Request.Method,
	}

	allow := make(chan bool)
	go ipPathRateLimit(rdb, c , requestModel, allow)
	ipPathAllow := <-allow

	if ipPathAllow == false{
		c.AbortWithStatusJSON(429, gin.H{"error": "A lot of request per second"})
	}

	c.Next()
}


func RateLimitPathMiddleware(rdb *redis.Client, c *gin.Context) {
	
	requestModel := &controller.Request{ 
		Ip: c.ClientIP(),
		Path: c.Request.URL.EscapedPath(),
		Method: c.Request.Method,
	}

	allow := make(chan bool)
	go pathRateLimit(rdb, c , requestModel, allow)
	pathAllow := <-allow

	if pathAllow == false{
		c.AbortWithStatusJSON(429, gin.H{"error": "A lot of request per second"})
	}

	c.Next()
}



func ipRateLimit(rdb *redis.Client, c *gin.Context, 
	requestModel *controller.Request ,allow chan bool) {

	key := requestModel.Ip
	allow <- requestIsAllow(rdb, key, envIpRateLimit , envIpRateLimitPerSecond)
}

func ipPathRateLimit(rdb *redis.Client, c *gin.Context, 
	requestModel *controller.Request ,allow chan bool) {

	key := requestModel.Ip +  requestModel.Path
	allow <- requestIsAllow(rdb, key, envIpPathRateLimit , envIpPathRateLimitPerSecond)
}

func pathRateLimit(rdb *redis.Client, c *gin.Context,
	requestModel *controller.Request, allow chan bool) {

	key := requestModel.Path
	allow <- requestIsAllow(rdb, key, envPathRateLimit, envPathRateLimitPerSecod)
}



func requestIsAllow (rdb *redis.Client, key string, ratelimit int , limitPerSecond int)  bool{

	txf := func(tx *redis.Tx) error {

		n, err := tx.Get(key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if n > ratelimit {
			return errors.New("A lot of request per second")
		}
		
		_, err = tx.TxPipelined(func(pipe redis.Pipeliner) error {
			if n == 0 {
				pipe.Set(key, 1 , time.Duration(limitPerSecond) * time.Second )
			} else {
				pipe.Incr(key)
			}
			return nil
		})
		return err
	}

	for i := 0; i < lockRetry ; i++ {
		err := rdb.Watch(txf, key)
		if err == nil {
			return true
		}
		if err == redis.TxFailedErr {
			continue
		}

		return false
	}

	return false

}