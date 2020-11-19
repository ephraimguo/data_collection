package utils

import (
	"data_collection/config"
	"data_collection/controllers/binance"
	"data_collection/controllers/deribit"
	"data_collection/database"
	"database/sql"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func Subscribe(platform string, db *sql.DB) {
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1) // os.Signal type channel
	signal.Notify(interrupt, os.Interrupt)

	var u url.URL
	if platform == "binance" {
		forceQuery, _ := strconv.ParseBool(config.Binance["ForceQuery"])
		u = url.URL{
			Scheme:     config.Binance["Scheme"],
			Host:       config.Binance["Host"],
			Path:       config.Binance["Path"],
			ForceQuery: forceQuery,
			RawQuery:   config.Binance["RawQuery"],
		}

	} else if platform == "deribit" {
		u = url.URL{
			Scheme: config.Deribit["Scheme"],
			Host:   config.Deribit["Host"],
			Path:   config.Deribit["Path"],
		}
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
				log.Println("read message error :", err)
			}

			switch platform {
			case "binance":
				r := binance.Parse(message)
				if r != nil {
					r.Quote.Platform = "binance"
					database.InsertSingleRecord(r.Quote, db)
				}
			case "deribit":
				r := deribit.Parse(message)
				if r != nil {
					r.Params.Quote.Platform = "deribit"
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
				err = conn.WriteMessage(websocket.TextMessage, []byte(config.Deribit["Payload"])) // heart beat
			} else if platform == "binance" {
				err = conn.WriteMessage(websocket.TextMessage, []byte(t.String())) // heart beat
			}

			if err != nil {
				log.Println("write text error:", err)
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
