// Package metrics provides types and functions for handling metrics
package metrics

import (
	"sort"
)

// Metric represents a single metric with metadata
type Metric struct {
	help        string
	mtype       string
	name        string
	suffix      string
	labels      map[string]bool
	unit        string
	vendor      string        // Identified vendor prefix (juniper, cisco, etc.)
	subsystem   string        // Subsystem identified from metric name
	category    string        // Category for grouping related metrics
	displayName string        // Optional display name for the metric (used for better UI)
}

// New creates a new metric with initial values
func New(name, help string, labels map[string]bool, mtype, suffix, unit string) *Metric {
	m := &Metric{
		name:        name,
		help:        help,
		mtype:       mtype,
		suffix:      suffix,
		unit:        unit,
		vendor:      "",
		subsystem:   "",
		category:    "",
		displayName: "",
	}
	
	if labels == nil {
		labels = make(map[string]bool)
	}
	m.labels = labels

	return m
}

// Help returns the help text for this metric
func (m *Metric) Help() string {
	return m.help
}

// SetHelp sets the help text for this metric
func (m *Metric) SetHelp(help string) {
	m.help = help
}

// Type returns the metric type (counter, gauge, summary)
func (m *Metric) Type() string {
	return m.mtype
}

// SetType sets the metric type
func (m *Metric) SetType(mtype string) {
	m.mtype = mtype
}

// Name returns the metric name
func (m *Metric) Name() string {
	return m.name
}

// SetName sets the metric name
func (m *Metric) SetName(name string) {
	m.name = name
}

// Unit returns the metric unit
func (m *Metric) Unit() string {
	return m.unit
}

// SetUnit sets the metric unit
func (m *Metric) SetUnit(unit string) {
	m.unit = unit
}

// Suffix returns the metric name suffix
func (m *Metric) Suffix() string {
	return m.suffix
}

// SetSuffix sets the metric name suffix
func (m *Metric) SetSuffix(suffix string) {
	m.suffix = suffix
}

// Labels returns a sorted list of label names
func (m *Metric) Labels() []string {
	if len(m.labels) == 0 {
		return nil
	}
	
	labels := make([]string, 0, len(m.labels))
	for label := range m.labels {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	
	return labels
}

// SetLabels sets the entire labels map
func (m *Metric) SetLabels(labels map[string]bool) {
	m.labels = labels
}

// AddLabel adds a new label
func (m *Metric) AddLabel(label string) {
	if m.labels == nil {
		m.labels = make(map[string]bool)
	}
	m.labels[label] = true
}

// HasLabel checks if a specific label exists
func (m *Metric) HasLabel(label string) bool {
	if m.labels == nil {
		return false
	}
	_, exists := m.labels[label]
	return exists
}

// LabelCount returns the number of labels
func (m *Metric) LabelCount() int {
	return len(m.labels)
}

// FullName returns the complete metric name including suffix
func (m *Metric) FullName() string {
	if m.suffix == "" {
		return m.name
	}
	return m.name + m.suffix
}

// Vendor returns the identified vendor prefix
func (m *Metric) Vendor() string {
	return m.vendor
}

// SetVendor sets the vendor prefix
func (m *Metric) SetVendor(vendor string) {
	m.vendor = vendor
}

// Subsystem returns the identified metric subsystem
func (m *Metric) Subsystem() string {
	return m.subsystem
}

// SetSubsystem sets the metric subsystem
func (m *Metric) SetSubsystem(subsystem string) {
	m.subsystem = subsystem
}

// Category returns the metric category
func (m *Metric) Category() string {
	return m.category
}

// SetCategory sets the metric category
func (m *Metric) SetCategory(category string) {
	m.category = category
}

// DisplayName returns the display name for this metric
func (m *Metric) DisplayName() string {
	if m.displayName != "" {
		return m.displayName
	}
	return m.name
}

// SetDisplayName sets a custom display name for this metric
func (m *Metric) SetDisplayName(displayName string) {
	m.displayName = displayName
}