package main

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

var Current_Price = 100.0

type order struct {
	Price  float64 `json:"price"`
	Amount int16   `json:"amount"`
	Stock  string  `json:"stock"`
	Type   string  `json:"type"`
}

/*
Types of order
MARKET_BUY
MARKET_SELL
LIMIT_BUY
LIMIT_SELL
*/

// var complete_orders = []order{}
var pending_orders = []order{}

func getPrice(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Current_Price)
}

func buyProcess(newOrder order) {
	var shares_Bought int = 0
	var executedOrders = []order{}
	order_completed := false

	for {
		for index, element := range pending_orders {
			if element.Type == "MARKET_SELL" || element.Type == "LIMIT_SELL" {
				if newOrder.Amount >= element.Amount {
					shares_Bought += int(element.Amount)
					newOrder.Amount -= int16(shares_Bought)
					pending_orders = append(pending_orders[:index], pending_orders[index+1:]...)

					order_completed = true
					//remove element from the pending orders array
				} else {
					shares_Bought = int(newOrder.Amount)
					element.Amount -= int16(shares_Bought)
					newOrder.Amount = 0
				}
				//need to save the different prices that the order executes at and ensure that all the shares are bought
				executedOrder := order{
					Price:  element.Price,
					Amount: element.Amount,
					Stock:  newOrder.Stock,
					Type:   newOrder.Type,
				}
				executedOrders = append(executedOrders, executedOrder)
				//then the user has bought it from this person
			}

		}
		if order_completed {
			break
		}
	}
	return executedOrders
}

func buyStock(c *gin.Context) {
	var newOrder order

	if err := c.BindJSON(&newOrder); err != nil {
		return
	}

	switch newOrder.Type {
	case "MARKET_BUY":
		//do shite
		pending_orders = append(pending_orders, newOrder)
		sort.Slice(pending_orders, func(i, j int) bool {
			return pending_orders[i].Price < pending_orders[j].Price
		})
		break

	case "LIMIT_BUY":
		//do more shit
	default:
		println("FAIL BUY ORDER")

	}

	//complete_orders = append(complete_orders, newOrder)
	c.IndentedJSON(http.StatusCreated, newOrder)

}

func main() {

	router := gin.Default()
	router.GET("/price", getPrice)
	router.POST("/buy", buyStock)

	router.Run("localhost:8080")
	fmt.Println("Hello world")
}
