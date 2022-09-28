package main

import (
	"context"
	"github.com/go-redis/redis/v9"
	"strconv"
)

type RDB struct {
	RDB     *redis.Client
	Config  RedisConfig
	Context context.Context
}

func connectRDB(con RedisConfig) (RDB, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     con.Address + ":" + strconv.Itoa(con.Port),
		Username: con.Username,
		Password: con.Password,
		DB:       con.DataBase,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return RDB{}, err
	}
	return RDB{RDB: rdb, Config: con, Context: ctx}, nil
}

func (rdb RDB) addRule(rule Rule) {

}

func (rdb RDB) updateTicket(ticket *Ticket) {

}
