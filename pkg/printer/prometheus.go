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
	push "github.com/prometheus/client_golang/prometheus/push"
	"github.com/rs/zerolog/log"
)

// Prometheus is a output format with prometheus metrics
type Prometheus struct {
	writer         io.Writer
	labels         map[string]string
	pushgatewayURL string
}

// NewPrometheus returns a Prometheus Printer with color optionally disabled
func NewPrometheus(w io.Writer, c *config.Prometheus) *Prometheus {
	return &Prometheus{
		writer:         w,
		labels:         c.Labels,
		pushgatewayURL: c.PushgatewayURL,
	}
}

func (t *Prometheus) PrintSuccessExitMessage() bool {
	return true
}

// Print prints the file results or send it to pushgateway if URL provided
func (t *Prometheus) Print(fs *result.FileResults) error {

	for _, r := range fs.Results {
		pos := fmt.Sprintf("%d:%d-%d",
			r.GetStartPosition().Line,
			r.GetStartPosition().Column,
			r.GetEndPosition().Column)

		labelString := ""

		if len(t.labels) > 0 {
			first := true
			for k, v := range t.labels {
				if first {
					labelString += ", "
					first = false
				}
				labelString += k + "=" + v + ", "
			}
			labelString = labelString[:len(labelString)-2]
		}
		fmt.Fprintf(t.writer, "woke_result{file=\"%s:%s\", term=\"%s\"%s} 1 \n",
			fs.Filename, pos, r.GetRuleName(), labelString)

		if t.pushgatewayURL != "" {
			pusher := push.New(t.pushgatewayURL, "woke")
			woke_result := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "woke_result",
				Help: "Inclusive woke result",
			})
			for k, v := range t.labels {
				pusher = pusher.Grouping(k, v)
			}
			pusher = pusher.Grouping("term", r.GetRuleName()).
				Grouping("instance", fs.Filename+":"+pos).
				Collector(woke_result)

			if err := pusher.Push(); err != nil {
				log.Error().Err(err).Msg("Could not push woke_result to Pushgateway")
			} else {
				log.Debug().Msg("push woke_result to Pushgateway done")
			}
		}

	}

	return nil
}

func (t *Prometheus) Start() {
}

func (t *Prometheus) End() {
}
