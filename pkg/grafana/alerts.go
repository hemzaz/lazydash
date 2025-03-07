package grafana

import (
	"fmt"
	"strings"

	"github.com/hemzaz/lazydash/pkg/metrics"
)

// AlertDefinition defines a Grafana alert configuration
type AlertDefinition struct {
	Name                string            `json:"name"`
	Message             string            `json:"message"`
	Handler             int               `json:"handler"`
	NoDataState         string            `json:"noDataState"`
	ExecutionErrorState string            `json:"executionErrorState"`
	Frequency           string            `json:"frequency"`
	For                 string            `json:"for"`
	Conditions          []AlertCondition  `json:"conditions"`
	Notifications       []AlertNotification `json:"notifications,omitempty"`
}

// AlertCondition defines a single alert condition
type AlertCondition struct {
	Type      string          `json:"type"`
	Target    AlertTarget     `json:"target"`
	Evaluator AlertEvaluator  `json:"evaluator"`
	Reducer   AlertReducer    `json:"reducer"`
	Operator  AlertOperator   `json:"operator,omitempty"`
}

// AlertTarget defines the data source and query for the alert
type AlertTarget struct {
	RefID          string `json:"refId"`
	Datasource     string `json:"datasource"`
	Expr           string `json:"expr,omitempty"`
}

// AlertEvaluator defines how to evaluate the alert
type AlertEvaluator struct {
	Type  string    `json:"type"`
	Params []float64 `json:"params"`
}

// AlertReducer defines how to reduce multiple values
type AlertReducer struct {
	Type     string   `json:"type"`
	Params   []string `json:"params,omitempty"`
}

// AlertOperator defines the operator between conditions
type AlertOperator struct {
	Type string `json:"type"`
}

// AlertNotification defines a notification channel
type AlertNotification struct {
	ID  int `json:"id"`
}

// createAlertDefinition creates a new alert definition
func createAlertDefinition(name, query string, warning, criticalThreshold float64, notify bool) *AlertDefinition {
	// Create basic alert
	alert := &AlertDefinition{
		Name:                fmt.Sprintf("Alert for %s", name),
		Message:             fmt.Sprintf("%s is outside acceptable range", name),
		Handler:             1,
		NoDataState:         "no_data",
		ExecutionErrorState: "alerting",
		Frequency:           "60s", // Check every minute
		For:                 "5m",  // Must be true for 5 minutes
		Conditions: []AlertCondition{
			{
				Type: "query",
				Target: AlertTarget{
					RefID:      "A",
					Datasource: "Prometheus",
					Expr:       query,
				},
				Evaluator: AlertEvaluator{
					Type:   "gt",
					Params: []float64{criticalThreshold},
				},
				Reducer: AlertReducer{
					Type: "avg",
				},
			},
		},
	}
	
	// Add warning condition if specified
	if warning > 0 && warning != criticalThreshold {
		warningCondition := AlertCondition{
			Type: "query",
			Target: AlertTarget{
				RefID:      "B",
				Datasource: "Prometheus",
				Expr:       query,
			},
			Evaluator: AlertEvaluator{
				Type:   "gt",
				Params: []float64{warning},
			},
			Reducer: AlertReducer{
				Type: "avg",
			},
			Operator: AlertOperator{
				Type: "or",
			},
		}
		alert.Conditions = append(alert.Conditions, warningCondition)
	}
	
	// Add notifications if requested
	if notify {
		alert.Notifications = []AlertNotification{
			{
				ID: 1, // Default notification channel
			},
		}
	}
	
	return alert
}

// addAlertToPanel adds an alert definition to a panel
func addAlertToPanel(panel *Panel, alert *AlertDefinition) {
	// Add alert to panel
	panel.Alert = alert
	
	// Make sure panel has the required fields for alerts
	panel.Datasource = "Prometheus"
	
	// Ensure all targets in alert have corresponding panel targets
	for _, condition := range alert.Conditions {
		found := false
		for _, target := range panel.Targets {
			if target.RefID == condition.Target.RefID {
				found = true
				break
			}
		}
		
		if !found {
			// Add target from alert
			panel.Targets = append(panel.Targets, PanelTarget{
				RefID:      condition.Target.RefID,
				Datasource: condition.Target.Datasource,
				Expr:       condition.Target.Expr,
				Format:     "time_series",
			})
		}
	}
}

// generateAlertThreshold creates alert thresholds based on metric type and value patterns
func generateAlertThreshold(metric *metrics.Metric) *AlertThreshold {
	name := metric.Name()
	metricType := metric.Type()
	
	// Different defaults based on metric type
	switch metricType {
	case "counter":
		// For counters, we typically alert on rate of change
		if strings.Contains(name, "error") || strings.Contains(name, "fail") {
			return &AlertThreshold{
				Metric:  name,
				Warning: 0.1,  // 0.1 errors per second
				Error:   1.0,  // 1 error per second
				Notify:  true,
			}
		}
		return nil // No default alerts for other counters
		
	case "gauge":
		// For gauges, set reasonable defaults based on name patterns
		if strings.Contains(name, "cpu") {
			return &AlertThreshold{
				Metric:  name,
				Warning: 80.0,  // 80% CPU
				Error:   95.0,  // 95% CPU
				Notify:  true,
			}
		} else if strings.Contains(name, "memory") || strings.Contains(name, "mem") {
			return &AlertThreshold{
				Metric:  name,
				Warning: 85.0,  // 85% memory
				Error:   95.0,  // 95% memory
				Notify:  true,
			}
		} else if strings.Contains(name, "disk") {
			return &AlertThreshold{
				Metric:  name,
				Warning: 80.0,  // 80% disk usage
				Error:   90.0,  // 90% disk usage
				Notify:  true,
			}
		}
		return nil
		
	case "summary":
		// For summaries (often latency), set thresholds based on common patterns
		if strings.Contains(name, "latency") || strings.Contains(name, "duration") || strings.Contains(name, "time") {
			return &AlertThreshold{
				Metric:  name,
				Warning: 1000.0,  // 1000ms
				Error:   2000.0,  // 2000ms
				Notify:  true,
			}
		}
		return nil
		
	default:
		return nil
	}
}

// AlertThreshold defines threshold levels for generating alerts
type AlertThreshold struct {
	Metric  string  // Which metric to alert on
	Warning float64 // Warning threshold
	Error   float64 // Error threshold
	Notify  bool    // Whether to send notifications
}