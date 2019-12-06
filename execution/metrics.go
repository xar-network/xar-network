/*

Copyright 2019 All in Bits, Inc
Copyright 2019 Xar Network

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package execution

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	ProcessingTime  metrics.Histogram
	OrdersProcessed metrics.Histogram
}

var exMetrics = &Metrics{
	ProcessingTime: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "dex",
		Subsystem: "execution",
		Name:      "execution_time",
		Help:      "Time for all match, and fill operations to complete.",
		Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
	}, []string{}),
	OrdersProcessed: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
		Namespace: "dex",
		Subsystem: "execution",
		Name:      "orders_processed",
		Help:      "Number of orders processed.",
		Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
	}, []string{}),
}

func PrometheusMetrics() *Metrics {
	return exMetrics
}
