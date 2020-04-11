package server

import (
	"auction_system/bidder/config"
	"auction_system/bidder/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Auction struct {
	AuctionID string `json:"auction_id"`
}

type BidResponse struct {
	BidderID string  `json:"bidder_id"`
	BidValue float32 `json:"bid_value"`
}

//respond with bid amount after some delay configured
func PlaceBids(w http.ResponseWriter, r *http.Request) {
	log.Printf("Placing Bid!!!")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("PlaceBids: Error while extracting body err:%v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No body found"))
		return
	}

	var auctionData Auction
	err = json.Unmarshal(body, &auctionData)
	if err != nil {
		log.Printf("PlaceBids: Error while extracting body, body:%v \t err:%v", body, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad JSON format: auction_id needed"))
		return
	}

	bid_delay_str := config.GetConfigParamString("bid_delay")
	bid_delay, _ := strconv.Atoi(bid_delay_str)
	time.Sleep(time.Duration(bid_delay) * time.Millisecond) //sleep for delay

	var Bid = BidResponse{BidderID: config.GetConfigParamString("bidder_id"), BidValue: utils.GetRandomNumber()}
	resp, err := json.Marshal(Bid)
	if err != nil {
		log.Printf("PlaceBids: marshalling failed for response, err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	w.Write(resp)
	log.Printf("Bid Placed!!!")
}
