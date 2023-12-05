package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aopoltorzhicky/go_kraken/futureswebsocket"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	kraken := futureswebsocket.NewKraken(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET"))
	go kraken.ConnectAndSubscribeToOrderBooks([]string{"PF_XBTUSD"})

	for {
		select {
		case <-signals:
			log.Print("Stopping...")
			return
		case update := <-kraken.Listen():
			switch update.Feed {
			case futureswebsocket.BOOK_SNAPSHOT:
				log.Printf("Book snapshot received for %s", update.ProductId)

			case futureswebsocket.BOOK:
				log.Printf("Book update received for %s", update.ProductId)

			}
		}
	}
}
