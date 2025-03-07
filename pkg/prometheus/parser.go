// Package prometheus provides functionality for parsing Prometheus metrics
package prometheus

import (
	"io"
	"strings"

	"github.com/hemzaz/lazydash/internal/config"
	"github.com/hemzaz/lazydash/pkg/metrics"
	PromLabel "github.com/prometheus/prometheus/pkg/labels"
	PromParse "github.com/prometheus/prometheus/pkg/textparse"
	"github.com/rs/zerolog/log"
)

// ParseMetrics parses Prometheus text metrics into a metric registry
func ParseMetrics(data []byte) *metrics.Registry {
	return ParseMetricsWithConfig(data, nil)
}

// ParseMetricsWithConfig parses Prometheus text metrics with vendor-specific config
func ParseMetricsWithConfig(data []byte, cfg *config.Config) *metrics.Registry {
	p := PromParse.NewPromParser(data)
	registry := metrics.NewRegistry()
	
	// Determine if we should detect vendor prefixes
	detectVendors := false
	var knownPrefixes []string
	var customPrefixes []string
	detectJuniper := false
	detectCisco := false
	
	if cfg != nil && cfg.VendorConfig != nil && cfg.VendorConfig.Enabled {
		detectVendors = true
		knownPrefixes = cfg.VendorConfig.KnownPrefixes
		customPrefixes = cfg.VendorConfig.CustomPrefixes
		detectJuniper = cfg.VendorConfig.JuniperEnabled
		detectCisco = cfg.VendorConfig.CiscoEnabled
		
		log.Info().Bool("vendor_detection", true).
			Int("known_prefixes", len(knownPrefixes)).
			Int("custom_prefixes", len(customPrefixes)).
			Bool("juniper", detectJuniper).
			Bool("cisco", detectCisco).
			Msg("Vendor-specific metric detection enabled")
	}

	for {
		et, err := p.Next()
		if err == io.EOF {
			break
		}

		//May be parsed out of order
		switch et {
		case PromParse.EntryHelp:
			m, h := p.Help()
			registry.Set(string(m), metrics.New(string(m), string(h), nil, "", "", "short"))

		case PromParse.EntryType:
			m, typ := p.Type()
			if metric := registry.Get(string(m)); metric != nil {
				metric.SetType(string(typ))
			}

		case PromParse.EntrySeries:
			labels := &PromLabel.Labels{}
			p.Metric(labels)

			labelmap := labels.Map()
			name := labelmap["__name__"]

			// Unify metrics key for simple access
			var suffix string
			if strings.HasSuffix(name, "_bucket") {
				suffix = "_bucket"
				name = strings.TrimSuffix(name, "_bucket")
			} else if strings.HasSuffix(name, "_sum") {
				suffix = "_sum"
				name = strings.TrimSuffix(name, "_sum")
			} else if strings.HasSuffix(name, "_count") {
				suffix = "_count"
				name = strings.TrimSuffix(name, "_count")
			}

			// Create metric if it doesn't exist
			if !registry.Has(name) {
				registry.Set(name, metrics.New(name, "", nil, "", suffix, "short"))
			}
			
			metric := registry.Get(name)
			metric.SetSuffix(suffix)

			// Try to detect appropriate unit type based on metric name
			if strings.Contains(name, "_seconds") {
				metric.SetUnit("s")
			} else if strings.Contains(name, "_milliseconds") {
				metric.SetUnit("ms")
			} else if strings.Contains(name, "_bytes") {
				metric.SetUnit("decbytes")
			} else if strings.Contains(name, "_percent") || strings.Contains(name, "_ratio") {
				metric.SetUnit("percent")
			} else if strings.Contains(name, "_count") {
				metric.SetUnit("short")
			}
			
			// Detect vendor prefixes if enabled
			if detectVendors {
				detectVendorPrefix(metric, knownPrefixes, customPrefixes)
				
				// Special handling for specific vendors
				if detectJuniper && metric.Vendor() == "juniper" {
					parseJuniperMetric(metric)
				}
				
				if detectCisco && metric.Vendor() == "cisco" {
					parseCiscoMetric(metric)
				}
			}

			// Add all labels to the metric
			for k := range labelmap {
				if k != "__name__" && k != "" {
					metric.AddLabel(k)
				}
			}
		}
	}

	return registry
}

// detectVendorPrefix identifies the vendor from metric prefixes
func detectVendorPrefix(metric *metrics.Metric, knownPrefixes, customPrefixes []string) {
	name := metric.Name()
	help := metric.Help()
	
	// Check if this is a Juniper JTIMON metric
	if strings.Contains(help, "JTIMON Metric") || strings.HasPrefix(name, "_") {
		// JTIMON metrics from Juniper typically start with underscore 
		// and have structured names like _interfaces_interface_...
		metric.SetVendor("juniper")
		
		// Extract subsystem from the JTIMON metric structure
		if strings.HasPrefix(name, "_components_") {
			metric.SetSubsystem("components")
			if strings.Contains(name, "_cpu_") {
				metric.SetCategory("cpu")
			} else if strings.Contains(name, "_memory_") {
				metric.SetCategory("memory")
			} else if strings.Contains(name, "_temperature_") {
				metric.SetCategory("temperature")
			} else if strings.Contains(name, "_transceiver_") {
				metric.SetCategory("transceiver")
			} else if strings.Contains(name, "_power_") {
				metric.SetCategory("power")
			}
		} else if strings.HasPrefix(name, "_interfaces_") {
			metric.SetSubsystem("interfaces")
			if strings.Contains(name, "_ethernet_") {
				metric.SetCategory("ethernet")
			} else if strings.Contains(name, "_aggregation_") {
				metric.SetCategory("lag")
			} else if strings.Contains(name, "_counters_") {
				metric.SetCategory("counters")
			}
		}
		return
	}
	
	// Check known prefixes
	for _, prefix := range knownPrefixes {
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			vendor := strings.TrimSuffix(prefix, "_")
			metric.SetVendor(vendor)
			
			// Extract subsystem (part after vendor_ prefix before the next underscore)
			remaining := strings.TrimPrefix(name, prefix)
			parts := strings.SplitN(remaining, "_", 2)
			if len(parts) > 0 {
				metric.SetSubsystem(parts[0])
			}
			
			return
		}
	}
	
	// Check custom prefixes
	for _, prefix := range customPrefixes {
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			vendor := strings.TrimSuffix(prefix, "_")
			metric.SetVendor(vendor)
			return
		}
	}
}

