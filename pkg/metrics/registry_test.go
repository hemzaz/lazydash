package metrics

import (
	"reflect"
	"testing"
)

func TestRegistryNew(t *testing.T) {
	registry := NewRegistry()
	
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	
	if registry.Count() != 0 {
		t.Errorf("Expected empty registry, got count %d", registry.Count())
	}
}

func TestRegistrySetGet(t *testing.T) {
	registry := NewRegistry()
	
	// Create a test metric
	metric := New("test_metric", "help", nil, "counter", "", "")
	
	// Set the metric
	registry.Set(metric.Name(), metric)
	
	// Check count
	if registry.Count() != 1 {
		t.Errorf("Expected count 1, got %d", registry.Count())
	}
	
	// Get the metric
	retrieved := registry.Get(metric.Name())
	if retrieved == nil {
		t.Fatalf("Get(%q) returned nil", metric.Name())
	}
	
	// Check it's the same metric
	if retrieved.Name() != metric.Name() {
		t.Errorf("Retrieved metric name %q, expected %q", retrieved.Name(), metric.Name())
	}
}

func TestRegistryHas(t *testing.T) {
	registry := NewRegistry()
	
	metric := New("test_metric", "help", nil, "counter", "", "")
	registry.Set(metric.Name(), metric)
	
	// Check existing metric
	if !registry.Has("test_metric") {
		t.Error("Has('test_metric') returned false for existing metric")
	}
	
	// Check non-existent metric
	if registry.Has("nonexistent") {
		t.Error("Has('nonexistent') returned true for non-existent metric")
	}
}

func TestRegistryList(t *testing.T) {
	registry := NewRegistry()
	
	// Add some metrics
	metrics := []*Metric{
		New("metric1", "help", nil, "counter", "", ""),
		New("metric2", "help", nil, "gauge", "", ""),
		New("metric3", "help", nil, "counter", "", ""),
	}
	
	for _, m := range metrics {
		registry.Set(m.Name(), m)
	}
	
	// Get the list
	list := registry.List()
	
	// Check length
	if len(list) != 3 {
		t.Errorf("Expected list length 3, got %d", len(list))
	}
	
	// Check it's sorted
	expected := []string{"metric1", "metric2", "metric3"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expected sorted list %v, got %v", expected, list)
	}
}

func TestRegistryListByType(t *testing.T) {
	registry := NewRegistry()
	
	// Add metrics of different types
	registry.Set("counter1", New("counter1", "", nil, "counter", "", ""))
	registry.Set("counter2", New("counter2", "", nil, "counter", "", ""))
	registry.Set("gauge1", New("gauge1", "", nil, "gauge", "", ""))
	registry.Set("summary1", New("summary1", "", nil, "summary", "", ""))
	
	// Get counters
	counters := registry.ListByType("counter")
	if len(counters) != 2 {
		t.Errorf("Expected 2 counters, got %d", len(counters))
	}
	
	// Check counter names
	counterNames := []string{counters[0].Name(), counters[1].Name()}
	expectedNames := []string{"counter1", "counter2"}
	if !reflect.DeepEqual(counterNames, expectedNames) {
		t.Errorf("Expected counter names %v, got %v", expectedNames, counterNames)
	}
	
	// Get gauges
	gauges := registry.ListByType("gauge")
	if len(gauges) != 1 {
		t.Errorf("Expected 1 gauge, got %d", len(gauges))
	}
	
	// Get non-existent type
	none := registry.ListByType("nonexistent")
	if len(none) != 0 {
		t.Errorf("Expected 0 metrics for non-existent type, got %d", len(none))
	}
}

func TestRegistryForEach(t *testing.T) {
	registry := NewRegistry()
	
	// Add metrics
	registry.Set("metric1", New("metric1", "", nil, "", "", ""))
	registry.Set("metric2", New("metric2", "", nil, "", "", ""))
	
	// Count with ForEach
	count := 0
	registry.ForEach(func(name string, metric *Metric) {
		count++
	})
	
	if count != 2 {
		t.Errorf("ForEach visited %d metrics, expected 2", count)
	}
}

func TestRegistryFilter(t *testing.T) {
	registry := NewRegistry()
	
	// Add metrics with different names
	registry.Set("prefix_metric1", New("prefix_metric1", "", nil, "", "", ""))
	registry.Set("prefix_metric2", New("prefix_metric2", "", nil, "", "", ""))
	registry.Set("other_metric", New("other_metric", "", nil, "", "", ""))
	
	// Filter for prefix
	filtered := registry.Filter(func(name string, metric *Metric) bool {
		return len(name) >= 7 && name[:7] == "prefix_"
	})
	
	if filtered.Count() != 2 {
		t.Errorf("Expected 2 filtered metrics, got %d", filtered.Count())
	}
	
	if !filtered.Has("prefix_metric1") || !filtered.Has("prefix_metric2") {
		t.Errorf("Filtered registry is missing expected metrics")
	}
}

func TestRegistryVendorFunctions(t *testing.T) {
	registry := NewRegistry()
	
	// Create metrics with different vendors
	m1 := New("vendor1_metric1", "", nil, "", "", "")
	m1.SetVendor("vendor1")
	
	m2 := New("vendor1_metric2", "", nil, "", "", "")
	m2.SetVendor("vendor1")
	
	m3 := New("vendor2_metric1", "", nil, "", "", "")
	m3.SetVendor("vendor2")
	
	m4 := New("no_vendor_metric", "", nil, "", "", "")
	
	// Add to registry
	registry.Set(m1.Name(), m1)
	registry.Set(m2.Name(), m2)
	registry.Set(m3.Name(), m3)
	registry.Set(m4.Name(), m4)
	
	// Test ListByVendor
	vendor1Metrics := registry.ListByVendor("vendor1")
	if len(vendor1Metrics) != 2 {
		t.Errorf("Expected 2 vendor1 metrics, got %d", len(vendor1Metrics))
	}
	
	// Test ListVendors
	vendors := registry.ListVendors()
	expectedVendors := []string{"vendor1", "vendor2"}
	if !reflect.DeepEqual(vendors, expectedVendors) {
		t.Errorf("Expected vendors %v, got %v", expectedVendors, vendors)
	}
	
	// Test FilterByVendor
	vendor2Registry := registry.FilterByVendor("vendor2")
	if vendor2Registry.Count() != 1 {
		t.Errorf("Expected 1 vendor2 metric, got %d", vendor2Registry.Count())
	}
	if !vendor2Registry.Has("vendor2_metric1") {
		t.Errorf("FilterByVendor missing expected metric")
	}
}