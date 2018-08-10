package wsc

import (
	"fmt"
	"log"
	"ntc-gwsc/util"
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
// var addr = "localhost:15801"
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

// 	uws.url = url.URL{Scheme: "ws", Host: addr, Path: "/ws/v1/tk"}
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

func (wsc *UWSClient) recvTK() {
	util.TCF{
		Try: func() {
			defer wsc.Close()
			defer close(wsc.done)
			for {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvTK Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *UWSClient) sendTK() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case t := <-ticker.C:
					//err := uws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))

					msec := t.UnixNano() / 1000000

					///// 1. Ticker24h Data.
					data := `{"t":` + fmt.Sprint(msec) + `,"d":[{"q":"704.79245092","s":"GTO_BTC","c":"0.00001528","pc":"0.92470277","t":1533873600000,"ch":"0.00000014","h":"0.00001619","l":"0.00001510","n":"0.00001528","o":"0.00001518"},{"q":"0.00000000","s":"GTO_KNOW","c":"3.24640000","pc":"0.00000000","t":1533873600000,"ch":"0.00000000","h":"3.24640000","l":"3.24640000","n":"3.24640000","o":"3.24640000"},{"q":"8639796.02353590","s":"ETH_USDT","c":"358.57700000","pc":"0.26513576","t":1533873600000,"ch":"0.95000000","h":"369.52100000","l":"354.71300000","n":"359.25700000","o":"361.38300000"},{"q":"219.19671119","s":"MEDX_BTC","c":"0.00000113","pc":"8.65384615","t":1533873600000,"ch":"0.00000009","h":"0.00000117","l":"0.00000108","n":"0.00000113","o":"0.00000115"},{"q":"0.00000000","s":"MEDX_KNOW","c":"0.09000000","pc":"0.00000000","t":1533873600000,"ch":"0.00000000","h":"0.09000000","l":"0.09000000","n":"0.09000000","o":"0.09000000"},{"q":"0.44078389","s":"PIB_BTC","c":"0.00000024","pc":"-4.00000000","t":1533873600000,"ch":"-0.00000001","h":"0.00000025","l":"0.00000024","n":"0.00000024","o":"0.00000025"},{"q":"5962.39505723","s":"GTO_ETH","c":"0.00027489","pc":"3.20630749","t":1533873600000,"ch":"0.00000854","h":"0.00028852","l":"0.00026535","n":"0.00027489","o":"0.00026594"},{"q":"0.00000000","s":"XEM_KNOW","c":"3.69900000","pc":"0.00000000","t":1533873600000,"ch":"0.00000000","h":"3.69900000","l":"3.69900000","n":"3.69900000","o":"3.69900000"},{"q":"0.00000000","s":"NPXSXEM_KNOW","c":"0.05870000","pc":"0.00000000","t":1533873600000,"ch":"0.00000000","h":"0.05870000","l":"0.05870000","n":"0.05870000","o":"0.05870000"},{"q":"654.68724974","s":"KNOW_BTC","c":"0.00000719","pc":"-5.39473684","t":1533873600000,"ch":"-0.00000041","h":"0.00000764","l":"0.00000718","n":"0.00000719","o":"0.00000764"},{"q":"76.16164350","s":"NPXSXEM_BTC","c":"0.00000064","pc":"68.42105263","t":1533873600000,"ch":"0.00000026","h":"0.00000064","l":"0.00000033","n":"0.00000064","o":"0.00000035"},{"q":"6171.55098716","s":"KNOW_ETH","c":"0.00013053","pc":"0.13041810","t":1533873600000,"ch":"0.00000017","h":"0.00013163","l":"0.00012901","n":"0.00013052","o":"0.00013062"},{"q":"0.00000000","s":"XPX_KNOW","c":"0.15000000","pc":"0.00000000","t":1533873600000,"ch":"0.00000000","h":"0.15000000","l":"0.15000000","n":"0.15000000","o":"0.15000000"},{"q":"1313.21595908","s":"NPXSXEM_ETH","c":"0.00001148","pc":"67.98825256","t":1533873600000,"ch":"0.00000463","h":"0.00001150","l":"0.00000638","n":"0.00001144","o":"0.00000680"},{"q":"4623.10399762","s":"TRX_ETH","c":"0.00006955","pc":"1.47359206","t":1533873600000,"ch":"0.00000101","h":"0.00007713","l":"0.00006730","n":"0.00006955","o":"0.00006857"},{"q":"887.50735317","s":"ETH_BTC","c":"0.05587000","pc":"-1.94805194","t":1533873600000,"ch":"-0.00111000","h":"0.05717000","l":"0.05528000","n":"0.05587000","o":"0.05708000"},{"q":"10758251.09466018","s":"BTC_USDT","c":"6423.80600000","pc":"2.10543181","t":1533873600000,"ch":"132.46000000","h":"6617.24600000","l":"6272.58000000","n":"6423.80600000","o":"6339.67800000"},{"q":"4871064.26230000","s":"KNOW_USDT","c":"0.04750000","pc":"1.06382978","t":1533873600000,"ch":"0.00050000","h":"0.05100000","l":"0.04610000","n":"0.04750000","o":"0.04720000"},{"q":"1457.27874423","s":"MEDX_ETH","c":"0.00001909","pc":"-0.57291666","t":1533873600000,"ch":"-0.00000011","h":"0.00001953","l":"0.00001797","n":"0.00001909","o":"0.00001920"},{"q":"902.86134111","s":"TRX_BTC","c":"0.00000388","pc":"-0.25706940","t":1533873600000,"ch":"-0.00000001","h":"0.00000415","l":"0.00000378","n":"0.00000388","o":"0.00000390"},{"q":"2565.47101698","s":"XPX_ETH","c":"0.00001553","pc":"7.39972337","t":1533873600000,"ch":"0.00000107","h":"0.00001600","l":"0.00001413","n":"0.00001553","o":"0.00001445"},{"q":"1528.68461514","s":"XEM_ETH","c":"0.00033936","pc":"0.67341066","t":1533873600000,"ch":"0.00000227","h":"0.00037155","l":"0.00032530","n":"0.00033936","o":"0.00033572"},{"q":"0.00000000","s":"TRX_KNOW","c":"0.75930000","pc":"0.00000000","t":1533873600000,"ch":"0.00000000","h":"0.75930000","l":"0.75930000","n":"0.75930000","o":"0.75930000"},{"q":"8034162.68900000","s":"GTO_USDT","c":"0.09700000","pc":"6.59340659","t":1533873600000,"ch":"0.00600000","h":"0.11800000","l":"0.08500000","n":"0.09700000","o":"0.10000000"},{"q":"798.73801019","s":"DGX_ETH","c":"0.10190200","pc":"0.80324463","t":1533873600000,"ch":"0.00081200","h":"0.10191000","l":"0.10111000","n":"0.10190200","o":"0.10131000"},{"q":"140.49534349","s":"XPX_BTC","c":"0.00000088","pc":"6.02409638","t":1533873600000,"ch":"0.00000005","h":"0.00000091","l":"0.00000081","n":"0.00000088","o":"0.00000083"},{"q":"142.37684601","s":"XEM_BTC","c":"0.00001898","pc":"-0.26274303","t":1533873600000,"ch":"-0.00000005","h":"0.00001934","l":"0.00001832","n":"0.00001898","o":"0.00001908"}],"e":"ticker24h"}`
					err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(data))
					if err != nil {
						log.Println("write:", err)
						return
					}
				case <-wsc.interrupt:
					log.Println("interrupt")
					// To cleanly close a connection, a client should send a close
					// frame and wait for the server to close the connection.
					err := wsc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					select {
					case <-wsc.done:
					case <-time.After(time.Second):
					}
					wsc.Close()
					return
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.sendTK Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func NewTKWSClient() *UWSClient {
	var tkwsc *UWSClient
	// var err error
	tkwsc, _ = NewInstanceWSC(NameTKWSC, "ws", "localhost:15801", "/ws/v1/tk")
	//defer uws.Close()
	return tkwsc
}

func (tkwsc *UWSClient) StartTKWSClient() {
	// Thread receive message.
	go tkwsc.recvTK()
	// Thread send message.
	go tkwsc.sendTK()
}
