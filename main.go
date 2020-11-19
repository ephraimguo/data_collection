package main

import (
	controllers "data_collection/cotrollers"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

// var addr = flag.String("addr", "fstream.binance.com/stream?streams=btcusdt@depth10/ethusdt@depth10@500ms", "http service address")

func main() {
	//db := controllers.GetConnection()

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1) // os.Signal type channel
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "wss",
		Host: "fstream.binance.com",
		Path: "stream",
		ForceQuery: true,
		RawQuery: "streams=btcusdt@depth10/ethusdt@depth10@500ms",
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
			}

			//log.Printf("recv: %s", message)
			r := controllers.Parse(message)
			controllers.InsertSingleRecord(r.Quote)
			// log.Printf("r: %+v", *r)

		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- done:
			return
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}