package main

import (
	"fmt"
	birdsocket "github.com/czerwonk/bird_socket"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

func birdQueryWithTimeout(socketPath string, query string, timeout time.Duration) ([]byte, error) {
	resultChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		reply, err := birdsocket.Query(socketPath, query)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- reply
	}()

	select {
	case reply := <-resultChan:
		return reply, nil
	case err := <-errChan:
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("birdsocket query timed out")
	}
}

func birdHandler(prefix *string) []string {
	socket := "/var/run/bird/bird.ctl"

	birdQuery := fmt.Sprintf("show route for %v primary all", *prefix)
	log.Debugf("Making BIRD query: %s", birdQuery)

	timeout := 15 * time.Second
	replyBinary, err := birdQueryWithTimeout(socket, birdQuery, timeout)
	if err != nil {
		log.Error(err)
		return []string{defaultASN}
	}
	reply := string(replyBinary)
	log.Debugf("%s", reply)

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
