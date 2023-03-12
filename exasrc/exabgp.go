package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func MessageHandler() {

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

	respondCommand(event)
}

func respondCommand(event Event) {
	log.Debugf("%+v", event)
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
	}

}
