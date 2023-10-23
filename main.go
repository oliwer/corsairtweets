// Twitter IRC Bot for #Corsair
// Written by Olivier Duclos (odyssey)

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/plugins/golang"
	"github.com/StalkR/goircbot/plugins/imdb"
	"github.com/StalkR/goircbot/plugins/invite"
	"github.com/StalkR/goircbot/plugins/sed"
	"github.com/StalkR/goircbot/plugins/up"
	"github.com/StalkR/goircbot/plugins/urban"
	"github.com/StalkR/goircbot/plugins/weather"
	"github.com/oliwer/corsairtweets/lastseen"
	"github.com/oliwer/corsairtweets/timein"
)

var (
	host     = flag.String("host", "irc.oftc.net", "Server host[:port]")
	ssl      = flag.Bool("ssl", true, "Connect with SSL")
	nick     = flag.String("nick", "twittard", "Bot nick")
	ident    = flag.String("ident", "corsairtwitterebot", "Bot ident")
	channels = flag.String("channels", "#Corsair", "Channels to join, comma-separated")
	ignored  = strings.Split(os.Getenv("IGNORE_NICKS"), ",")
	owmkey   = os.Getenv("OPENWEATHERMAP_KEY")
)

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done() // it not allowed to call wg.Wait() before it has been used.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-c
		log.Printf("Received signal %s. Exiting...\n", sig.String())
		wg.Wait()
		os.Exit(1)
	}()

	b := bot.NewBot(*host, *ssl, *nick, *ident, strings.Split(*channels, ","))
	golang.Register(b)
	imdb.Register(b)
	invite.Register(b)
	lastseen.Register(b, ignored, &wg)
	sed.Register(b)
	timein.Register(b)
	up.Register(b)
	urban.Register(b)
	weather.Register(b, owmkey)
	b.Run()
}
