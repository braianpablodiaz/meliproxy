package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v7"
    "io/ioutil"
    "log"
    "net/http"
    "encoding/json"
)

type RequestHandlerFunction func(repository *redis.Client, c *gin.Context , api string)


type Request struct {
	Ip string
	Path string
    Method string
    Api string
}

func Proxy(r *redis.Client, c *gin.Context, api string) {

    requestModel := Request{ 
        Ip: c.ClientIP(),
        Path: c.Request.URL.EscapedPath(),
        Method: c.Request.Method,
    }

    client := http.Client{ }
    request, err := http.NewRequest(requestModel.Method, api  + requestModel.Path , c.Request.Body)
    
	if err != nil {
        errorResponse(c, err.Error())
    }

    response, err := client.Do(request)

    if err != nil {
        errorResponse(c, err.Error())
    }

    defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	
    if err != nil {
        errorResponse(c, err.Error())
	}
    
    var jsonResponse map[string]interface{}
    if err := json.Unmarshal([]byte( string(responseData) ), &jsonResponse); err != nil {
        errorResponse(c, err.Error())
    }
    
    c.JSON(response.StatusCode, jsonResponse)
}


func errorResponse(c *gin.Context, err string) {
    log.Print(err)
    c.AbortWithStatusJSON(500, gin.H{"error": err})
}
