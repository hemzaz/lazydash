package metrics

import (
	"reflect"
	"testing"
)

func TestMetricNew(t *testing.T) {
	name := "test_metric"
	help := "Test metric help"
	labels := map[string]bool{"label1": true, "label2": true}
	mtype := "counter"
	suffix := "_total"
	unit := "bytes"

	metric := New(name, help, labels, mtype, suffix, unit)

	if metric.Name() != name {
		t.Errorf("Expected name %q, got %q", name, metric.Name())
	}
	if metric.Help() != help {
		t.Errorf("Expected help %q, got %q", help, metric.Help())
	}
	if metric.Type() != mtype {
		t.Errorf("Expected type %q, got %q", mtype, metric.Type())
	}
	if metric.Suffix() != suffix {
		t.Errorf("Expected suffix %q, got %q", suffix, metric.Suffix())
	}
	if metric.Unit() != unit {
		t.Errorf("Expected unit %q, got %q", unit, metric.Unit())
	}
}

func TestMetricNewWithNilLabels(t *testing.T) {
	metric := New("test", "help", nil, "counter", "", "")
	
	if metric.labels == nil {
		t.Error("Expected non-nil labels map when initialized with nil")
	}
	if metric.LabelCount() != 0 {
		t.Errorf("Expected 0 labels, got %d", metric.LabelCount())
	}
}

func TestMetricLabels(t *testing.T) {
	metric := New("test", "help", nil, "counter", "", "")
	
	// Add labels
	metric.AddLabel("label1")
	metric.AddLabel("label2")
	metric.AddLabel("label3")
	
	// Check label count
	if count := metric.LabelCount(); count != 3 {
		t.Errorf("Expected 3 labels, got %d", count)
	}
	
	// Check specific label
	if !metric.HasLabel("label2") {
		t.Error("Expected metric to have 'label2'")
	}
	if metric.HasLabel("nonexistent") {
		t.Error("Expected metric to not have 'nonexistent'")
	}
	
	// Check labels are sorted
	expected := []string{"label1", "label2", "label3"}
	actual := metric.Labels()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected labels %v, got %v", expected, actual)
	}
}

func TestMetricSetters(t *testing.T) {
	metric := New("test", "help", nil, "counter", "", "")
	
	// Test setters
	metric.SetName("new_name")
	if metric.Name() != "new_name" {
		t.Errorf("Expected name 'new_name', got %q", metric.Name())
	}
	
	metric.SetHelp("new help")
	if metric.Help() != "new help" {
		t.Errorf("Expected help 'new help', got %q", metric.Help())
	}
	
	metric.SetType("gauge")
	if metric.Type() != "gauge" {
		t.Errorf("Expected type 'gauge', got %q", metric.Type())
	}
	
	metric.SetSuffix("_suffix")
	if metric.Suffix() != "_suffix" {
		t.Errorf("Expected suffix '_suffix', got %q", metric.Suffix())
	}
	
	metric.SetUnit("seconds")
	if metric.Unit() != "seconds" {
		t.Errorf("Expected unit 'seconds', got %q", metric.Unit())
	}
	
	// Test label setters
	newLabels := map[string]bool{"replaced": true}
	metric.SetLabels(newLabels)
	if metric.LabelCount() != 1 || !metric.HasLabel("replaced") {
		t.Errorf("SetLabels did not replace labels correctly")
	}
}

func TestMetricFullName(t *testing.T) {
	// Without suffix
	metric1 := New("test_metric", "", nil, "", "", "")
	if metric1.FullName() != "test_metric" {
		t.Errorf("Expected 'test_metric', got %q", metric1.FullName())
	}
	
	// With suffix
	metric2 := New("test_metric", "", nil, "", "_total", "")
	if metric2.FullName() != "test_metric_total" {
		t.Errorf("Expected 'test_metric_total', got %q", metric2.FullName())
	}
}

func TestMetricVendorSubsystemCategory(t *testing.T) {
	metric := New("test", "", nil, "", "", "")
	
	// Test vendor
	metric.SetVendor("juniper")
	if metric.Vendor() != "juniper" {
		t.Errorf("Expected vendor 'juniper', got %q", metric.Vendor())
	}
	
	// Test subsystem
	metric.SetSubsystem("networking")
	if metric.Subsystem() != "networking" {
		t.Errorf("Expected subsystem 'networking', got %q", metric.Subsystem())
	}
	
	// Test category
	metric.SetCategory("hardware")
	if metric.Category() != "hardware" {
		t.Errorf("Expected category 'hardware', got %q", metric.Category())
	}
}

func TestMetricDisplayName(t *testing.T) {
	// Default is metric name
	metric := New("test_metric", "", nil, "", "", "")
	if metric.DisplayName() != "test_metric" {
		t.Errorf("Expected display name 'test_metric', got %q", metric.DisplayName())
	}
	
	// Custom display name
	metric.SetDisplayName("Test Metric")
	if metric.DisplayName() != "Test Metric" {
		t.Errorf("Expected display name 'Test Metric', got %q", metric.DisplayName())
	}
}