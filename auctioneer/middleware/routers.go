package middleware

import (
	"auction_system/auctioneer/server"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

//return and register routes with handlers
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

//middleware to dump incoming requests
func middleware(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_dump, _ := httputil.DumpRequest(r, true)
		log.Printf("Request Dump: %v", string(req_dump))

		inner.ServeHTTP(w, r)
	})
}

//register routes and corresponding handlers
var routes = Routes{
	Route{
		"ListAPI",
		"GET",
		"/api",
		server.ListAPI,
	},
	Route{
		"StartAuction",
		"POST",
		"/auction",
		server.StartAuction,
	},
	Route{
		"RegisterBidder",
		"POST",
		"/bidder",
		server.RegisterBidder,
	},
}
