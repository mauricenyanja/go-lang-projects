package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	baseURL = "http://localhost:808"
)

type LoadBalancer struct {
	RevProxy httputil.ReverseProxy
}

type EndPoints struct {
	List []*url.URL
}

func (e *EndPoints) Shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func MakeLoadBalancer(amount int) {
	// Instantiate Objects
	var lb LoadBalancer
	var ep Endpoints

	// server and router
	router := http.NewServeMux()
	server := http.Serve{
		Addr:    ":8090",
		Handler: router,
	}
	// Creating the endpoints
	for i := 0; i < amount; i++ {
		ep.List = append(ep.List, createEndpoint(baseURL, i))
	}

	// Handler Functions
	router.HandleFunc("/loadbalancer", makeRequest(&lb, &ep))

	// Listen and Server
	log.Fatal(server.ListenAndServe())
}



func makeRequest(lb *LoadBalancer,ep *EndPoints) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter, r *http.Request){
		for !testServer(ep.List[0].String()){
			ep.Shuffle()
		}
		lb.RevProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.Shuffle()
		lb.RevProxy.ServeHTTP(w,r)
	}
}
func createEndpoint(endpoint string,idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	url,_ :=url.Parse(link)
	return url
}
