// Package grafana provides functionality for creating Grafana dashboards
package grafana

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hemzaz/lazydash/internal/config"
	"github.com/hemzaz/lazydash/pkg/metrics"
	"github.com/hemzaz/lazydash/pkg/query"
	"github.com/rs/zerolog/log"
)

// RefreshIntervals defines available dashboard refresh intervals
var RefreshIntervals = []string{"5s", "10s", "30s", "1m", "5m", "15m", "30m", "1h", "2h", "1d"}

// TimeOptions defines available dashboard time range options
var TimeOptions = []string{"5m", "15m", "1h", "3h", "6h", "12h", "24h", "2d", "3d", "4d", "7d", "30d"}

// TimeRange contains a time range for dashboards
type TimeRange struct {
	From string `json:"from"` // e.g., "now-6h"
	To   string `json:"to"`   // e.g., "now"
}

// TimePicker defines dashboard time picker options
type TimePicker struct {
	RefreshIntervals []string `json:"refresh_intervals,omitempty"`
	TimeOptions      []string `json:"time_options,omitempty"`
}

// Dashboard represents a Grafana dashboard
type Dashboard struct {
	ID            int         `json:"id"`
	UID           string      `json:"uid,omitempty"`
	Title         string      `json:"title"`
	Tags          []string    `json:"tags"`
	Timezone      string      `json:"timezone"`
	Editable      bool        `json:"editable"`
	FolderID      int         `json:"folderId,omitempty"`
	FolderUID     string      `json:"folderUid,omitempty"`
	FolderTitle   string      `json:"folderTitle,omitempty"`
	Description   string      `json:"description,omitempty"`
	HideControls  bool        `json:"hideControls"`
	GraphTooltip  int         `json:"graphTooltip"`
	Panels        []Panel     `json:"panels"`
	Time          TimeRange   `json:"time"`
	TimePicker    TimePicker  `json:"timepicker"`
	Templating    Templating  `json:"templating,omitempty"`
	Annotations   Annotations `json:"annotations,omitempty"`
	SchemaVersion int         `json:"schemaVersion"`
	Version       int         `json:"version"`
}

// Templating represents dashboard template variables
type Templating struct {
	Enable bool            `json:"enable,omitempty"`
	List   []TemplateVar   `json:"list,omitempty"`
}

// TemplateVar represents a single template variable
type TemplateVar struct {
	Name           string             `json:"name,omitempty"`
	Label          string             `json:"label,omitempty"`
	Type           string             `json:"type,omitempty"`
	Datasource     string             `json:"datasource,omitempty"`
	Query          string             `json:"query,omitempty"`
	Refresh        int                `json:"refresh,omitempty"`
	IncludeAll     bool               `json:"includeAll,omitempty"`
	Multi          bool               `json:"multi,omitempty"`
	AllValue       string             `json:"allValue,omitempty"`
	Current        TemplateVarState   `json:"current,omitempty"`
	Options        []TemplateOption   `json:"options,omitempty"`
}

