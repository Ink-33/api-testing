/*
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package runner

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type prometheusReporter struct {
	remote string
	sync   bool
}

// NewPrometheusWriter creates a new PrometheusWriter
func NewPrometheusWriter(remote string, sync bool) TestReporter {
	return &prometheusReporter{
		remote: remote,
		sync:   sync,
	}
}

// PutRecord puts the test result into the Prometheum Pushgateway
func (w *prometheusReporter) PutRecord(record *ReportRecord) {
	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		defer wait.Done()

		name := fmt.Sprintf("response_time_%s_%s", record.Group, record.Name)

		var responseTime prometheus.Gauge
		responseTime = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "api_testing",
			Name:      name,
			Help:      "The response time in milliseconds of the API.",
		})
		responseTime.Set(float64(record.EndTime.Sub(record.BeginTime).Milliseconds()))

		if err := push.New(w.remote, "api-testing").
			Collector(responseTime).
			Push(); err != nil {
			fmt.Println("Could not push completion time to Pushgateway:", err)
		}
	}()
	if w.sync {
		wait.Wait()
	}
	return
}

func (w *prometheusReporter) GetAllRecords() []*ReportRecord {
	// no support
	return nil
}

func (r *prometheusReporter) ExportAllReportResults() (result ReportResultSlice, err error) {
	// no support
	return
}