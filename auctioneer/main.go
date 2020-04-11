package main

import (
	"auction_system/auctioneer/middleware"
	"fmt"
	"net/http"
)

const Bind_Port = "3000"

func main() {
	fmt.Println("------------------auctioneer started----------------")

	router := middleware.NewRouter()
	http.ListenAndServe(":"+Bind_Port, router) //listen to port
}
