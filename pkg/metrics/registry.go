package metrics

import (
	"sort"
)

// Registry represents a collection of metrics indexed by name
type Registry struct {
	metrics map[string]*Metric
}

// NewRegistry creates a new metric registry
func NewRegistry() *Registry {
	return &Registry{
		metrics: make(map[string]*Metric),
	}
}

// Set adds or replaces a metric in the registry
func (r *Registry) Set(name string, metric *Metric) {
	r.metrics[name] = metric
}

// Get retrieves a metric from the registry
func (r *Registry) Get(name string) *Metric {
	return r.metrics[name]
}

// Has checks if a metric exists in the registry
func (r *Registry) Has(name string) bool {
	_, exists := r.metrics[name]
	return exists
}

// List returns a sorted list of all metric names
func (r *Registry) List() []string {
	list := make([]string, 0, len(r.metrics))
	for k := range r.metrics {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}

// ListByType returns metrics filtered by type
func (r *Registry) ListByType(metricType string) []*Metric {
	var result []*Metric
	for _, name := range r.List() {
		metric := r.Get(name)
		if metric.Type() == metricType {
			result = append(result, metric)
		}
	}
	return result
}

// Count returns the number of metrics in the registry
func (r *Registry) Count() int {
	return len(r.metrics)
}

// ForEach executes a function for each metric in the registry
func (r *Registry) ForEach(fn func(name string, metric *Metric)) {
	for _, name := range r.List() {
		fn(name, r.metrics[name])
	}
}

// Filter returns a new registry with metrics that match the filter function
func (r *Registry) Filter(fn func(name string, metric *Metric) bool) *Registry {
	result := NewRegistry()
	r.ForEach(func(name string, metric *Metric) {
		if fn(name, metric) {
			result.Set(name, metric)
		}
	})
	return result
}

// ListByVendor returns metrics filtered by vendor
func (r *Registry) ListByVendor(vendor string) []*Metric {
	var result []*Metric
	for _, name := range r.List() {
		metric := r.Get(name)
		if metric.Vendor() == vendor {
			result = append(result, metric)
		}
	}
	return result
}

// ListVendors returns a list of all vendors found in the metrics
func (r *Registry) ListVendors() []string {
	vendors := make(map[string]bool)
	r.ForEach(func(name string, metric *Metric) {
		if metric.Vendor() != "" {
			vendors[metric.Vendor()] = true
		}
	})
	
	result := make([]string, 0, len(vendors))
	for vendor := range vendors {
		result = append(result, vendor)
	}
	sort.Strings(result)
	return result
}

// FilterByVendor returns a new registry with metrics that match the vendor
func (r *Registry) FilterByVendor(vendor string) *Registry {
	return r.Filter(func(name string, metric *Metric) bool {
		return metric.Vendor() == vendor
	})
}