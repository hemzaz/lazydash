// Package query provides query building functionality for Prometheus metrics
package query

import (
	"fmt"
	"strings"

	"github.com/hemzaz/lazydash/internal/config"
	"github.com/hemzaz/lazydash/pkg/metrics"
)

// Builder creates PromQL queries for different metric types
type Builder struct {
	config *config.Config
}

// NewBuilder creates a new PromQL query builder with configuration
func NewBuilder(cfg *config.Config) *Builder {
	return &Builder{
		config: cfg,
	}
}

// BuildQuery creates a PromQL query for a metric based on its type
func (b *Builder) BuildQuery(metric *metrics.Metric) string {
	// Special handling for JTIMON metrics from Juniper devices
	if metric.Vendor() == "juniper" && strings.HasPrefix(metric.Name(), "_") {
		return b.buildJuniperJtimonQuery(metric)
	}

	switch metric.Type() {
	case "counter":
		return b.buildCounterQuery(metric)
	case "gauge":
		return b.buildGaugeQuery(metric)
	case "summary":
		return b.buildSummaryQuery(metric)
	default:
		return metric.FullName()
	}
}

// buildCounterQuery builds a rate-based query for counter metrics
func (b *Builder) buildCounterQuery(metric *metrics.Metric) string {
	tmpl := b.config.CounterExprTmpl
	metricName := metric.Name() + metric.Suffix()
	return strings.Replace(tmpl, b.config.Delimiter, metricName, -1)
}

// buildGaugeQuery builds a query for gauge metrics
func (b *Builder) buildGaugeQuery(metric *metrics.Metric) string {
	tmpl := b.config.GaugeExprTmpl
	metricName := metric.Name() + metric.Suffix()
	return strings.Replace(tmpl, b.config.Delimiter, metricName, -1)
}

// buildSummaryQuery builds a query for summary metrics
func (b *Builder) buildSummaryQuery(metric *metrics.Metric) string {
	tmpl := b.config.SummaryExprTmpl
	return strings.Replace(tmpl, b.config.Delimiter, metric.Name(), -1)
}

// FormatLegend creates a legend format string based on metric labels
func (b *Builder) FormatLegend(metric *metrics.Metric, fallback string) string {
	labels := metric.Labels()
	
	if len(labels) == 0 {
		if fallback != "" {
			return fallback
		}
		return "Job:[{{job}}]"
	}

	parts := make([]string, 0, len(labels))
	for _, label := range labels {
		parts = append(parts, fmt.Sprintf("%s:[{{%s}}]", label, label))
	}
	
	return strings.Join(parts, " ")
}

// buildJuniperJtimonQuery builds special queries for Juniper JTIMON metrics
func (b *Builder) buildJuniperJtimonQuery(metric *metrics.Metric) string {
	metricName := metric.Name()
	
	// Different query logic based on the type of JTIMON metric
	if strings.Contains(metricName, "_counters_") {
		// For counter-type metrics, use rate to show changes over time
		if strings.Contains(metricName, "_in_") || strings.Contains(metricName, "_out_") || 
		   strings.Contains(metricName, "_frames_") || strings.Contains(metricName, "_drops_") ||
		   strings.Contains(metricName, "_errors_") || strings.Contains(metricName, "_packets_") {
			// Use rate for counter metrics
			return fmt.Sprintf("rate(%s[5m])", metricName)
		}
	}
	
	// For CPU, memory, and utilization metrics, just show the value
	if strings.Contains(metricName, "_cpu_") || strings.Contains(metricName, "_memory_") ||
	   strings.Contains(metricName, "_utilization_") || strings.Contains(metricName, "_temperature_") {
		// For certain metrics, filter on device to get per-device values
		return fmt.Sprintf("%s", metricName)
	}
	
	// Default query for other JTIMON metrics
	return metricName
}

// GetLegend returns the appropriate legend format for a metric type
func (b *Builder) GetLegend(metric *metrics.Metric) string {
	// Special legend format for Juniper JTIMON metrics
	if metric.Vendor() == "juniper" && strings.HasPrefix(metric.Name(), "_") {
		// Build a legend format using the most relevant labels
		labels := metric.Labels()
		if len(labels) > 0 {
			// Priority labels to include
			// Try to include device name and component/interface name in the legend
			priorityLabels := []string{"device", "_interfaces_interface__name", "_components_component__name"}
			legends := []string{}
			
			for _, label := range priorityLabels {
				for _, metricLabel := range labels {
					if metricLabel == label {
						legends = append(legends, fmt.Sprintf("%s:[{{%s}}]", metricLabel, metricLabel))
						break
					}
				}
			}
			
			if len(legends) > 0 {
				return strings.Join(legends, " ")
			}
		}
		
		// Default for JTIMON metrics 
		return "Device:[{{device}}]"
	}

	switch metric.Type() {
	case "counter":
		return b.FormatLegend(metric, b.config.CounterLegend)
	case "gauge":
		return b.FormatLegend(metric, b.config.GaugeLegend)
	case "summary":
		return b.FormatLegend(metric, b.config.SummaryLegend)
	default:
		return b.FormatLegend(metric, "")
	}
}