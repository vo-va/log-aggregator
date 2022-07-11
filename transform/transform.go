package transform

import "log-aggregator/types"

type Transformer func(rec *types.Record) (*types.Record, error)
