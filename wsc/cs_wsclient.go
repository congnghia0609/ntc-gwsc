package wsc

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// type UWSClient struct {
// 	interrupt chan os.Signal
// 	done      chan struct{}
// 	url       url.URL
// 	conn      *websocket.Conn
// }

// // var addr = "localhost:8080"
// var addr = "localhost:15601"
// var instance *UWSClient
// var once sync.Once

// // GetInstanceUWS Singleton UWSClient
// func GetInstanceUWS() *UWSClient {
// 	fmt.Println("=================== UWSClient.GetInstanceUWS ===================")
// 	once.Do(func() {
// 		var uws UWSClient
// 		var err error
// 		uws, err = InitUWS()
// 		if err != nil {
// 			fmt.Println("[xxxxxx] ERROR:", err)
// 			//log.Fatal("dial:", err)
// 		}
// 		instance = &uws
// 	})
// 	return instance
// }

// // NewInstanceUWS new object UWSClient
// func NewInstanceUWS() *UWSClient {
// 	fmt.Println("=================== UWSClient.NewInstanceUWS ===================")
// 	var uws UWSClient
// 	var err error
// 	uws, err = InitUWS()
// 	if err != nil {
// 		fmt.Println("[xxxxxx] ERROR:", err)
// 		//log.Fatal("dial:", err)
// 	}
// 	instance = &uws
// 	return instance
// }

// // InitUWS UWSClient
// func InitUWS() (UWSClient, error) {
// 	var uws UWSClient
// 	var err error
// 	uws.interrupt = make(chan os.Signal, 1)
// 	signal.Notify(uws.interrupt, os.Interrupt)
// 	uws.done = make(chan struct{})

// 	uws.url = url.URL{Scheme: "ws", Host: addr, Path: "/ws/v1/cs/ETH_BTC@1h"}
// 	log.Printf("connecting to %s", uws.url.String())

// 	uws.conn, _, err = websocket.DefaultDialer.Dial(uws.url.String(), nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}

// 	return uws, err
// }

// // InitUWS2 UWSClient
// func InitUWS2(host string, path string) (UWSClient, error) {
// 	var uws UWSClient
// 	var err error
// 	uws.interrupt = make(chan os.Signal, 1)
// 	signal.Notify(uws.interrupt, os.Interrupt)
// 	uws.done = make(chan struct{})

// 	uws.url = url.URL{Scheme: "ws", Host: host, Path: path}
// 	log.Printf("connecting to %s", uws.url.String())

// 	uws.conn, _, err = websocket.DefaultDialer.Dial(uws.url.String(), nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}

// 	return uws, err
// }

// // NewInstanceUWS2 new object UWSClient
// func NewInstanceUWS2(host string, path string) *UWSClient {
// 	fmt.Println("=================== UWSClient.NewInstanceUWS2 ===================")
// 	var uws UWSClient
// 	var err error
// 	uws, err = InitUWS2(host, path)
// 	if err != nil {
// 		fmt.Println("[xxxxxx] ERROR:", err)
// 		//log.Fatal("dial:", err)
// 	}
// 	instance = &uws
// 	return instance
// }

// func (uws UWSClient) Close() {
// 	if uws.conn != nil {
// 		uws.conn.Close()
// 	}
// }

// // Reconnect is method auto-reconnect to websocket server.
// func Reconnect(host string, path string) *UWSClient {
// 	// create ticket time for every 3 seconds
// 	ticker := time.NewTicker(time.Duration(3) * time.Second)
// 	var count = 1
// 	for _ = range ticker.C {
// 		var uws *UWSClient
// 		uws = NewInstanceUWS2(host, path)

// 		// Create new websocket
// 		fmt.Printf("\nRetry Connect : %d times\n", count)
// 		count = count + 1
// 		if uws != nil {
// 			instance = uws
// 			return instance
// 		}
// 	}
// 	return nil
// }

func StartCSWSClient() {
	var uws *UWSClient
	// var err error
	uws, _ = NewInstanceWSC(NameCSWSC, "ws", "localhost:15601", "/ws/v1/cs/ETH_BTC@1h")
	defer uws.Close()

	// Thread receive message.
	go func() {
		defer uws.Close()
		defer close(uws.done)
		for {
			_, message, err := uws.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// Thread send message.
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			//err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))

			msec := t.UnixNano() / 1000000

			///// 1. Candlesticks Data.
			data := `{"tt":"1h","s":"ETH_BTC","t":` + fmt.Sprint(msec) + `,"e":"kline","k":{"c":"0.00028022","t":1533715200000,"v":"905062.00000000","h":"0.00028252","l":"0.00027787","o":"0.00027919"}}`

			err := uws.conn.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-uws.interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := uws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-uws.done:
			case <-time.After(time.Second):
			}
			uws.Close()
			return
		}
	}
}