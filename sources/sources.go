package sources

import "log-aggregator/types"

// Source is anything that can produce logs, e.g. Journald
type Source interface {
	Start(chan<- *types.Record)
	Stop()
}
