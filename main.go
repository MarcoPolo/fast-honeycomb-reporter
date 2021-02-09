package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/ddo/go-fast"
	"github.com/go-ping/ping"
	"github.com/honeycombio/libhoney-go"
)

const waitBetweenDownloads = 30 * time.Minute
const waitVariance = 1 * time.Minute

func main() {
	err := libhoney.Init(libhoney.Config{
		WriteKey: os.Getenv("HONEYCOMB_API_KEY"),
		Dataset:  "Mammoth Lakes Internet",
	})
	fmt.Println("Err", err)
	defer libhoney.Close()

	go func() {
		for {
			pingTest()
		}
	}()

	for {
		downloadSpeedTest()
		secs := time.Duration(waitVariance.Seconds()*rand.Float64()) * time.Second
		time.Sleep(waitBetweenDownloads + secs)
	}
}

func downloadSpeedTest() {
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

	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		panic(err)
	}

}

func pingTest() {
	pinger, err := ping.NewPinger("1.1.1.1")
	if err != nil {
		panic(err)
	}

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		event := libhoney.NewEvent()
		event.Add(map[string]interface{}{
			"latency_ms": pkt.Rtt.Milliseconds(),
		})
		err := event.Send()

		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v err=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, err)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.Run()
	if err != nil {
		panic(err)
	}

}