// TemplateVarState represents the current state of a template variable
type TemplateVarState struct {
	Text  string   `json:"text,omitempty"`
	Value string   `json:"value,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// TemplateOption represents an option for a template variable
type TemplateOption struct {
	Selected bool   `json:"selected,omitempty"`
	Text     string `json:"text,omitempty"`
	Value    string `json:"value,omitempty"`
}

// Annotations represents dashboard annotations
type Annotations struct {
	List []Annotation `json:"list,omitempty"`
}

// Annotation represents a single annotation
type Annotation struct {
	Datasource string `json:"datasource,omitempty"`
	Enable     bool   `json:"enable"`
	Name       string `json:"name,omitempty"`
	IconColor  string `json:"iconColor,omitempty"`
	Query      string `json:"query,omitempty"`
	Type       string `json:"type,omitempty"`
}

// NewDashboard creates a new dashboard with default settings
func NewDashboard(title string) *Dashboard {
	return &Dashboard{
		Title:         title,
		Tags:          []string{"prometheus", "generated"},
		Timezone:      "browser",
		Editable:      true,
		HideControls:  false,
		GraphTooltip:  0,
		Panels:        []Panel{},
		Time:          TimeRange{From: "now-6h", To: "now"},
		TimePicker:    TimePicker{RefreshIntervals: RefreshIntervals, TimeOptions: TimeOptions},
		SchemaVersion: 22,
		Version:       0,
	}
}

// SetDescription sets the dashboard description
func (d *Dashboard) SetDescription(desc string) {
	d.Description = desc
}

// AddPanel adds a panel to the dashboard
func (d *Dashboard) AddPanel(panel Panel) {
	// Set the panel ID
	panel.ID = len(d.Panels) + 1
	d.Panels = append(d.Panels, panel)
}

// DumpJSON outputs the dashboard as JSON
func (d *Dashboard) DumpJSON(pretty bool) {
	var b []byte
	var err error

	if pretty {
		b, err = json.MarshalIndent(d, "", "  ")
	} else {
		b, err = json.Marshal(d)
	}

	if err != nil {
		log.Fatal().Err(err).Msg("Could not marshal dashboard to JSON")
	}

	fmt.Println(string(b))
}

// Generate creates a dashboard based on the given metrics
func (d *Dashboard) Generate(metrics *metrics.Registry, cfg *config.Config, queryBuilder *query.Builder) {
	d.Description = cfg.Description

	// Check how to organize the dashboard
	if cfg.VendorConfig != nil && cfg.VendorConfig.Enabled && cfg.VendorConfig.GroupByVendor {
		vendors := metrics.ListVendors()
		if len(vendors) > 0 {
			d.generateWithVendorGroups(metrics, cfg, queryBuilder)
			return
		}
	}
	
	if cfg.AutoCorrelate {
		d.generateWithCorrelation(metrics, cfg, queryBuilder)
	} else if cfg.LabelGrouping != nil && len(cfg.LabelGrouping.GroupByLabels) > 0 {
		d.generateWithLabelGroups(metrics, cfg, queryBuilder)
	} else {
		d.generateStandard(metrics, cfg, queryBuilder)
	}
}

// generateStandard creates a standard dashboard with all metrics
func (d *Dashboard) generateStandard(registry *metrics.Registry, cfg *config.Config, queryBuilder *query.Builder) {
	pgrid := NewPanelGridLayout()

	// Process each metric
	registry.ForEach(func(name string, metric *metrics.Metric) {
		// Create panel for this metric
		panel := createPanelForMetric(metric, cfg, queryBuilder)
		
		// Set grid position
		panel.SetGridPos(pgrid.X, pgrid.Y, 8, 12)
		
		// Add panel to dashboard
		d.AddPanel(*panel)
		
		// Update grid position for next panel
		if (pgrid.Count % 2) < 1 {
			pgrid.UpdateX(0)
			pgrid.UpdateY(9)
		} else {
			pgrid.UpdateX(12)
		}
	})
}

// generateWithLabelGroups creates a dashboard with panels grouped by labels
func (d *Dashboard) generateWithLabelGroups(registry *metrics.Registry, cfg *config.Config, queryBuilder *query.Builder) {
	yPos := 0
	
	// Create a map to group metrics by label values
	groups := make(map[string][]*metrics.Metric)
	
	// Group metrics by their label values
	registry.ForEach(func(name string, metric *metrics.Metric) {
		groupKey := ""
		for _, labelName := range cfg.LabelGrouping.GroupByLabels {
			if metric.HasLabel(labelName) {
				if groupKey != "" {
					groupKey += ":"
				}
				groupKey += labelName
			}
		}
		
		if groupKey == "" {
			groupKey = "other"
		}
		
		groups[groupKey] = append(groups[groupKey], metric)
	})
	
	// Process each group
	for groupName, groupMetrics := range groups {
		// Add a row header if enabled
		if cfg.LabelGrouping.SeparateRows {
			rowPanel := &Panel{
				Title:       strings.Replace(groupName, ":", " ", -1),
				Type:        "row",
				Description: "Metrics grouped by " + groupName,
				GridPos: PanelGridPos{
					X: 0,
					Y: yPos,
					W: 24,
					H: 1,
				},
			}
			d.AddPanel(*rowPanel)
			yPos += 1
		}
		
		// Process metrics in this group
		xPos := 0
		for i, metric := range groupMetrics {
			// Create panel for this metric
			panel := createPanelForMetric(metric, cfg, queryBuilder)
			
			// Calculate panel width
			panelWidth := 12
			if cfg.LabelGrouping.PanelsPerRow > 0 {
				panelWidth = 24 / cfg.LabelGrouping.PanelsPerRow
			}
			
			// Set grid position
			panel.SetGridPos(xPos, yPos, 8, panelWidth)
			
			// Add panel to dashboard
			d.AddPanel(*panel)
			
			// Update position for next panel
			xPos += panelWidth
			
			// Check if we need to start a new row
			if cfg.LabelGrouping.PanelsPerRow > 0 && (i+1) % cfg.LabelGrouping.PanelsPerRow == 0 {
				xPos = 0
				yPos += 8
			}
		}
		
		// Ensure we're at the start of a row for the next group
		if xPos > 0 {
			yPos += 8
		} else if len(groupMetrics) > 0 {
			yPos += 1
		}
	}
}

// generateWithCorrelation creates a dashboard with panels grouped by auto-correlation
func (d *Dashboard) generateWithCorrelation(registry *metrics.Registry, cfg *config.Config, queryBuilder *query.Builder) {
	// This would implement auto-correlation using the correlation.go logic
	// For now, just use standard generation
	d.generateStandard(registry, cfg, queryBuilder)
}

// generateWithVendorGroups creates a dashboard with panels grouped by vendor
func (d *Dashboard) generateWithVendorGroups(registry *metrics.Registry, cfg *config.Config, queryBuilder *query.Builder) {
	yPos := 0
	vendors := registry.ListVendors()
	
	// Add section for metrics with no vendor identified
	if nonVendorMetrics := registry.Filter(func(name string, metric *metrics.Metric) bool {
		return metric.Vendor() == ""
	}); nonVendorMetrics.Count() > 0 {
		// Add a row header for general metrics
		rowPanel := &Panel{
			Title:       "General Metrics",
			Type:        "row",
			Description: "General metrics with no vendor prefix",
			GridPos: PanelGridPos{
				X: 0,
				Y: yPos,
				W: 24,
				H: 1,
			},
		}
		d.AddPanel(*rowPanel)
		yPos += 1
		
		// Add metrics to this row
		xPos := 0
		metricCount := 0
		nonVendorMetrics.ForEach(func(name string, metric *metrics.Metric) {
			// Create panel for this metric
			panel := createPanelForMetric(metric, cfg, queryBuilder)
			
			// Calculate panel width
			panelWidth := 12
			if cfg.LabelGrouping != nil && cfg.LabelGrouping.PanelsPerRow > 0 {
				panelWidth = 24 / cfg.LabelGrouping.PanelsPerRow
			}
			
			// Set grid position
			panel.SetGridPos(xPos, yPos, 8, panelWidth)
			
			// Add panel to dashboard
			d.AddPanel(*panel)
			
			// Update position for next panel
			xPos += panelWidth
			metricCount++
			
			// Check if we need to start a new row
			if cfg.LabelGrouping != nil && cfg.LabelGrouping.PanelsPerRow > 0 && metricCount % cfg.LabelGrouping.PanelsPerRow == 0 {
				xPos = 0
				yPos += 8
			}
		})
		
		// Ensure we're at the start of a row for the next vendor
		if xPos > 0 {
			yPos += 8
		}
	}
	
	// Now add sections for each vendor
	for _, vendor := range vendors {
		vendorMetrics := registry.FilterByVendor(vendor)
		if vendorMetrics.Count() == 0 {
			continue
		}
		
		// Add a row header for this vendor
		vendorTitle := vendor
		switch vendor {
		case "juniper":
			vendorTitle = "Juniper Networks"
		case "cisco":
			vendorTitle = "Cisco Systems"
		case "arista":
			vendorTitle = "Arista Networks"
		case "huawei":
			vendorTitle = "Huawei Technologies"
		case "paloalto":
			vendorTitle = "Palo Alto Networks"
		case "fortinet":
			vendorTitle = "Fortinet"
		case "f5":
			vendorTitle = "F5 Networks"
		case "checkpoint":
			vendorTitle = "Check Point"
		}
		
		rowPanel := &Panel{
			Title:       vendorTitle,
			Type:        "row",
			Description: "Metrics for " + vendorTitle,
			GridPos: PanelGridPos{
				X: 0,
				Y: yPos,
				W: 24,
				H: 1,
			},
		}
		d.AddPanel(*rowPanel)
		yPos += 1
		
		// If there's category grouping available, organize by category
		categories := make(map[string][]*metrics.Metric)
		vendorMetrics.ForEach(func(name string, metric *metrics.Metric) {
			category := metric.Category()
			if category == "" {
				category = "other"
			}
			categories[category] = append(categories[category], metric)
		})
		
		// If we have categories, organize by them
		if len(categories) > 1 {
			// Process each category
			for category, categoryMetrics := range categories {
				if len(categoryMetrics) == 0 {
					continue
				}
				
				// Add a category row
				categoryTitle := strings.ToUpper(category[:1]) + category[1:]
				categoryPanel := &Panel{
					Title:       categoryTitle,
					Type:        "row",
					Description: categoryTitle + " metrics for " + vendorTitle,
					GridPos: PanelGridPos{
						X: 0,
						Y: yPos,
						W: 24,
						H: 1,
					},
				}
				d.AddPanel(*categoryPanel)
				yPos += 1
				
				// Add metrics in this category
				xPos := 0
				metricCount := 0
				for _, metric := range categoryMetrics {
					// Create panel for this metric
					panel := createPanelForMetric(metric, cfg, queryBuilder)
					
					// Calculate panel width
					panelWidth := 12
					if cfg.LabelGrouping != nil && cfg.LabelGrouping.PanelsPerRow > 0 {
						panelWidth = 24 / cfg.LabelGrouping.PanelsPerRow
					}
					
					// Set grid position
					panel.SetGridPos(xPos, yPos, 8, panelWidth)
					
					// Add panel to dashboard
					d.AddPanel(*panel)
					
					// Update position for next panel
					xPos += panelWidth
					metricCount++
					
					// Check if we need to start a new row
					if cfg.LabelGrouping != nil && cfg.LabelGrouping.PanelsPerRow > 0 && metricCount % cfg.LabelGrouping.PanelsPerRow == 0 {
						xPos = 0
						yPos += 8
					}
				}
				
				// Ensure we're at the start of a row for the next category
				if xPos > 0 {
					yPos += 8
				}
			}
		} else {
			// No categories, just show all vendor metrics flat
			xPos := 0
			metricCount := 0
			vendorMetrics.ForEach(func(name string, metric *metrics.Metric) {
				// Create panel for this metric
				panel := createPanelForMetric(metric, cfg, queryBuilder)
				
				// Calculate panel width
				panelWidth := 12
				if cfg.LabelGrouping != nil && cfg.LabelGrouping.PanelsPerRow > 0 {
					panelWidth = 24 / cfg.LabelGrouping.PanelsPerRow
				}
				
				// Set grid position
				panel.SetGridPos(xPos, yPos, 8, panelWidth)
				
				// Add panel to dashboard
				d.AddPanel(*panel)
				
				// Update position for next panel
				xPos += panelWidth
				metricCount++
				
				// Check if we need to start a new row
				if cfg.LabelGrouping != nil && cfg.LabelGrouping.PanelsPerRow > 0 && metricCount % cfg.LabelGrouping.PanelsPerRow == 0 {
					xPos = 0
					yPos += 8
				}
			})
			
			// Ensure we're at the start of a row for the next vendor
			if xPos > 0 {
				yPos += 8
			}
		}
	}
}

// Helper function to create a panel for a metric
func createPanelForMetric(metric *metrics.Metric, cfg *config.Config, queryBuilder *query.Builder) *Panel {
	// Use display name if available, otherwise format the metric name
	var title string
	if metric.DisplayName() != metric.Name() {
		// Use custom display name if it was set
		title = strings.Replace(metric.DisplayName(), "_", " ", -1)
	} else {
		title = strings.Replace(metric.Name(), "_", " ", -1)
	}
	
	// Create base panel
	panel := NewPanel(title)
	panel.SetDescription(metric.Help())
	panel.SetUnit(metric.Unit())
	
	// Set up legend
	if cfg.Table {
		panel.Legend = PanelLegend{
			Show:         true,
			Current:      true,
			Values:       true,
			AlignAsTable: true,
		}
	}
	
	// Build query
	expr := queryBuilder.BuildQuery(metric)
	panel.SetMetricExpr(expr)
	
	// Set legend format
	panel.SetLegendFormat(queryBuilder.GetLegend(metric))
	
	// Determine and set visualization type
	if cfg.Visualizations != nil {
		visualType := getVisualizationType(metric, cfg)
		configureVisualization(panel, metric, visualType)
	} else if cfg.Gauges && metric.Type() == "gauge" {
		panel.SetType("gauge")
	} else {
		panel.SetType("graph")
	}
	
	// Add alerts if enabled and applicable
	if cfg.GenerateAlerts {
		alertThreshold := generateAlertThreshold(metric)
		if alertThreshold != nil {
			alertDef := createAlertDefinition(
				metric.Name(),
				expr,
				alertThreshold.Warning,
				alertThreshold.Error,
				alertThreshold.Notify,
			)
			addAlertToPanel(panel, alertDef)
		}
	}
	
	return panel
}