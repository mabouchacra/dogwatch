package main

import (
	"fmt"

	"github.com/mabouchacra/dogwatch/api"
)

func main() {
	metrics := api.GetLogBasedMetricVolumeWithLimit(500)
	fmt.Println(metrics)
}
