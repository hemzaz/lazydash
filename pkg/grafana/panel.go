package grafana

import (
	"github.com/hemzaz/lazydash/internal/config"
	"github.com/hemzaz/lazydash/pkg/metrics"
)

// Panel represents a Grafana dashboard panel
type Panel struct {
	GridPos       PanelGridPos  `json:"gridPos"`
	Type          string        `json:"type"`
	Title         string        `json:"title"`
	ID            int           `json:"id"`
	Mode          string        `json:"mode,omitempty"`
	Content       string        `json:"content,omitempty"`
	Targets       []PanelTarget `json:"targets,omitempty"`
	Description   string        `json:"description,omitempty"`
	Legend        PanelLegend   `json:"legend,omitempty"`
	Bars          bool          `json:"bars,omitempty"`
	DashLength    int           `json:"dashLength,omitempty"`
	Dashes        bool          `json:"dashes,omitempty"`
	Fill          int           `json:"fill,omitempty"`
	FillGradient  int           `json:"fillGradient,omitempty"`
	Lines         bool          `json:"lines,omitempty"`
	LinesWidth    int           `json:"linewidth,omitempty"`
	NullPointMode string        `json:"nullPointMode,omitempty"`
	Percentage    bool          `json:"percentage,omitempty"`
	PointRadius   int           `json:"pointradius,omitempty"`
	Points        bool          `json:"points,omitempty"`
	Render        string        `json:"renderer,omitempty"`
	SpaceLength   int           `json:"spaceLength,omitempty"`
	Stack         bool          `json:"stack,omitempty"`
	SteppedLine   bool          `json:"steppedLine,omitempty"`
	YAxes         []PanelYAxes  `json:"yaxes,omitempty"`
	YAxis         PanelYAxis    `json:"yaxis,omitempty"`
	XAxis         PanelXAxis    `json:"xaxis,omitempty"`
	ToolTip       PanelToolTip  `json:"tooltip,omitempty"`
	Options       PanelOptions  `json:"options,omitempty"`
	Datasource    string        `json:"datasource,omitempty"`
	Alert         *AlertDefinition `json:"alert,omitempty"`
	
	// Additional fields for new visualization types
	DataFormat      string        `json:"dataFormat,omitempty"`
	HideZeroBuckets bool          `json:"hideZeroBuckets,omitempty"`
	HighlightCards  bool          `json:"highlightCards,omitempty"`
	Color           HeatmapColor  `json:"color,omitempty"`
	Columns         []TableColumn `json:"columns,omitempty"`
	Transform       string        `json:"transform,omitempty"`
	Sort            TableSort     `json:"sort,omitempty"`
}

// PanelOptions contains options for different panel types
type PanelOptions struct {
	ShowThresholdMarkers bool              `json:"showThresholdMarkers,omitempty"`
	ShowThresholdLabels  bool              `json:"showThresholdLabels,omitempty"`
	FieldOptions         PanelFieldOptions `json:"fieldOptions,omitempty"`
	Orientation          string            `json:"orientation,omitempty"`
	ColorMode            string            `json:"colorMode,omitempty"`
	GraphMode            string            `json:"graphMode,omitempty"`
	TextMode             string            `json:"textMode,omitempty"`
	DisplayMode          string            `json:"displayMode,omitempty"`
}

// PanelFieldOptions contains field configuration for panels
type PanelFieldOptions struct {
	Values   bool                      `json:"values,omitempty"`
	Calcs    []string                  `json:"calcs,omitempty"`
	Defaults PanelFieldOptionsDefaults `json:"defaults,omitempty"`
}

// PanelFieldOptionsDefaults contains default settings for field options
type PanelFieldOptionsDefaults struct {
	Thresholds []PanelFieldOptionsThreshold `json:"thresholds,omitempty"`
	Unit       string                       `json:"unit,omitempty"`
}

// PanelFieldOptionsThreshold contains threshold settings
type PanelFieldOptionsThreshold struct {
	Value int    `json:"value,omitempty"`
	Color string `json:"color,omitempty"`
}

// PanelGridPos represents the position of a panel
type PanelGridPos struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
	H int `json:"h,omitempty"`
	W int `json:"w,omitempty"`
}

