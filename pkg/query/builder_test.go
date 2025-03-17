package query

import (
	"testing"
)

func TestNewPromQLBuilder(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	if builder == nil {
		t.Fatal("NewPromQLBuilder returned nil")
	}
	
	if builder.metric != "test_metric" {
		t.Errorf("Expected metric name to be 'test_metric', got %q", builder.metric)
	}
	
	if builder.selectors == nil {
		t.Error("Expected selectors map to be initialized")
	}
	
	if len(builder.functions) != 0 {
		t.Errorf("Expected functions to be empty, got %d items", len(builder.functions))
	}
	
	if len(builder.groupBy) != 0 {
		t.Errorf("Expected groupBy to be empty, got %d items", len(builder.groupBy))
	}
}

func TestWithLabel(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	// Add a label
	builder.WithLabel("instance", "localhost:9090")
	
	// Check the label was added
	if val, ok := builder.selectors["instance"]; !ok || val != "localhost:9090" {
		t.Errorf("Expected label 'instance' with value 'localhost:9090', got %v", builder.selectors)
	}
	
	// Check method chaining
	returnedBuilder := builder.WithLabel("job", "prometheus")
	if returnedBuilder != builder {
		t.Error("WithLabel did not return the builder for chaining")
	}
}

func TestWithLabels(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	// Add multiple labels
	labels := map[string]string{
		"instance": "localhost:9090",
		"job":      "prometheus",
		"env":      "production",
	}
	
	builder.WithLabels(labels)
	
	// Check all labels were added
	if len(builder.selectors) != len(labels) {
		t.Errorf("Expected %d labels, got %d", len(labels), len(builder.selectors))
	}
	
	for k, v := range labels {
		if val, ok := builder.selectors[k]; !ok || val != v {
			t.Errorf("Expected label %q with value %q, got %q", k, v, val)
		}
	}
}

func TestWithFunction(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	// Add a function without args
	builder.WithFunction("sum")
	
	if len(builder.functions) != 1 || builder.functions[0] != "sum" {
		t.Errorf("Expected functions to contain 'sum', got %v", builder.functions)
	}
	
	// Add a function with args
	builder.WithFunction("quantile", 0.95)
	
	if len(builder.functions) != 2 || builder.functions[1] != "quantile(0.95)" {
		t.Errorf("Expected functions to contain 'quantile(0.95)', got %v", builder.functions)
	}
	
	// Add a function with multiple args
	builder.WithFunction("histogram_quantile", 0.99, "some_metric")
	
	if len(builder.functions) != 3 || builder.functions[2] != "histogram_quantile(0.99,some_metric)" {
		t.Errorf("Expected functions to contain 'histogram_quantile(0.99,some_metric)', got %v", builder.functions)
	}
}

func TestWithRate(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	// Add rate with time range
	builder.WithRate("5m")
	
	if builder.timeRange != "5m" {
		t.Errorf("Expected timeRange to be '5m', got %q", builder.timeRange)
	}
	
	if len(builder.functions) != 1 || builder.functions[0] != "rate" {
		t.Errorf("Expected functions to contain 'rate', got %v", builder.functions)
	}
}

func TestWithGroupBy(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	// Add group by with single label
	builder.WithGroupBy("instance")
	
	if len(builder.groupBy) != 1 || builder.groupBy[0] != "instance" {
		t.Errorf("Expected groupBy to contain 'instance', got %v", builder.groupBy)
	}
	
	// Add group by with multiple labels
	builder.WithGroupBy("job", "env")
	
	if len(builder.groupBy) != 3 || builder.groupBy[1] != "job" || builder.groupBy[2] != "env" {
		t.Errorf("Expected groupBy to contain ['instance', 'job', 'env'], got %v", builder.groupBy)
	}
}

func TestWithOffset(t *testing.T) {
	builder := NewPromQLBuilder("test_metric")
	
	// Add offset
	builder.WithOffset("1h")
	
	if builder.offset != "1h" {
		t.Errorf("Expected offset to be '1h', got %q", builder.offset)
	}
}

