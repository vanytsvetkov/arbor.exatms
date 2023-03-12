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

	routeB, err := birdsocket.Query(socket, fmt.Sprintf("show route for %v primary all", *prefix))
	if err != nil {
		log.Error(err)
		return []string{defaultASN}
	}
	route := string(routeB)
	log.Debugf("%s", route)

	// common case: ^\s+BGP.(\w+):\s+(.+)\s*$
	re := regexp.MustCompile(`BGP\.as_path:\s([\s\d+]*)`)

	match := re.FindStringSubmatch(route)
	if len(match) > 1 {
		group := strings.TrimSpace(match[1])
		//log.Debugf("as_path: %s", group)
		return strings.Split(group, " ")
	} else {
		return []string{defaultASN}
	}

}
