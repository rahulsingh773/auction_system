package server

import (
	"auction_system/auctioneer/utils"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Auction struct {
	AuctionID string `json:"auction_id" validate:"nonzero"`
}

type Bidder struct {
	BidderID   string `json:"bidder_id"`
	BidderPort string `json:"bidder_port"`
}

type Bid struct {
	BidderID string  `json:"bidder_id"`
	BidValue float32 `json:"bid_value"`
}

var Bidders map[string]Bidder

//return swagger yml file
func ListAPI(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("../api-swagger.yml")
	status := http.StatusOK
	if nil != err {
		status = http.StatusNotFound
		log.Printf("API document api-swagger.yml not found")
	}
	w.WriteHeader(status)

	if nil == err {
		_, err := io.Copy(w, f)
		if nil != err {
			log.Printf("Failed to copy API document api-swagger.yml to HTTP response -\n\t%s", err)
		}
		f.Close()
	}
}

//starts auction, collect bids and return response
func StartAuction(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting Auction!!!")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("StartAuction: Error while extracting body err:%v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No body found"))
		return
	}

	var auctionData Auction
	err = json.Unmarshal(body, &auctionData)
	if err != nil {
		log.Printf("StartAuction: Error while extracting body, body:%v \t err:%v", body, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad JSON format: auction_id needed"))
		return
	}

	var wg sync.WaitGroup
	bids := make(map[string]float32) //map to store bids from bidders

	for _, bidder := range Bidders {
		wg.Add(1)
		go func(bidder Bidder) {
			defer wg.Done()
			bid_chan := make(chan Bid)
			go RequestBids(auctionData.AuctionID, bidder.BidderPort, bid_chan) //send request to bidder to get bids
			select {
			case bid := <-bid_chan:
				log.Printf("received bid from bidder: %v, value: %v", bid.BidderID, bid.BidValue)
				bids[bid.BidderID] = bid.BidValue //save the bid to map
			case <-time.After(200 * time.Millisecond): //wait only for 200 ms for bidder response
				log.Printf("bid failed from bidder: %v", bidder.BidderID)
			}
		}(bidder)
	}

	wg.Wait()
	log.Printf("Auction: bids: %#v", bids)

	winner_id, win_bid := findWinner(bids) //compare bids and return max bidding amount bidder
	if win_bid != -1 {
		var win_data = Bid{BidderID: winner_id, BidValue: win_bid}
		resp, _ := json.Marshal(win_data)
		w.Write([]byte(resp))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No bids available"))
	}

}

//register bidder details
func RegisterBidder(w http.ResponseWriter, r *http.Request) {
	log.Printf("Registering Bidder!!!")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("RegisterBidder: Error while extracting body err:%v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No body found"))
		return
	}

	var bidderData Bidder
	err = json.Unmarshal(body, &bidderData)
	if err != nil {
		log.Printf("RegisterBidder: Error while extracting body, body:%v \t err:%v", body, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad JSON format: bidder details needed"))
		return
	}

	if len(Bidders) == 0 {
		Bidders = make(map[string]Bidder, 0)
	}

	if _, ok := Bidders[bidderData.BidderID]; ok {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("bidder already exist"))
		return
	}

	Bidders[bidderData.BidderID] = bidderData //save the bidder's data
	log.Printf("Bidder Registered, data: %#v", Bidders)
	w.Write([]byte("Bidder Registered"))
}

//send request to bidder to respond with bids
func RequestBids(auction_id, bidder_port string, bid_chan chan<- Bid) {
	var bid_data = Auction{AuctionID: auction_id}
	data, _ := json.Marshal(bid_data)
	resp, code, err := utils.SendRestRequest(http.MethodPost, "http://localhost:"+bidder_port+"/bids", string(data))
	log.Printf("resp: %v, code: %v, err: %v", string(resp), code, err)
	if err == nil && code == http.StatusOK {
		var bid Bid
		json.Unmarshal(resp, &bid)
		bid_chan <- bid
	}
}

//compare bids and find winner
func findWinner(bids map[string]float32) (string, float32) {
	bidder_id, bid_value := "", float32(-1)
	for id, value := range bids {
		if bid_value == -1 || value > bid_value {
			bidder_id = id
			bid_value = value
		}
	}
	log.Printf("Winner: id: %v, value: %v", bidder_id, bid_value)
	return bidder_id, bid_value
}
