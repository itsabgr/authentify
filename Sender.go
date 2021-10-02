package authentify

import (
	"context"
	"fmt"
	"io"
	"time"
)

type Sender interface {
	Send(ctx context.Context, to string, code Code, expireAt time.Time) error
	Name() string
	io.Closer
}

type SendersMap = map[string]Sender

func SendersToMap(senders ...Sender) (SendersMap, error) {
	sendersMap := make(SendersMap)
	for _, sender := range senders {
		_, ok := sendersMap[sender.Name()]
		if ok {
			panic(fmt.Errorf("duplicate sender %s", sender.Name()))
		}
		sendersMap[sender.Name()] = sender
	}
	return sendersMap, nil
}
