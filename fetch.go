package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// fetch collects all debug information every tick, sums the values, and
// sends it as a single msr down the returned chan.
//
// fetch closes the returned chan when notified to do so using ctx.
func fetch(ctx context.Context, tick time.Duration, targets []string) chan msr {
	t := time.NewTicker(tick)
	msrs := make(chan msr)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(msrs)
				return
			case <-t.C:
				ms := make(chan msrErr)
				for _, target := range targets {
					go func(t string) {
						cx, cancel := context.WithTimeout(ctx, 2*time.Second)
						defer cancel()
						m, err := getMetrics(cx, t)
						ms <- msrErr{m, err}
					}(target)
				}

				errs := false
				r := msr{}
				for range targets {
					m := <-ms
					if m.e != nil {
						log.Printf("%v", m.e)
						errs = true
					}
					r.Add(m.m)
				}
				if !errs {
					msrs <- r
				}
			}
		}
	}()
	return msrs
}

// getMetrics performs http operations against a single target, parses out an
// msr, then returns it and any errors encountered.
func getMetrics(ctx context.Context, target string) (msr, error) {
	r := msr{}

	u := fmt.Sprintf("http://%s:12345/debug/metrics", target)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return r, errors.Wrap(err, "making http request object")
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, errors.Wrap(err, "json decode")
	}
	return r, nil
}
