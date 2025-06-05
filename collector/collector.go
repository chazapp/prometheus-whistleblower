package collector

import (
	"errors"
	"slices"

	"github.com/prometheus/client_golang/prometheus"
)

type Label struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Metric struct {
	Metric string  `json:"metric"`
	Labels []Label `json:"labels"`
	Value  int     `json:"value"`
	ID     int     `json:"id",omitempty`
}

type WhistleblowerCollector struct {
	Metrics     []Metric
	Descriptors []*prometheus.Desc
}

func NewWhistleblowerCollector() WhistleblowerCollector {
	return WhistleblowerCollector{
		Metrics:     make([]Metric, 0),
		Descriptors: make([]*prometheus.Desc, 0),
	}
}

func (w *WhistleblowerCollector) Collect(ch chan<- prometheus.Metric) {
	for idx, item := range w.Metrics {
		labelValues := func(labels []Label) []string {
			lv := make([]string, len(labels))
			for i, l := range labels {
				lv[i] = l.Value
			}
			return lv
		}(item.Labels)
		ch <- prometheus.MustNewConstMetric(w.Descriptors[idx], prometheus.CounterValue, float64(item.Value), labelValues...)
	}
}

func (w *WhistleblowerCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, item := range w.Descriptors {
		ch <- item
	}
}

func (w *WhistleblowerCollector) AddMetric(metric Metric) {
	metric.ID = len(w.Metrics)
	w.Metrics = append(w.Metrics, metric)
	labels := func(labels []Label) []string {
		lv := make([]string, len(labels))
		for i, l := range labels {
			lv[i] = l.Label
		}
		return lv
	}(metric.Labels)
	w.Descriptors = append(w.Descriptors, prometheus.NewDesc(metric.Metric, "", labels, nil))
}

func (w *WhistleblowerCollector) DeleteMetric(id int) error {
	for idx, item := range w.Metrics {
		if item.ID == id {
			w.Metrics = slices.Delete(w.Metrics, idx, idx+1)
			w.Descriptors = slices.Delete(w.Descriptors, idx, idx+1)
			return nil
		}
	}
	return errors.New("unknown metric ID")
}
