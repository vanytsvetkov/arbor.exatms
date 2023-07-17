package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

type HoldTimer struct {
	Current time.Duration
	Max     time.Duration
}

var (
	HoldTime = HoldTimer{
		Current: time.Second * 180,
		Max:     time.Second * 180,
	}
	RIB = make(map[string]bool)
)

func heartBeat() {
	for range time.Tick(time.Second * 1) {
		if HoldTime.Current > 0 {
			HoldTime.Current -= time.Second * 1
			log.Debugf("Hold timer eq. %v", HoldTime.Current)
		} else if HoldTime.Current <= 0 && len(RIB) > 0 {
			log.Info("Hold Time Exceeded. Withdrawing all!")
			withdrawAll()
			HoldTime.Current = 0
		}
	}
}

func withdrawAll() {
	log.Debugf("my RIB is %v", RIB)
	for pfx := range RIB {
		withdrawCmd := fmt.Sprintf("neighbor %s withdraw route %s",
			neighbor,
			pfx,
		)
		fmt.Println(withdrawCmd)
		log.Info(withdrawCmd)
	}

	RIB = make(map[string]bool)
}

func MessageHandler() {

	go heartBeat()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if len(msg) > 0 {
			log.Debug(msg)
			if msg != "done" && msg != "error" {
				processMessage(&msg)
			}
		}
	}

	if scanner.Err() != nil {
		log.Error(scanner.Err())
	}

}

func processMessage(msg *string) {

	var event Event
	err := json.Unmarshal([]byte(*msg), &event)
	if err != nil {
		log.Error(err)
		return
	}

	if event.Type == "keepalive" {
		processKeepAlive()
	} else if event.Type == "update" {
		processUpdate(event)
	}
}

func processKeepAlive() {
	HoldTime.Current = HoldTime.Max
	log.Debugf("Got KEEPALIVE, set Hold Timer eq. %v", HoldTime.Max)
}

func processUpdate(event Event) {

	HoldTime.Current = HoldTime.Max
	log.Debugf("Got UPDATE, set Hold Timer eq. %v", HoldTime.Max)

	// Make announce
	for nexthop, announces := range event.Neighbor.Message.Update.Announce.IPv4Unicast {
		for _, announce := range announces {

			asPath := birdHandler(&announce.NLRI)

			var communities []string
			for _, cmt := range event.Neighbor.Message.Update.Attribute.Community {
				communities = append(communities, fmt.Sprintf("%d:%d", cmt[0], cmt[1]))
			}

			announceCmd := fmt.Sprintf("neighbor %s announce route %s next-hop %s community %v local-preference %d as-path %v",
				neighbor,
				announce.NLRI,
				nexthop,
				communities,
				event.Neighbor.Message.Update.Attribute.LocalPreference,
				asPath,
			)
			fmt.Println(announceCmd)
			log.Info(announceCmd)
			RIB[announce.NLRI] = true
			log.Debugf("RIB is %v", RIB)
		}
	}

	// Make withdraw
	for _, withdraw := range event.Neighbor.Message.Update.Withdraw.IPv4Unicast {
		withdrawCmd := fmt.Sprintf("neighbor %s withdraw route %s",
			neighbor,
			withdraw.NLRI,
		)
		fmt.Println(withdrawCmd)
		log.Info(withdrawCmd)
		delete(RIB, withdraw.NLRI)
		log.Debugf("RIB is %v", RIB)
	}

}
