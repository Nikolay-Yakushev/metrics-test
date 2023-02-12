package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	ReqBodySize *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer, name string) *Metrics{
    m := &Metrics{
        ReqBodySize: prometheus.NewHistogramVec(
                        prometheus.HistogramOpts{
                            Namespace: name,
                            Name:      "request_body_size",
                            Help:      "Request body size",
                            Buckets:   []float64{0, 0.25, 0.5, 0.75, 1},
                        }, 
                        []string{"method"}),

        }
    return m
}