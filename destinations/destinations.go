package destinations

import "log-aggregator/types"

// Destination is something we can write logs to
type Destination interface {
	Start(<-chan *types.Record, chan<- types.Cursor)
}
