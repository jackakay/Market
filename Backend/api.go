package main

import (
	"fmt"
	"net/http"

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

func marketBuyProcess(newOrder order) []order {
	var shares_needed int = int(newOrder.Amount)
	var shares_Bought int = 0
	var executedOrders = []order{}

	for {
		//attempting to buy shares from somebody who has market sell or limit sell orders
		for index, element := range pending_orders {

			if element.Type == "MARKET_SELL" || element.Type == "LIMIT_SELL" {
				//if the user is buying more shares than this order has available
				if newOrder.Amount > element.Amount {
					executedOrder := order{
						Price:  element.Price,
						Amount: element.Amount,
						Type:   newOrder.Type,
						Stock:  newOrder.Stock,
					}
					executedOrders = append(executedOrders, executedOrder)
					shares_Bought += int(element.Amount)
					fmt.Println(" 1. Bought ", element.Amount, " shares at price ", element.Price, ". ", shares_needed-shares_Bought, " shares remaining")

					newOrder.Amount -= int16(element.Amount)

					element.Amount = 0
					pending_orders = append(pending_orders[:index], pending_orders[index+1:]...)

					//remove element from the pending orders array
				} else {
					executedOrder := order{
						Price:  element.Price,
						Amount: newOrder.Amount,
						Type:   newOrder.Type,
						Stock:  newOrder.Stock,
					}
					executedOrders = append(executedOrders, executedOrder)
					shares_Bought += int(newOrder.Amount)
					fmt.Println(" 2. Bought ", (element.Amount - newOrder.Amount), " shares at price ", element.Price, ". ", shares_needed-shares_Bought, " shares remaining")

					element.Amount -= int16(newOrder.Amount)

					newOrder.Amount = 0
				}

				//need to save the different prices that the order executes at and ensure that all the shares are bought

				//then the user has bought it from this person
			}
			if newOrder.Amount == 0 || shares_Bought == shares_needed {
				break
			}
		}
		if newOrder.Amount == 0 || shares_Bought == shares_needed {
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
		//pending_orders = append(pending_orders, newOrder)
		//sort.Slice(pending_orders, func(i, j int) bool {
		//return pending_orders[i].Price < pending_orders[j].Price
		//})
		var executedOrders = marketBuyProcess(newOrder)
		c.IndentedJSON(http.StatusOK, executedOrders)

	case "LIMIT_BUY":
		//do more shit
	default:
		println("FAIL BUY ORDER")

	}

	//complete_orders = append(complete_orders, newOrder)
	//c.IndentedJSON(http.StatusCreated, newOrder)

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func populateOrders() []order {
	startingprice := 95
	orders := []order{}
	for i := 0; i < 10; i++ {
		newOrder := order{
			Price:  float64(startingprice + i),
			Amount: 10,
			Stock:  "JD",
			Type:   "LIMIT_SELL",
		}
		orders = append(orders, newOrder)
	}
	return orders
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/price", getPrice)
	router.POST("/buy", buyStock)

	pending_orders = append(pending_orders, populateOrders()...)

	router.Run("localhost:8080")
	fmt.Println("Hello world")
}

/*
two ways of making this work could either be figure out how to do asynchronous requests so the user doesnt freeze
or what i think is better have the user submit buy or sell orders to the server, which will handle everything, but the user will have to be
constantly trying to communicate with the server the get this information
*/