// PanelTarget represents a query target for a panel
type PanelTarget struct {
	Expr         string `json:"expr,omitempty"`
	RefID        string `json:"refId,omitempty"`
	LegendFormat string `json:"legendFormat,omitempty"`
	Format       string `json:"format,omitempty"`
	Datasource   string `json:"datasource,omitempty"`
}

// PanelLegend contains legend configuration
type PanelLegend struct {
	Avg          bool `json:"avg,omitempty"`
	Current      bool `json:"current,omitempty"`
	Max          bool `json:"max,omitempty"`
	Min          bool `json:"min,omitempty"`
	Show         bool `json:"show,omitempty"`
	Total        bool `json:"total,omitempty"`
	Values       bool `json:"values,omitempty"`
	HideEmpty    bool `json:"hideEmpty,omitempty"`
	HideZero     bool `json:"hideZero,omitempty"`
	RightSide    bool `json:"rightSide,omitempty"`
	SideWidth    int  `json:"sideWidth,omitempty"`
	AlignAsTable bool `json:"alignAsTable,omitempty"`
}

// PanelToolTip contains tooltip configuration
type PanelToolTip struct {
	Shared    bool   `json:"shared,omitempty"`
	Sort      int    `json:"sort,omitempty"`
	ValueType string `json:"value_type,omitempty"`
}

// PanelXAxis contains x-axis configuration
type PanelXAxis struct {
	Buckets int    `json:"buckets,omitempty"`
	Mode    string `json:"mode,omitempty"`
	Name    string `json:"name,omitempty"`
	Show    bool   `json:"show,omitempty"`
}

// PanelYAxes contains y-axis configuration
type PanelYAxes struct {
	Decimals int    `json:"decimals,omitempty"`
	Format   string `json:"format,omitempty"`
	Label    string `json:"label,omitempty"`
	LogBase  int    `json:"logBase,omitempty"`
	Max      int    `json:"max,omitempty"`
	Min      int    `json:"min,omitempty"`
	Show     bool   `json:"show,omitempty"`
}

// PanelYAxis contains additional y-axis options
type PanelYAxis struct {
	Align      bool `json:"align,omitempty"`
	AlignLevel int  `json:"alignLevel,omitempty"`
}

// HeatmapColor defines heatmap color settings
type HeatmapColor struct {
	Mode       string  `json:"mode,omitempty"`
	CardColor  string  `json:"cardColor,omitempty"`
	ColorScale string  `json:"colorScale,omitempty"`
	Exponent   float64 `json:"exponent,omitempty"`
	Min        int     `json:"min,omitempty"`
}

// TableColumn defines a column in a table panel
type TableColumn struct {
	Text  string `json:"text,omitempty"`
	Value string `json:"value,omitempty"`
}

// TableSort defines sorting for a table panel
type TableSort struct {
	Col  int  `json:"col,omitempty"`
	Desc bool `json:"desc,omitempty"`
}

// PanelGridLayout represents the layout state when positioning panels
type PanelGridLayout struct {
	X     int
	Y     int
	Count int
}

// NewPanelGridLayout creates a new panel grid layout
func NewPanelGridLayout() *PanelGridLayout {
	return &PanelGridLayout{X: 0, Y: 0, Count: 0}
}

// UpdateY updates the Y position
func (p *PanelGridLayout) UpdateY(y int) {
	p.Y = p.Y + y
	p.Count = p.Count + 1
}

// UpdateX updates the X position
func (p *PanelGridLayout) UpdateX(x int) {
	if x > 0 {
		p.X = p.X + x
		p.Count = p.Count + 1
	} else {
		p.X = 0
	}
}

// NewPanel creates a new panel with basic settings
func NewPanel(title string) *Panel {
	return &Panel{
		Title:       title,
		Type:        "graph",
		Description: "",
		Targets: []PanelTarget{
			{
				Expr:         "",
				RefID:        "A",
				LegendFormat: "",
				Format:       "time_series",
			},
		},
		YAxes: []PanelYAxes{
			{
				Format:  "short",
				LogBase: 1,
				Show:    true,
			},
			{
				Show: false,
			},
		},
		XAxis: PanelXAxis{
			Mode: "time",
			Show: true,
		},
		Legend: PanelLegend{
			Show:    true,
			Current: false,
			Values:  false,
		},
		Options: PanelOptions{
			FieldOptions: PanelFieldOptions{
				Defaults: PanelFieldOptionsDefaults{
					Unit: "short",
				},
			},
		},
		Bars:       false,
		Lines:      true,
		LinesWidth: 1,
		Fill:       1,
	}
}

