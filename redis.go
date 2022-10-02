package main

import (
	"context"
	"database/sql"
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

func (rdb RDB) addRule(i int, rule RuleJ) {
	id := strconv.Itoa(i)
	rdb.RDB.Pipelined(rdb.Context, func(rpl redis.Pipeliner) error {
		rpl.HSetNX(rdb.Context, id, "routes", rule.Routes)
		rpl.HSetNX(rdb.Context, id, "airlines", rule.Airlines)
		rpl.HSetNX(rdb.Context, id, "agencies", rule.Agencies)
		rpl.HSetNX(rdb.Context, id, "suppliers", rule.Suppliers)
		rpl.HSetNX(rdb.Context, id, "type", rule.Type)
		rpl.HSetNX(rdb.Context, id, "value", rule.Value)
		return nil
	})
}

func (rdb RDB) updateTicket(ticket *Ticket) {
	keys, err := rdb.RDB.HKeys(rdb.Context, "").Result()
	if err != nil {
		return
	}
	set := false
	route := Route{
		Origin:      sql.NullString{String: ticket.Origin, Valid: true},
		Destination: sql.NullString{String: ticket.Destination, Valid: true},
	}
	for _, id := range keys {
		var ruleJ RuleJ
		var rule Rule
		err = rdb.RDB.HGetAll(rdb.Context, id).Scan(&ruleJ)
		if err != nil {
			return
		}
		err = ruleJ.Unmarshal(&rule)
		if err != nil {
			return
		}
		if !checkRouteSlice(rule.Routes, route) {
			continue
		}
		if !checkStrSlice(rule.Airlines, ticket.Airline) {
			continue
		}
		if !checkStrSlice(rule.Agencies, ticket.Agency) {
			continue
		}
		if !checkStrSlice(rule.Suppliers, ticket.Supplier) {
			continue
		}
		if !set {
			ticket.Markup, _ = strconv.Atoi(id)
			if rule.Type == "PERCENTAGE" {
				ticket.PayablePrice = ticket.BasePrice * ((100 + rule.Value) / 100)
			} else {
				ticket.PayablePrice = ticket.BasePrice + rule.Value
			}
			set = true
		} else {
			var val int
			if rule.Type == "PERCENTAGE" {
				val = ticket.BasePrice * ((100 + rule.Value) / 100)
			} else {
				val = ticket.BasePrice + rule.Value
			}
			if val > ticket.PayablePrice {
				ticket.PayablePrice = val
				ticket.Markup, _ = strconv.Atoi(id)
			}
		}
	}
}

func checkStrSlice(slice []string, str string) bool {
	if slice == nil {
		return true
	}
	for _, sliceStr := range slice {
		if sliceStr == str {
			return true
		}
	}
	return false
}

func checkRouteSlice(slice []Route, r Route) bool {
	if slice == nil {
		return true
	}
	for _, sliceRoute := range slice {
		if checkRoute(sliceRoute, r) {
			return true
		}
	}
	return false
}

func checkRoute(r1 Route, r2 Route) bool {
	if r1.Origin.Valid && r1.Origin.String != r2.Origin.String {
		return false
	}
	if r1.Destination.Valid && r1.Destination.String != r2.Destination.String {
		return false
	}
	return true
}
