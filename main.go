package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

var tick = flag.Duration("t", 5*time.Second, "time between polls")

func main() {
	flag.Parse()
	targets := flag.Args()

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		cancel()
	}()

	msrs := fetch(ctx, *tick, targets)

	var n, o float64
	last := <-msrs
	for m := range msrs {
		d := m.Delta(last)
		last = m
		o += float64(d.V1)
		n += float64(d.Authed)
		log.Printf("(%5.2f) %+v", n/(n+o), d)
	}
}
