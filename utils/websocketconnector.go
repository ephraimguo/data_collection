package utils

import (
	"data_collection/controllers/binance"
	"data_collection/controllers/deribit"
	"data_collection/database"
	"database/sql"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func Subscribe(platform string, db *sql.DB) {
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1) // os.Signal type channel
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{
		Scheme: "wss",
		Host:   "www.deribit.com",
		Path:   "/ws/api/v2",
		//ForceQuery: true,
		//RawQuery: "",
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

			switch platform {
			case "binance":
				r := binance.Parse(message)
				database.InsertSingleRecord(r.Quote, db)
			case "deribit":
				r := deribit.Parse(message)
				fmt.Printf("r: %+v\n", r)
				if r != nil {
					database.InsertSingleRecord(r.Params.Quote, db)
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			var err error
			if platform == "deribit" {
				payload := "{\"method\":\"public/subscribe\",\"params\":{\"channels\":[\"quote.BTC-PERPETUAL\",\"quote.ETH-PERPETUAL\"]},\"jsonrpc\":\"2.0\",\"id\":7}"
				err = conn.WriteMessage(websocket.TextMessage, []byte(payload))
			} else {
				err = conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			}
			if err != nil {
				log.Println("write text :", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure,
					"",
				),
			)
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
