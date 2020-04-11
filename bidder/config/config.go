package config

import "os"

var Config map[string]interface{}

func init() {
	Config = make(map[string]interface{})

	Config["auctioneer_url"] = os.Getenv("AUCTIONEER_URL") //auctioneer_url value
	Config["bind_port"] = os.Getenv("BIND_PORT")           //bind_port to listen
	Config["bid_delay"] = os.Getenv("BID_DELAY")           //bind delay to return response
}

func GetConfigParamString(param string) string {
	return Config[param].(string)
}

func SetConfigParam(param, value string) {
	Config[param] = value
}