// SetGridPos sets the panel's grid position
func (p *Panel) SetGridPos(x, y, h, w int) {
	p.GridPos.X = x
	p.GridPos.Y = y
	p.GridPos.H = h
	p.GridPos.W = w
}

// SetType sets the panel's visualization type
func (p *Panel) SetType(ptype string) {
	p.Type = ptype
}

// SetUnit sets the panel's unit
func (p *Panel) SetUnit(unit string) {
	// Set unit for graph panels
	if len(p.YAxes) > 0 {
		p.YAxes[0].Format = unit
	}
	
	// Set unit for other panel types
	p.Options.FieldOptions.Defaults.Unit = unit
}

// SetLegendFormat sets the legend format string
func (p *Panel) SetLegendFormat(format string) {
	for i := range p.Targets {
		p.Targets[i].LegendFormat = format
	}
}

// SetDescription sets the panel description
func (p *Panel) SetDescription(description string) {
	p.Description = description
}

// SetMetricExpr sets the panel's query expression
func (p *Panel) SetMetricExpr(expr string) {
	for i := range p.Targets {
		p.Targets[i].Expr = expr
	}
}

// getVisualizationType determines the best visualization for a metric
func getVisualizationType(metric *metrics.Metric, cfg *config.Config) config.VisualizationType {
	// If there's a specific override in the config, use that
	if cfg.Visualizations != nil {
		switch metric.Type() {
		case "counter":
			if cfg.Visualizations.CounterType != "" {
				return cfg.Visualizations.CounterType
			}
		case "gauge":
			if cfg.Visualizations.GaugeType != "" {
				return cfg.Visualizations.GaugeType
			}
		case "summary":
			if cfg.Visualizations.SummaryType != "" {
				return cfg.Visualizations.SummaryType
			}
		}
		
		// If there's a default type specified, use that
		if cfg.Visualizations.DefaultType != "" {
			return cfg.Visualizations.DefaultType
		}
	}
	
	// Otherwise, determine based on metric properties
	switch metric.Type() {
	case "counter":
		return config.VisualizationGraph
		
	case "gauge":
		// If config setting for gauge panels
		if cfg.Gauges {
			return config.VisualizationGauge
		}
		
		// If we should use stat panel for gauges and it has few labels
		if cfg.Visualizations != nil && cfg.Visualizations.UseStatForGauges && metric.LabelCount() <= 1 {
			return config.VisualizationStat
		}
		
		return config.VisualizationGraph
		
	case "summary":
		// If it's a histogram and we want to use heatmaps
		if metric.Suffix() == "_bucket" && cfg.Visualizations != nil && cfg.Visualizations.UseHeatmapForHistograms {
			return config.VisualizationHeatmap
		}
		
		return config.VisualizationGraph
		
	default:
		return config.VisualizationGraph
	}
}

// configureVisualization sets up a panel with the appropriate visualization
func configureVisualization(panel *Panel, metric *metrics.Metric, visualType config.VisualizationType) {
	panel.SetType(string(visualType))
	
	// Apply specific configuration based on visualization type
	switch visualType {
	case config.VisualizationHeatmap:
		configureHeatmapPanel(panel, metric)
		
	case config.VisualizationTable:
		configureTablePanel(panel, metric)
		
	case config.VisualizationStat:
		configureStatPanel(panel, metric)
		
	case config.VisualizationBarGauge:
		configureBarGaugePanel(panel, metric)
		
	case config.VisualizationGauge:
		// Gauge panel specific settings
		panel.Options.FieldOptions.Values = false
		panel.Options.FieldOptions.Calcs = []string{"lastNotNull"}
		panel.Options.ShowThresholdLabels = true
		panel.Options.ShowThresholdMarkers = true
		
	case config.VisualizationGraph:
		// Graph is the default, but we can add specific graph settings
		panel.Bars = false
		panel.Lines = true
		panel.LinesWidth = 1
		panel.Fill = 1
		panel.Dashes = false
		panel.Points = false
	}
}

