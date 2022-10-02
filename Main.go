package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"log"
	"strconv"
)

type a struct {
	A []string `json:"a"`
}
type b struct {
	B types.JSONText `json:"b"`
}

func main() {
	config, err := loadConfig("config.yml")
	if err != nil {
		panic(err)
	}
	serverConfig := config.ServerConfig
	sqlConfig := config.SQLConfig
	redisConfig := config.RedisConfig
	var rdb RDB
	var db DB
	rdb, err = connectRDB(redisConfig)
	if err != nil {
		panic(err)
	}
	db, err = connectDB(sqlConfig, rdb)
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

	router.GET(serverConfig.RuleAddress, func(ctx *gin.Context) {
		addRule(ctx, db, rdb)
	})
	router.POST(serverConfig.RuleAddress, func(ctx *gin.Context) {
		addRule(ctx, db, rdb)
	})

	router.GET(serverConfig.TicketAddress, func(ctx *gin.Context) {
		updateTicket(ctx, rdb)
	})
	router.POST(serverConfig.TicketAddress, func(ctx *gin.Context) {
		updateTicket(ctx, rdb)
	})

	panic(router.Run("localhost:" + strconv.Itoa(serverConfig.Port)))
}
