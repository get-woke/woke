/**
 * Copyright 2022 Cisco and its affiliates
 * All rights reserved.
**/

package printer

import (
	"fmt"
	"io"

	config "github.com/get-woke/woke/pkg/config"
	"github.com/get-woke/woke/pkg/result"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// Pushgateway is a way to push metrics within prometheus
type Pushgateway struct {
	writer         io.Writer
	labels         map[string]string
	pushgatewayURL string
}

// NewPushgateway returns a Pushgateway Printer with color optionally disabled
func NewPushgateway(w io.Writer, c *config.Config) *Pushgateway {
	return &Pushgateway{
		writer:         w,
		labels:         c.Labels,
		pushgatewayURL: c.PushgatewayURL,
	}
}

func (t *Pushgateway) PrintSuccessExitMessage() bool {
	return true
}

// Print prints the file results
func (t *Pushgateway) Print(fs *result.FileResults) error {

	for _, r := range fs.Results {
		pos := fmt.Sprintf("%d:%d-%d",
			r.GetStartPosition().Line,
			r.GetStartPosition().Column,
			r.GetEndPosition().Column)

		labelString := ""
		woke_result := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "woke_result",
			Help: "Inclusive woke result",
		})
		pusher := push.New(t.pushgatewayURL, "woke")
		if len(t.labels) > 0 {
			first := true
			for k, v := range t.labels {
				if first {
					labelString += ", "
					first = false
				}
				labelString += k + "=" + v + ", "
				pusher = pusher.Grouping(k, v)
			}
			labelString = labelString[:len(labelString)-2]
		}

		woke_result.Set(1)
		pusher = pusher.Grouping("term", r.GetRuleName()).
			Grouping("instance", fs.Filename+":"+pos).
			Collector(woke_result)

		if err := pusher.Push(); err != nil {
			fmt.Println("Could not push woke_result to Pushgateway:", err)
		} else {
			fmt.Printf("push woke_result to Pushgateway: %v", pusher)
		}
	}

	return nil
}

func (t *Pushgateway) Start() {
}

func (t *Pushgateway) End() {
}
