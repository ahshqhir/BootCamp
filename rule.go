package main

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

type Rule struct {
	Routes    []Route  `json:"routes"`
	Airlines  []string `json:"airlines"`
	Agencies  []string `json:"agencies"`
	Suppliers []string `json:"suppliers"`
	Type      string   `json:"amountType"`
	Value     string   `json:"amountValue"`
}

func addRule(c *gin.Context, db DB, rdb RDB) {
	var rules []Rule
	err := c.BindJSON(&rules)
	if err != nil {
		return
	}
	for _, rule := range rules {
		db.addRule(rule)
		rdb.addRule(rule)
	}
}
