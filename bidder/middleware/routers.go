package middleware

import (
	"auction_system/bidder/server"
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

//dump incoming request
func middleware(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req_dump, _ := httputil.DumpRequest(r, true)
		log.Printf("Request Dump: %v", string(req_dump))

		inner.ServeHTTP(w, r)
	})
}

var routes = Routes{
	Route{
		"PlaceBids",
		"POST",
		"/bids",
		server.PlaceBids,
	},
}
