package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Ticket struct {
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	Airline      string `json:"airline"`
	Agency       string `json:"agency"`
	Supplier     string `json:"supplier"`
	BasePrice    int    `json:"basePrice"`
	Markup       int    `json:"markup"`
	PayablePrice int    `json:"payablePrice"`
}

func updateTicket(c *gin.Context, rdb RDB) {
	var tickets []Ticket
	err := c.BindJSON(&tickets)
	if err != nil {
		return
	}
	for i := range tickets {
		rdb.updateTicket(&tickets[i])
	}
	c.IndentedJSON(http.StatusOK, tickets)
}
