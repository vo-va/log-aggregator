package eleven

import (
	"os"

	"log-aggregator/transform"
	"log-aggregator/types"
)

const (
	EnvProduct    = "ELEVEN_PRODUCT"
	EnvComponent     = "ELEVEN_COMPONENT"
)

func New() transform.Transformer {

	product := os.Getenv(EnvProduct)
	component := os.Getenv(EnvComponent)

	return func(rec *types.Record) (*types.Record, error) {
		rec.Fields["product"] = product
		rec.Fields["component"] = component

		formattedTime := rec.Time.Format("2006-01-02T15:04:05.000")
		rec.Fields["when"] = formattedTime
		return rec, nil
	}
}
