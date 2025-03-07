# lazydash [![CircleCI](https://circleci.com/gh/hemzaz/lazydash/tree/master.svg?style=svg&circle-token=166cef586b42bb07d2e81ffaffaac8bd371970d2)](https://circleci.com/gh/hemzaz/lazydash/tree/master)

Auto generate Grafana dashboards based on prometheus metrics endpoints. Allows for quick prototyping of dashboards. 

# Notes

Version: 1.0.0

* Supports Counter, Gauge, Summary, and Histogram metrics
* Advanced panel organization with label grouping and auto-correlation
* Multiple visualization types (graphs, gauges, stats, tables, heatmaps)
* Automatic alert generation
* Folder organization in Grafana
* Assumes metrics adhere to prometheus metrics naming conventions and standards e.g:
```
# HELP builder_builds_triggered_total Number of triggered image builds
# TYPE builder_builds_triggered_total counter
builder_builds_triggered_total 0"
```
# Features

* Generate dashboards directly from any Prometheus metrics endpoint
* Post generated dashboards directly to Grafana via the API
* Intelligent panel organization by label or correlation
* Advanced visualization selection based on metric type
* Auto panel legend format with all labels from prometheus
* Panels use the metrics HELP field as the description
* Override default query expressions and legend formats
* Folder organization for dashboard management
* Alert generation based on metric patterns

<img src="lazydash.gif" alt="demo" width="800"/>
---

# Usage

```
usage: lazydash [<flags>]

Generate a Grafana dashboard from Prometheus metrics data via file, stdin or HTTP url

Flags:
  -h, --help             Show context-sensitive help (also try --help-long and --help-man).
      --version          Show application version.
  -f, --file=""          Parse metrics from file.
  -t, --title="Prometheus Dashboard"  
                         Dashboard title
      --description="Generated by Lazydash"  
                         Dashboard description
      --stdin            Read from stdin
      --url=""           Fetch Prometheus data from HTTP(S) url
  -p, --pretty           Print pretty indented JSON
  -g, --gauges           Render gauge values as gauge panel type instead of graph
      --table            Render legend as a table
      --set-counter-expr="sum(rate(:METRIC: [1m]))"  
                         Set custom meterics query expression for counter type metric
      --set-gauge-expr=":METRIC:"  
                         Set custom meterics query expression for gauge type metric
      --set-summary-expr=":METRIC:"  
                         Set custom meterics query expression for summary type metric
      --set-delimiter=":METRIC:"  
                         Set custom meterics delimiter used to insert metric name into expression
      --set-counter-legend="Job:[{{job}}]"  
                         Set the default counter panel legend format
      --set-gauge-legend="Job:[{{job}}]"  
                         Set the default counter panel legend format
  -H, --grafana-url=""   Set the grafana api url e.g http://grafana.example.com:3000
  -I, --insecure         Skip ssl certificate verification
  -T, --token=""         Set the grafana api token
      --folder=""        Set the Grafana folder name
      --folder-create    Create the folder if it doesn't exist
      --folder-description="Generated by Lazydash"  
                         Set the Grafana folder description
      --group-by= ...    Group panels by label
      --separate-rows    Create separate rows for each label group
      --panels-per-row=2 Number of panels per row (0 for auto)
      --auto-correlate   Automatically correlate related metrics
      --correlation-threshold=0.7  
                         Threshold for auto-correlation (0.0-1.0)
      --viz-counters=graph  
                         Visualization type for counters
      --viz-gauges=gauge  
                         Visualization type for gauges
      --viz-summaries=graph  
                         Visualization type for summaries
      --heatmap-for-histograms  
                         Use heatmap for histogram metrics
      --stat-for-gauges  Use stat panels for simple gauges
      --table-for-multilabels  
                         Use tables for metrics with many labels
      --generate-alerts  Automatically generate alerts for common metrics
```
---

# Advanced Features

## Organization
* **Label Grouping**: Group panels by metric labels `--group-by="job" --group-by="instance"`
* **Folder Organization**: Create and use Grafana folders `--folder="My Dashboard" --folder-description="Generated dashboards"`
* **Auto-correlation**: Group related metrics together `--auto-correlate --correlation-threshold=0.8`
* **Vendor Grouping**: Automatically detect and group vendor-specific metrics:
  * `--vendor-detect` - Enable vendor-specific metric detection
  * `--group-by-vendor` - Group metrics by vendor in the dashboard
  * `--vendor-prefix="myvendor_"` - Add custom vendor prefixes to detect
  * `--juniper` - Enable special handling for Juniper metrics
  * `--cisco` - Enable special handling for Cisco metrics

## Visualization
* **Multiple Visualization Types**: Choose from graphs, gauges, stats, tables, heatmaps
  * `--viz-counters=graph` 
  * `--viz-gauges=gauge`
  * `--viz-summaries=graph`
* **Automatic Visualization Selection**:
  * `--heatmap-for-histograms` - Use heatmaps for histogram metrics
  * `--stat-for-gauges` - Use stat panels for simple gauges
  * `--table-for-multilabels` - Use tables for metrics with many labels

## Alerting
* **Auto-generated Alerts**: `--generate-alerts` - Creates alert rules based on metric patterns

# Examples

## Pull metric types from prometheus HTTP endpoint and post to the grafana API
```
lazydash -p --url="http://localhost:9323/metrics" -H http://localhost:3000 -T "eyJrIjoiRzRTUGV1a2RWcjgzbklvVzdXenIySEhJWEJlSkx4UksiLCJuIjoidGVzdCIsImlkIjoxfQ==" --table
```

## Advanced dashboard with grouping and visualization selection
```
lazydash -p --url=http://localhost:9323/metrics --group-by="job" --auto-correlate --heatmap-for-histograms --stat-for-gauges --generate-alerts > dashboard.json
```

## Detect and group vendor-specific metrics (e.g., Juniper, Cisco)
```
lazydash -p --url=http://prometheus:9090/api/v1/query?query=up --vendor-detect --juniper --cisco --group-by-vendor --generate-alerts > network_devices.json
```

## Process Juniper JTIMON metrics from a file
```
lazydash -f JuniperMetrics.txt --vendor-detect --juniper --group-by-vendor --pretty > juniper_dashboard.json
```

## Override the counter expression and legend format
```
lazydash -p --url=http://localhost:9323/metrics --set-counter-expr="sum(rate(:METRIC: [1m])) by(instance)" --set-counter-legend="Instance:[{{instance}}]" > dash1.json
```

## Pipe via curl
```
curl -s http://localhost:9323/metrics | lazydash -t "Demo" -p
```

## Using file input
```
lazydash -f promdata.txt -p
```

---

# Build and Install 

Requirement: Go 1.21+

```
# Install from source
git clone https://github.com/hemzaz/lazydash.git
cd lazydash
go install ./cmd/lazydash

# Or use go install directly
go install github.com/hemzaz/lazydash/cmd/lazydash@latest
```

