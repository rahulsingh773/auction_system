package main

import (
	"auction_system/bidder/config"
	"auction_system/bidder/middleware"
	"auction_system/bidder/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Bidder struct {
	BidderID   string `json:"bidder_id"`
	BidderPort string `json:"bidder_port"`
}

var BidderID = utils.GetRandomString(7)

func main() {
	fmt.Println("------------------bidder started----------------")
	config.SetConfigParam("bidder_id", BidderID)

	router := middleware.NewRouter()

	port := config.GetConfigParamString("bind_port")
	RegisterWithAuctioner(port)

	http.ListenAndServe(":"+port, router)
}

func RegisterWithAuctioner(port string) {
	fmt.Println("---------------Registering With Auctioner------------------")
	bidder := Bidder{
		BidderID:   BidderID,
		BidderPort: port,
	}
	data, err := json.Marshal(bidder)
	if err != nil {
		log.Printf("RegisterWithAuctioner: data marshalling failed: %v", err)
	}
	utils.SendRestRequest(http.MethodPost, config.GetConfigParamString("auctioneer_url"), string(data))
}
