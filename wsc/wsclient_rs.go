/**
 *
 * @author nghiatc
 * @since Aug 8, 2018
 */

package wsc

import (
	"encoding/json"
	"fmt"
	"log"
	"ntc-gwsc/conf"
	"ntc-gwsc/util"
	"time"

	"github.com/gorilla/websocket"
)

func (wsc *NWSClient) recvRS() {
	util.TCF{
		Try: func() {
			defer wsc.Close()
			defer close(wsc.done)
			for {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					wsc.Reconnect()
					// return
				}
				//log.Printf("url: %s", wsc.url.String())
				log.Printf("recvRS: %s", message)
				if len(message) > 0 {
					var data map[string]interface{}
					json.Unmarshal(message, &data)
					if data["list_symbol"] != nil {
						listsymbol := data["list_symbol"].(string)
						log.Printf("======= recvRS listsymbol: %s", listsymbol)
						if len(listsymbol) > 0 {
							//wss.ReloadMapSymbol(listsymbol)
						}
					}
				}
			}
		},
		Catch: func(e util.Exception) {
			log.Printf("wsc.recvRS Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

func (wsc *NWSClient) sendRS() {
	util.TCF{
		Try: func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case t := <-ticker.C:
					//err := nws.conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
					msec := t.UnixNano() / 1000000
					///// 1. Ticker24h Data.
					data := `{"t":` + fmt.Sprint(msec) + `,"list_symbol":"BTC_USDT;ETH_USDT;KNOW_USDT;GTO_USDT"}`
					err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(data))
					if err != nil {
						log.Println("write:", err)
						//return
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
			log.Printf("wsc.sendRS Caught %v\n", e)
		},
		Finally: func() {
			//log.Println("Finally...")
		},
	}.Do()
}

// NewRSWSClient new instance RSWSClient of NWSClient
func NewRSWSClient() *NWSClient {
	var rswsc *NWSClient
	c := conf.GetConfig()
	scheme := c.GetString(NameCRWSC + ".wsc.scheme")
	address := c.GetString(NameCRWSC + ".wsc.host")
	path := c.GetString(NameCRWSC + ".wsc.path")
	log.Printf("################ RSWSClient[%s] start...", NameRSWSC)
	rswsc, _ = NewInstanceWSC(NameRSWSC, scheme, address, path)
	return rswsc
}

// StartRSWSClient start RSWSClient
func (wsc *NWSClient) StartRSWSClient() {
	// Thread receive message.
	go wsc.recvRS()
	// Thread send message.
	// go wsc.sendRS()
}
