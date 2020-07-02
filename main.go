package main

import (
	"fmt"
	"os"

	"github.com/ddo/go-fast"
	"github.com/honeycombio/libhoney-go"
)

func main() {
	libhoney.Init(libhoney.Config{
		WriteKey: os.Getenv("HONEYCOMB_API_KEY"),
		Dataset:  "Mammoth Lakes Internet",
	})

	fastCom := fast.New()

	// init
	err := fastCom.Init()
	if err != nil {
		panic(err)
	}

	// get urls
	urls, err := fastCom.GetUrls()
	if err != nil {
		panic(err)
	}

	// measure
	KbpsChan := make(chan float64)

	go func() {
		for Kbps := range KbpsChan {
			fmt.Println("Creating event")
			event := libhoney.NewEvent()
			event.Add(map[string]interface{}{
				"speed_kbps": Kbps,
			})
			event.Send()
			fmt.Printf("%.2f Kbps %.2f Mbps\n", Kbps, Kbps/1000)
		}

		fmt.Println("done")
		libhoney.Close()
	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		panic(err)
	}
}