func TestBuild(t *testing.T) {
	// Test simple metric
	t.Run("Simple metric", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		query := builder.Build()
		
		expected := "test_metric{}"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test with labels
	t.Run("With labels", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		builder.WithLabel("instance", "localhost:9090")
		builder.WithLabel("job", "prometheus")
		query := builder.Build()
		
		// We can't guarantee the order of labels, so check both possibilities
		expected1 := "test_metric{instance=\"localhost:9090\",job=\"prometheus\"}"
		expected2 := "test_metric{job=\"prometheus\",instance=\"localhost:9090\"}"
		
		if query != expected1 && query != expected2 {
			t.Errorf("Expected query to be either %q or %q, got %q", expected1, expected2, query)
		}
	})
	
	// Test with time range
	t.Run("With time range", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		builder.timeRange = "5m"
		query := builder.Build()
		
		expected := "test_metric{}[5m]"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test with offset
	t.Run("With offset", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		builder.WithOffset("1h")
		query := builder.Build()
		
		expected := "test_metric{} offset 1h"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test with function
	t.Run("With function", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		builder.WithFunction("sum")
		query := builder.Build()
		
		expected := "sum(test_metric{})"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test with multiple functions
	t.Run("With multiple functions", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		builder.WithRate("5m")
		builder.WithFunction("sum")
		query := builder.Build()
		
		expected := "rate(sum(test_metric{}[5m]))"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test with group by
	t.Run("With group by", func(t *testing.T) {
		builder := NewPromQLBuilder("test_metric")
		builder.WithFunction("sum")
		builder.WithGroupBy("instance")
		query := builder.Build()
		
		expected := "sum(test_metric{}) by (instance)"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test complex query
	t.Run("Complex query", func(t *testing.T) {
		builder := NewPromQLBuilder("http_requests_total")
		builder.WithLabel("job", "api-server")
		builder.WithRate("5m")
		builder.WithFunction("sum")
		builder.WithGroupBy("instance", "method")
		
		query := builder.Build()
		
		expected := "rate(sum(http_requests_total{job=\"api-server\"}[5m])) by (instance,method)"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	// Test BuildCounterRateQuery
	t.Run("BuildCounterRateQuery", func(t *testing.T) {
		labels := map[string]string{"job": "api-server"}
		query := BuildCounterRateQuery("http_requests_total", "5m", labels)
		
		expected := "rate(sum(http_requests_total{job=\"api-server\"}[5m]))"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test BuildGaugeQuery
	t.Run("BuildGaugeQuery without groupBy", func(t *testing.T) {
		labels := map[string]string{"job": "node-exporter"}
		query := BuildGaugeQuery("node_memory_Active_bytes", labels, nil)
		
		expected := "node_memory_Active_bytes{job=\"node-exporter\"}"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	t.Run("BuildGaugeQuery with groupBy", func(t *testing.T) {
		labels := map[string]string{"job": "node-exporter"}
		groupBy := []string{"instance"}
		query := BuildGaugeQuery("node_memory_Active_bytes", labels, groupBy)
		
		expected := "avg(node_memory_Active_bytes{job=\"node-exporter\"}) by (instance)"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test BuildHistogramQuery
	t.Run("BuildHistogramQuery", func(t *testing.T) {
		query := BuildHistogramQuery("http_request_duration_seconds", 0.95, "5m")
		
		expected := "rate(histogram_quantile(0.95(http_request_duration_seconds_bucket{}[5m])))"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
	
	// Test BuildErrorRateQuery
	t.Run("BuildErrorRateQuery", func(t *testing.T) {
		query := BuildErrorRateQuery("http_requests_errors_total", "http_requests_total", "5m")
		
		expected := "(rate(sum(http_requests_errors_total{}[5m])) / rate(sum(http_requests_total{}[5m]))) * 100"
		if query != expected {
			t.Errorf("Expected query %q, got %q", expected, query)
		}
	})
}