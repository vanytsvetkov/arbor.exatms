package main

import (
	"fmt"
	birdsocket "github.com/czerwonk/bird_socket"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func birdHandler(prefix *string) []string {
	socket := "/var/run/bird/bird.ctl"

	birdQuery := fmt.Sprintf("show route for %v primary all", *prefix)
	log.Debugf("Making BIRD query: %s", birdQuery)

	replyBinary, err := birdsocket.Query(socket, birdQuery)
	if err != nil {
		log.Error(err)
		return []string{defaultASN}
	}
	reply := string(replyBinary)
	log.Debugf("%s", reply)

	if reply == "Network not found" {
		return []string{defaultASN}
	} else {
		// common case: ^\s+BGP.(\w+):\s+(.+)\s*$
		re := regexp.MustCompile(`BGP\.as_path:\s([\s\d+]*)`)

		match := re.FindStringSubmatch(reply)
		if len(match) > 1 {
			group := strings.TrimSpace(match[1])
			//log.Debugf("as_path: %s", group)
			return strings.Split(group, " ")
		} else {
			return []string{defaultASN}
		}
	}

}
