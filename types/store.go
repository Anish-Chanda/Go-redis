package types

import "time"

type Store struct {
	Value string
	Exp   time.Time
}
