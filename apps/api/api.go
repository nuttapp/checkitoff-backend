package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/nuttapp/checkitoff-backend/apps/api/config"
	"github.com/nuttapp/checkitoff-backend/dal"
)

const (
	ProducerConnectionError = "Failed to connect to NSQ producer"
	ProducerPublishError    = "Failed to publish to NSQ producer"
)

type APIServer struct {
	Consumer *nsq.Consumer
	Ctx      *APIContext
}

func (s *APIServer) Start() {
	fmt.Printf("Starting consumer for topic: \"%s\", channel: \"%s\"...",
		s.Ctx.Cfg.NSQ.SubTopic, s.Ctx.Cfg.NSQ.SubChannel)
	// start listening for messages from nsq
	// start listening on endpoints
	consumer, err := nsq.NewConsumer(s.Ctx.Cfg.NSQ.SubTopic, s.Ctx.Cfg.NSQ.SubChannel, s.Ctx.NSQCfg)
	if err != nil {
		panic(err)
	}

	// var logBuf bytes.Buffer
	// logger := log.New(&logBuf, "", log.LstdFlags)
	// consumer.SetLogger(logger, nsq.LogLevelDebug)

	consumer.AddHandler(s)
	err = consumer.ConnectToNSQLookupd(s.Ctx.Cfg.NSQ.LookupdHTTPAddr)
	if err != nil {
		panic(err)
	}

	s.Consumer = consumer
}

func (s *APIServer) Stop() {
	// waiting for consumer to stop take about 10ms
	fmt.Println("Stopping producer...")
	s.Ctx.Producer.Stop()
	fmt.Println("Stopping consumer...")
	s.Consumer.Stop()
	<-s.Consumer.StopChan

	// disconnect producer
	// disconnect consumer
}

func (s *APIServer) CreateTopic(topic string) {
	if len(topic) == 0 {
		panic("topic cannot be empty")
	}
	if len(s.Ctx.Cfg.NSQ.LookupdHTTPAddr) == 0 {
		panic("LookupdHttpAddr cannot be empty")
	}
	url := fmt.Sprintf("http://%s/create_topic?topic=%s", s.Ctx.Cfg.NSQ.LookupdHTTPAddr, topic)
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Couldn't create topic \"%s\". Status Code %d", topic, res.StatusCode)
		panic(errMsg)
	}
	fmt.Printf("Createed topic %s. Status \"%s\"\n", topic, res.Status)
}

func (s *APIServer) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("received msg: %s\n", msg.ID)
	listMsg, err := dal.DeserializeListEvent(msg.Body)
	if err != nil {
		return err
	}
	s.Ctx.guard.Lock()
	replyChan, ok := s.Ctx.Messages[listMsg.ID]
	s.Ctx.guard.Unlock()

	if ok {
		replyChan <- msg
	} else {
		fmt.Println("could not find message")
	}

	// receive message from nsq
	// lookup id and get channel
	// if id exists send the received msg on channel
	// if id not exist log error
	// mh.testChan <- msg
	return nil
}

// validate token compare hash of nic id & token to what's stored in db
// save to db
type Req struct {
	ID   string                 // ID assigned to an individual request
	Msgs map[string]interface{} // Messages generated as part of the request
}

// Done sends the array of bytes to the client that initiated the request
func (r *Req) Done(data []byte) {
}

type APIContext struct {
	Requests map[string]Req
	Messages map[string]chan *nsq.Message
	guard    sync.RWMutex

	NSQCfg   *nsq.Config
	Cfg      *config.Config
	Producer *nsq.Producer
}

func (ctx *APIContext) Publish(reqID string, event []byte) error {
	err := ctx.Producer.Publish(ctx.Cfg.NSQ.PubTopic, event)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *APIContext) PublishWithReplyChan(eID string, event []byte) (chan *nsq.Message, error) {
	reply := make(chan *nsq.Message)
	ctx.guard.Lock()
	ctx.Messages[eID] = reply
	ctx.guard.Unlock()

	err := ctx.Producer.Publish(ctx.Cfg.NSQ.PubTopic, event)
	if err != nil {
		return nil, err
	}
	fmt.Println("Message published successfully...")
	return reply, err
}

// NewAPIServer(APIContext)
// Server.start

// InitContext Initializes an new API Context and
// returns a pointer to it.  Context is meant to span across
// API requests, it can be safely accessed accross goroutines
func NewContext(environment string) *APIContext {
	apiCfg := config.NewConfig()

	nsqCfg := nsq.NewConfig()
	nsqCfg.MaxInFlight = 10

	producer, err := nsq.NewProducer(apiCfg.NSQ.ProducerTCPAddr, nsqCfg)
	if err != nil {
		panic(fmt.Sprintf("%s: %s", ProducerConnectionError, err))
	}

	context := &APIContext{
		Messages: make(map[string]chan *nsq.Message),
		Producer: producer,
		Cfg:      apiCfg,
		NSQCfg:   nsqCfg,
	}

	return context
}

func main() {
	// Context = NewContext()

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	//
	// nsqCfg := nsq.NewConfig()
	// nsqCfg.MaxInFlight = 5
	//
	// conCfg := &common.ConsumerConfig{
	// 	Topic:           "api01",
	// 	Channel:         "client_events",
	// 	LookupdHTTPaddr: "127.0.0.1:4161",
	// 	Concurrency:     10,
	// }
	//
	// handler := &APIMessageHandler{}
	// mc := common.NewMessageConsumer(conCfg, nsqCfg, handler)
	// mc.Start()
	//
	// for {
	// 	select {
	// 	// On quit signal stop the api consumer
	// 	case <-signalChan:
	// 		fmt.Println("")
	// 		mc.Stop()
	// 	case err := <-mc.StopChan:
	// 		if err != nil {
	// 			fmt.Printf("MessageConsumer stopped with error: %s", err)
	// 		}
	// 		fmt.Println("END")
	// 		return
	// 	}
	// }
}

// type APIMessageHandler struct{}
//
// func (mh *APIMessageHandler) HandleMessage(m *nsq.Message) error {
// 	// Check if message is related to a connected client
// 	// send response to client
// 	fmt.Println(string(m.Body))
// 	return nil
// }
//
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
