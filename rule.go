package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	"net/http"
)

type Route struct {
	Origin      sql.NullString `json:"origin"`
	Destination sql.NullString `json:"destination"`
}

type Rule struct {
	Routes    []Route  `json:"routes"`
	Airlines  []string `json:"airlines"`
	Agencies  []string `json:"agencies"`
	Suppliers []string `json:"suppliers"`
	Type      string   `json:"amountType"`
	Value     int      `json:"amountValue"`
}

type RuleJ struct {
	Routes    types.JSONText `json:"routes"      redis:"routes"`
	Airlines  types.JSONText `json:"airlines"    redis:"airlines"`
	Agencies  types.JSONText `json:"agencies"    redis:"agencies"`
	Suppliers types.JSONText `json:"suppliers"   redis:"suppliers"`
	Type      string         `json:"amountType"  redis:"type"`
	Value     int            `json:"amountValue" redis:"value"`
}

func (ruleJ RuleJ) Unmarshal(rule *Rule) error {
	var err error
	*rule = Rule{}

	err = ruleJ.Routes.Unmarshal(&rule.Routes)
	if err != nil {
		return err
	}

	err = ruleJ.Airlines.Unmarshal(&rule.Airlines)
	if err != nil {
		return err
	}

	err = ruleJ.Agencies.Unmarshal(&rule.Agencies)
	if err != nil {
		return err
	}

	err = ruleJ.Suppliers.Unmarshal(&rule.Suppliers)
	if err != nil {
		return err
	}

	rule.Type = ruleJ.Type
	rule.Value = ruleJ.Value

	return nil
}

func addRule(c *gin.Context, db DB, rdb RDB) {
	var rules []RuleJ
	err := c.BindJSON(&rules)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "FAILED", "message": err.Error()})
		return
	}
	for _, ruleJ := range rules {
		var rule Rule
		err = ruleJ.Unmarshal(&rule)
		if err != nil {
			continue
		}
		if !db.validateRule(rule) {
			continue
		}
		i := db.addRule(ruleJ)
		if i != -1 {
			rdb.addRule(i, ruleJ)
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"status": "SUCCESS", "message": nil})
}