// parseJuniperMetric extracts additional information from Juniper metrics
func parseJuniperMetric(metric *metrics.Metric) {
	name := metric.Name()
	subsystem := metric.Subsystem()
	
	// For JTIMON metrics which are already categorized, set units appropriately
	if strings.HasPrefix(name, "_") {
		// Set units based on metric name
		if strings.Contains(name, "_cpu_") || strings.Contains(name, "_utilization_") {
			metric.SetUnit("percent")
		} else if strings.Contains(name, "_power_") {
			metric.SetUnit("dbm")
		} else if strings.Contains(name, "_temperature_") {
			metric.SetUnit("celsius")
		} else if strings.Contains(name, "_bytes") || strings.Contains(name, "_octets") {
			metric.SetUnit("decbytes")
		} else if strings.Contains(name, "_bits_") {
			metric.SetUnit("decbits")
		} else if strings.Contains(name, "_speed") {
			metric.SetUnit("bps")
		} else if strings.Contains(name, "_packets") || strings.Contains(name, "_frames") {
			metric.SetUnit("pps")
		} else if strings.Contains(name, "_memory_") {
			metric.SetUnit("bytes")
		} else if strings.Contains(name, "_counters_") {
			metric.SetUnit("short")
		}
		
		// Interface metrics handling
		if strings.Contains(name, "_interfaces_") {
			if metric.Category() == "" {
				if strings.Contains(name, "_error") || strings.Contains(name, "_discard") {
					metric.SetCategory("errors")
				} else if strings.Contains(name, "_counters_") {
					metric.SetCategory("counters")
				} else if strings.Contains(name, "_state_") {
					metric.SetCategory("state")
				} else if strings.Contains(name, "_statistics_") {
					metric.SetCategory("statistics")
				}
			}
		}
		
		// Set better display names for metrics
		if strings.Contains(name, "_instant") {
			// Remove the _instant suffix for better display
			displayName := strings.Replace(name, "_instant", "", 1)
			metric.SetDisplayName(displayName)
		} else {
			metric.SetDisplayName(name)
		}
		
		return
	}
	
	// Handle traditional Juniper metrics with standard naming
	switch {
	case strings.Contains(name, "interface"):
		metric.SetCategory("interfaces")
	case strings.Contains(name, "bgp") || strings.Contains(subsystem, "bgp"):
		metric.SetCategory("bgp")
	case strings.Contains(name, "ospf") || strings.Contains(subsystem, "ospf"):
		metric.SetCategory("ospf")
	case strings.Contains(name, "route") || strings.Contains(subsystem, "route"):
		metric.SetCategory("routing")
	case strings.Contains(name, "memory") || strings.Contains(subsystem, "memory"):
		metric.SetCategory("memory")
	case strings.Contains(name, "cpu") || strings.Contains(subsystem, "cpu"):
		metric.SetCategory("cpu")
	case strings.Contains(name, "temperature") || strings.Contains(subsystem, "temp"):
		metric.SetCategory("temperature")
		metric.SetUnit("celsius")
	}
	
	// Set better units for traditional Juniper metrics
	if strings.Contains(name, "bps") || strings.Contains(name, "bitrate") {
		metric.SetUnit("bps")
	} else if strings.Contains(name, "pps") || strings.Contains(name, "packetrate") {
		metric.SetUnit("pps")
	} else if strings.Contains(name, "error") || strings.Contains(name, "discard") {
		if metric.Type() == "counter" {
			metric.SetCategory("errors")
		}
	}
}

// parseCiscoMetric extracts additional information from Cisco metrics
func parseCiscoMetric(metric *metrics.Metric) {
	name := metric.Name()
	subsystem := metric.Subsystem()
	
	// Set more specific categories based on Cisco naming conventions
	switch {
	case strings.Contains(name, "interface"):
		metric.SetCategory("interfaces")
	case strings.Contains(name, "bgp") || strings.Contains(subsystem, "bgp"):
		metric.SetCategory("bgp")
	case strings.Contains(name, "ospf") || strings.Contains(subsystem, "ospf"):
		metric.SetCategory("ospf")
	case strings.Contains(name, "route") || strings.Contains(subsystem, "route"):
		metric.SetCategory("routing")
	case strings.Contains(name, "memory") || strings.Contains(subsystem, "memory"):
		metric.SetCategory("memory")
	case strings.Contains(name, "cpu") || strings.Contains(subsystem, "cpu"):
		metric.SetCategory("cpu")
	}
}