package query

import (
	"fmt"
	"strings"
)

// PromQLBuilder helps build advanced PromQL queries
type PromQLBuilder struct {
	metric    string
	selectors map[string]string
	functions []string
	timeRange string
	groupBy   []string
	offset    string
}

// NewPromQLBuilder creates a new PromQL query builder
func NewPromQLBuilder(metric string) *PromQLBuilder {
	return &PromQLBuilder{
		metric:    metric,
		selectors: make(map[string]string),
		groupBy:   []string{},
	}
}

// WithLabel adds a label selector to the query
func (b *PromQLBuilder) WithLabel(key, value string) *PromQLBuilder {
	b.selectors[key] = value
	return b
}

// WithLabels adds multiple label selectors at once
func (b *PromQLBuilder) WithLabels(labels map[string]string) *PromQLBuilder {
	for k, v := range labels {
		b.selectors[k] = v
	}
	return b
}

// WithFunction wraps the query with a function like rate, sum, avg, etc.
func (b *PromQLBuilder) WithFunction(fn string, args ...interface{}) *PromQLBuilder {
	fnStr := fn
	if len(args) > 0 {
		argStrings := make([]string, len(args))
		for i, arg := range args {
			argStrings[i] = fmt.Sprintf("%v", arg)
		}
		fnStr = fmt.Sprintf("%s(%s)", fn, strings.Join(argStrings, ","))
	}
	b.functions = append(b.functions, fnStr)
	return b
}

// WithRate adds a rate() function with the given time range
func (b *PromQLBuilder) WithRate(duration string) *PromQLBuilder {
	b.timeRange = duration
	return b.WithFunction("rate")
}

// WithGroupBy adds group by clause to the query
func (b *PromQLBuilder) WithGroupBy(labels ...string) *PromQLBuilder {
	b.groupBy = append(b.groupBy, labels...)
	return b
}

// WithOffset adds an offset to the query
func (b *PromQLBuilder) WithOffset(offset string) *PromQLBuilder {
	b.offset = offset
	return b
}

// Build constructs the final PromQL query string
func (b *PromQLBuilder) Build() string {
	// Start with the metric name
	query := b.metric
	
	// Add label selectors if any
	if len(b.selectors) > 0 {
		selectorParts := make([]string, 0, len(b.selectors))
		for k, v := range b.selectors {
			selectorParts = append(selectorParts, fmt.Sprintf("%s=\"%s\"", k, v))
		}
		query = fmt.Sprintf("%s{%s}", query, strings.Join(selectorParts, ","))
	} else {
		query = fmt.Sprintf("%s{}", query)
	}
	
	// Add time range for rate functions
	if b.timeRange != "" {
		query = fmt.Sprintf("%s[%s]", query, b.timeRange)
	}
	
	// Add offset if specified
	if b.offset != "" {
		query = fmt.Sprintf("%s offset %s", query, b.offset)
	}
	
	// Apply functions in reverse order (innermost first)
	for i := len(b.functions) - 1; i >= 0; i-- {
		fn := b.functions[i]
		if strings.HasSuffix(fn, ")") {
			// Function already has arguments
			query = strings.Replace(fn, ")", fmt.Sprintf("(%s))", query), 1)
		} else {
			query = fmt.Sprintf("%s(%s)", fn, query)
		}
	}
	
	// Add group by if specified
	if len(b.groupBy) > 0 {
		query = fmt.Sprintf("%s by (%s)", query, strings.Join(b.groupBy, ","))
	}
	
	return query
}

// Common PromQL patterns
func BuildCounterRateQuery(metric string, timeRange string, labels map[string]string) string {
	builder := NewPromQLBuilder(metric)
	for k, v := range labels {
		builder.WithLabel(k, v)
	}
	return builder.WithRate(timeRange).WithFunction("sum").Build()
}

func BuildGaugeQuery(metric string, labels map[string]string, groupBy []string) string {
	builder := NewPromQLBuilder(metric)
	for k, v := range labels {
		builder.WithLabel(k, v)
	}
	if len(groupBy) > 0 {
		builder.WithFunction("avg").WithGroupBy(groupBy...)
	}
	return builder.Build()
}

func BuildHistogramQuery(metric string, percentile float64, timeRange string) string {
	return NewPromQLBuilder(metric + "_bucket").
		WithRate(timeRange).
		WithFunction("histogram_quantile", percentile).
		Build()
}

func BuildErrorRateQuery(errorMetric, totalMetric, timeRange string) string {
	errorQuery := NewPromQLBuilder(errorMetric).WithRate(timeRange).WithFunction("sum").Build()
	totalQuery := NewPromQLBuilder(totalMetric).WithRate(timeRange).WithFunction("sum").Build()
	
	return fmt.Sprintf("(%s / %s) * 100", errorQuery, totalQuery)
}