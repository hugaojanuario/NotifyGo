package kafka

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func TestBrokers(brokers string) error {
	for _, broker := range strings.Split(brokers, ",") {
		broker = strings.TrimSpace(broker)
		conn, err := net.DialTimeout("tcp", broker, 5*time.Second)
		if err != nil {
			return fmt.Errorf("could not reach broker %s: %w", broker, err)
		}
		conn.Close()
	}
	return nil
}