// configureHeatmapPanel configures a panel as a heatmap
func configureHeatmapPanel(panel *Panel, metric *metrics.Metric) {
	panel.Options.FieldOptions.Values = true
	panel.Options.FieldOptions.Calcs = []string{"mean"}
	
	// Set appropriate X and Y axis settings
	panel.XAxis.Mode = "histogram"
	if len(panel.YAxes) > 0 {
		panel.YAxes[0].Format = metric.Unit()
		panel.YAxes[0].LogBase = 1
	}
	
	// Add heatmap-specific options
	panel.DataFormat = "tsbuckets"
	panel.HideZeroBuckets = true
	panel.HighlightCards = true
	
	// Set an appropriate color scheme
	panel.Color = HeatmapColor{
		Mode:       "spectrum",
		CardColor:  "#b4ff00",
		ColorScale: "sqrt",
		Exponent:   0.5,
		Min:        0,
	}
}

// configureTablePanel configures a panel as a table
func configureTablePanel(panel *Panel, metric *metrics.Metric) {
	panel.Options.FieldOptions.Values = true
	panel.Options.FieldOptions.Calcs = []string{"lastNotNull"}
	
	// Set appropriate column settings
	panel.Columns = []TableColumn{
		{
			Text:  "Time",
			Value: "time",
		},
		{
			Text:  "Value",
			Value: "value",
		},
	}
	
	// Add all labels as columns
	for _, label := range metric.Labels() {
		panel.Columns = append(panel.Columns, TableColumn{
			Text:  label,
			Value: "label_" + label,
		})
	}
	
	// Set table display options
	panel.Transform = "timeseries_to_columns"
	panel.Sort = TableSort{
		Col:  0,
		Desc: true,
	}
}

// configureStatPanel configures a panel as a stat panel
func configureStatPanel(panel *Panel, metric *metrics.Metric) {
	panel.Options.FieldOptions.Values = false
	panel.Options.FieldOptions.Calcs = []string{"lastNotNull"}
	panel.Options.FieldOptions.Defaults.Unit = metric.Unit()
	
	// Set thresholds based on metric type
	switch metric.Type() {
	case "counter":
		panel.Options.FieldOptions.Defaults.Thresholds = []PanelFieldOptionsThreshold{
			{Value: 0, Color: "green"},
			{Value: 100, Color: "orange"},
			{Value: 500, Color: "red"},
		}
	case "gauge":
		if metric.Unit() == "percent" {
			panel.Options.FieldOptions.Defaults.Thresholds = []PanelFieldOptionsThreshold{
				{Value: 0, Color: "green"},
				{Value: 80, Color: "orange"},
				{Value: 90, Color: "red"},
			}
		} else {
			panel.Options.FieldOptions.Defaults.Thresholds = []PanelFieldOptionsThreshold{
				{Value: 0, Color: "green"},
				{Value: 70, Color: "orange"},
				{Value: 90, Color: "red"},
			}
		}
	}
	
	// Set other stat panel options
	panel.Options.ColorMode = "value"
	panel.Options.GraphMode = "area"
	panel.Options.TextMode = "auto"
}

// configureBarGaugePanel configures a panel as a bar gauge
func configureBarGaugePanel(panel *Panel, metric *metrics.Metric) {
	panel.Options.FieldOptions.Values = false
	panel.Options.FieldOptions.Calcs = []string{"lastNotNull"}
	panel.Options.FieldOptions.Defaults.Unit = metric.Unit()
	
	// Set thresholds based on metric type (similar to stat panel)
	switch metric.Type() {
	case "counter":
		panel.Options.FieldOptions.Defaults.Thresholds = []PanelFieldOptionsThreshold{
			{Value: 0, Color: "green"},
			{Value: 100, Color: "orange"},
			{Value: 500, Color: "red"},
		}
	case "gauge":
		if metric.Unit() == "percent" {
			panel.Options.FieldOptions.Defaults.Thresholds = []PanelFieldOptionsThreshold{
				{Value: 0, Color: "green"},
				{Value: 80, Color: "orange"},
				{Value: 90, Color: "red"},
			}
		} else {
			panel.Options.FieldOptions.Defaults.Thresholds = []PanelFieldOptionsThreshold{
				{Value: 0, Color: "green"},
				{Value: 70, Color: "orange"},
				{Value: 90, Color: "red"},
			}
		}
	}
	
	// Set orientation
	panel.Options.Orientation = "horizontal"
	panel.Options.DisplayMode = "gradient"
}