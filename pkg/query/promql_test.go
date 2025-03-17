package query

import (
	"testing"

	"github.com/hemzaz/lazydash/internal/config"
	"github.com/hemzaz/lazydash/pkg/metrics"
)

func TestNewBuilder(t *testing.T) {
	cfg := config.New()
	builder := NewBuilder(cfg)
	
	if builder == nil {
		t.Fatal("NewBuilder returned nil")
	}
	
	if builder.config != cfg {
		t.Error("Builder config not set correctly")
	}
}

func TestBuildQuery(t *testing.T) {
	cfg := config.New()
	builder := NewBuilder(cfg)
	
	// Test counter metric
	t.Run("Counter metric", func(t *testing.T) {
		metric := metrics.New("http_requests_total", "Counter of HTTP requests", nil, "counter", "_total", "")
		query := builder.BuildQuery(metric)
		
		expected := "sum(rate(http_requests_total_total [1m]))"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test gauge metric
	t.Run("Gauge metric", func(t *testing.T) {
		metric := metrics.New("node_memory_usage", "Memory usage", nil, "gauge", "", "bytes")
		query := builder.BuildQuery(metric)
		
		expected := "node_memory_usage"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test summary metric
	t.Run("Summary metric", func(t *testing.T) {
		metric := metrics.New("http_request_duration_seconds", "HTTP request duration", nil, "summary", "", "seconds")
		query := builder.BuildQuery(metric)
		
		expected := "http_request_duration_seconds"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test unknown type
	t.Run("Unknown type", func(t *testing.T) {
		metric := metrics.New("custom_metric", "Custom metric", nil, "unknown", "", "")
		query := builder.BuildQuery(metric)
		
		expected := "custom_metric"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test Juniper metric
	t.Run("Juniper metric", func(t *testing.T) {
		metric := metrics.New("_juniper_interfaces_counters_in_octets", "Juniper interface counters", nil, "counter", "", "")
		metric.SetVendor("juniper")
		query := builder.BuildQuery(metric)
		
		expected := "rate(_juniper_interfaces_counters_in_octets[5m])"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test Juniper CPU metric
	t.Run("Juniper CPU metric", func(t *testing.T) {
		metric := metrics.New("_juniper_system_cpu_utilization", "Juniper CPU utilization", nil, "gauge", "", "")
		metric.SetVendor("juniper")
		query := builder.BuildQuery(metric)
		
		expected := "_juniper_system_cpu_utilization"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
}

func TestBuildCounterQuery(t *testing.T) {
	cfg := config.New()
	builder := NewBuilder(cfg)
	
	// Default counter expression
	t.Run("Default counter expression", func(t *testing.T) {
		metric := metrics.New("http_requests_total", "Counter of HTTP requests", nil, "counter", "_total", "")
		query := builder.buildCounterQuery(metric)
		
		expected := "sum(rate(http_requests_total_total [1m]))"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Custom counter expression
	t.Run("Custom counter expression", func(t *testing.T) {
		cfg.CounterExprTmpl = "max(rate(:METRIC:[5m]))"
		builder := NewBuilder(cfg)
		
		metric := metrics.New("http_requests_total", "Counter of HTTP requests", nil, "counter", "_total", "")
		query := builder.buildCounterQuery(metric)
		
		expected := "max(rate(http_requests_total_total[5m]))"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
}

func TestBuildGaugeQuery(t *testing.T) {
	cfg := config.New()
	builder := NewBuilder(cfg)
	
	// Default gauge expression
	t.Run("Default gauge expression", func(t *testing.T) {
		metric := metrics.New("node_memory_usage", "Memory usage", nil, "gauge", "", "bytes")
		query := builder.buildGaugeQuery(metric)
		
		expected := "node_memory_usage"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Custom gauge expression
	t.Run("Custom gauge expression", func(t *testing.T) {
		cfg.GaugeExprTmpl = "avg(:METRIC:)"
		builder := NewBuilder(cfg)
		
		metric := metrics.New("node_memory_usage", "Memory usage", nil, "gauge", "", "bytes")
		query := builder.buildGaugeQuery(metric)
		
		expected := "avg(node_memory_usage)"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
}

func TestFormatLegend(t *testing.T) {
	cfg := config.New()
	builder := NewBuilder(cfg)
	
	// No labels, using fallback
	t.Run("No labels, using fallback", func(t *testing.T) {
		metric := metrics.New("test_metric", "", nil, "", "", "")
		legend := builder.FormatLegend(metric, "Fallback")
		
		expected := "Fallback"
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
	
	// No labels, no fallback
	t.Run("No labels, no fallback", func(t *testing.T) {
		metric := metrics.New("test_metric", "", nil, "", "", "")
		legend := builder.FormatLegend(metric, "")
		
		expected := "Job:[{{job}}]"
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
	
	// With labels
	t.Run("With labels", func(t *testing.T) {
		metric := metrics.New("test_metric", "", nil, "", "", "")
		metric.AddLabel("instance")
		metric.AddLabel("job")
		legend := builder.FormatLegend(metric, "")
		
		// We can't predict the exact order of labels
		if legend != "instance:[{{instance}}] job:[{{job}}]" && 
		   legend != "job:[{{job}}] instance:[{{instance}}]" {
			t.Errorf("Expected legend with instance and job, got %q", legend)
		}
	})
}

func TestGetLegend(t *testing.T) {
	cfg := config.New()
	builder := NewBuilder(cfg)
	
	// Counter legend
	t.Run("Counter legend", func(t *testing.T) {
		metric := metrics.New("requests_total", "", nil, "counter", "", "")
		legend := builder.GetLegend(metric)
		
		expected := "Job:[{{job}}]" // Default counter legend
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
	
	// Gauge legend
	t.Run("Gauge legend", func(t *testing.T) {
		metric := metrics.New("memory_usage", "", nil, "gauge", "", "")
		legend := builder.GetLegend(metric)
		
		expected := "Job:[{{job}}]" // Default gauge legend
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
	
	// Summary legend
	t.Run("Summary legend", func(t *testing.T) {
		metric := metrics.New("latency", "", nil, "summary", "", "")
		legend := builder.GetLegend(metric)
		
		expected := "Job:[{{job}}]" // Default summary legend
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
	
	// Juniper legend
	t.Run("Juniper legend", func(t *testing.T) {
		metric := metrics.New("_juniper_metric", "", nil, "", "", "")
		metric.SetVendor("juniper")
		legend := builder.GetLegend(metric)
		
		expected := "Device:[{{device}}]" // Default Juniper legend
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
	
	// Juniper legend with priority labels
	t.Run("Juniper legend with priority labels", func(t *testing.T) {
		metric := metrics.New("_juniper_metric", "", nil, "", "", "")
		metric.SetVendor("juniper")
		metric.AddLabel("device")
		metric.AddLabel("_interfaces_interface__name")
		legend := builder.GetLegend(metric)
		
		expected := "device:[{{device}}] _interfaces_interface__name:[{{_interfaces_interface__name}}]"
		if legend != expected {
			t.Errorf("Expected legend %q, got %q", expected, legend)
		}
	})
}