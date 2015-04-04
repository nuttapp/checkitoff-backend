package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/common"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	nsqCfg := nsq.NewConfig()
	nsqCfg.MaxInFlight = 5

	conCfg := &common.ConsumerConfig{
		Topic:           "api01",
		Channel:         "client_events",
		LookupdHTTPaddr: "127.0.0.1:4161",
		Concurrency:     10,
	}

	handler := &APIMessageHandler{}
	mc := common.NewMessageConsumer(conCfg, nsqCfg, handler)
	mc.Start()

	for {
		select {
		// On quit signal stop the api consumer
		case <-signalChan:
			fmt.Println("")
			mc.Stop()
		case err := <-mc.StopChan:
			if err != nil {
				fmt.Printf("MessageConsumer stopped with error: %s", err)
			}
			fmt.Println("END")
			return
		}
	}
}

type APIMessageHandler struct{}

func (mh *APIMessageHandler) HandleMessage(m *nsq.Message) error {
	// Check if message is related to a connected client
	// send response to client
	fmt.Println(string(m.Body))
	return nil
}

// func main() {
// 	log.Println("[38;5;177mCheck-itoff api START...[0m")
// 	router := httprouter.New()
//
// 	commonHandlers := alice.New(loggingHandler)
// 	var _ = commonHandlers
//
// 	var routes = []R{
// 		R{"Home", "GET", "/", Index},
// 		R{"List Controller", "POST", "/lists/", Index},
// 		R{"List Controller", "GET", "/lists/:id", controllers.ListControllerCreate},
// 		R{"List Controller", "DELETE", "/lists/:id", controllers.ListControllerCreate},
// 	}
//
// 	for _, route := range routes {
// 		router.Handle(route.Method, route.Path, route.Handle)
// 	}
//
// 	chain := alice.New(auth, corsHandler).Then(router)
//
// 	log.Fatal(http.ListenAndServe(":8080", chain))
// }
//
//
//
//
//
//
//
//
// func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	res := []byte("foo")
// 	w.WriteHeader(http.StatusOK)
// 	_, err := w.Write(res)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
//
// // R is used to store routes
// type R struct {
// 	GroupName string
// 	Method    string
// 	Path      string
// 	httprouter.Handle
// }
//
// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "You are on the about page.")
// }
//
// func loggingHandler(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		t1 := time.Now()
// 		next.ServeHTTP(w, r)
// 		t2 := time.Now()
// 		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
// 	}
//
// 	return http.HandlerFunc(fn)
// }
//
// func corsHandler(h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, req *http.Request) {
// 		if origin := req.Header.Get("Origin"); origin == "" {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 			w.Header().Set("Access-Control-Allow-Headers",
// 				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 		}
// 		// Stop here if its Preflighted OPTIONS request
// 		if req.Method == "OPTIONS" {
// 			return
// 		}
// 		h.ServeHTTP(w, req)
// 	}
// 	return http.HandlerFunc(fn)
// }
//
// func auth(h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, req *http.Request) {
// 		// if true {
// 		// 	w.WriteHeader(http.StatusUnauthorized)
// 		// 	return
// 		// }
// 		h.ServeHTTP(w, req)
// 	}
//
// 	return http.HandlerFunc(fn)
// }
//