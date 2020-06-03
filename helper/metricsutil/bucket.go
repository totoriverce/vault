package metricsutil

import (
	"sort"
	"time"
)

var bucketBoundaries = []struct {
	Value time.Duration
	Label string
}{
	{1 * time.Minute, "1m"},
	{10 * time.Minute, "10m"},
	{20 * time.Minute, "20m"},
	{1 * time.Hour, "1h"},
	{2 * time.Hour, "2h"},
	{24 * time.Hour, "1d"},
	{2 * 24 * time.Hour, "2d"},
	{7 * 24 * time.Hour, "7d"},
	{30 * 24 * time.Hour, "30d"},
}

const overflowBucket = "+Inf"

// TTLBucket computes the label to apply for a token TTL.
func TTLBucket(ttl time.Duration) string {
	upperBound := sort.Search(
		len(bucketBoundaries),
		func(i int) bool {
			return ttl <= bucketBoundaries[i].Value
		},
	)
	if upperBound >= len(bucketBoundaries) {
		return overflowBucket
	} else {
		return bucketBoundaries[upperBound].Label
	}

}
