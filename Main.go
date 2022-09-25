package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	config, err := loadConfig("config.yml")
	if err != nil {
		panic(err)
	}
	serverConfig := config.ServerConfig
	sqlConfig := config.SQLConfig
	_, err = connectDB(sqlConfig)
	if err != nil {
		panic(err)
	}
	log.Println("log")
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	err = router.SetTrustedProxies([]string{})
	if err != nil {
		panic(err)
	}

	router.GET(serverConfig.RuleAddress, addRule)
	router.POST(serverConfig.RuleAddress, addRule)
	//router.GET(serverConfig.TicketAddress, )
	//router.GET(serverConfig.TicketAddress, )

	panic(router.Run("localhost:" + strconv.Itoa(serverConfig.Port)))
}
