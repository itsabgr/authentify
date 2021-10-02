package authentify

import "time"

type Sender interface {
	Send(to string, code Code, expireAt time.Time) error
	Name() string
}
