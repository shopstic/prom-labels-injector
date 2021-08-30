package label_injector

import (
	"shopstic/prom-labels-injector/config"
	"strings"
)

func InjectLabels(metrics string, settings *config.LabelInjectorConfig) string {
	var result strings.Builder
	labels := ToLabels(settings)
	lines := strings.Split(metrics, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result.WriteString("\n")
		} else if strings.HasPrefix(trimmed, "#") {
			result.WriteString(trimmed)
			result.WriteString("\n")
		} else {
			result.WriteString(ProcessMetricLine(trimmed, labels))
			result.WriteString("\n")
		}
	}
	return result.String()
}

func ProcessMetricLine(line string, labels string) string {
	index := strings.LastIndex(line, "}")
	var result strings.Builder
	if index == -1 {
		// add labels
		i := 0
		for i < len(line) {
			c := line[i]
			// metric name may contain ASCII letters and digits, as well as underscores and colons
			if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == ':') {
				break
			}
			i += 1
		}
		if i < len(line) {
			result.WriteString(line[:i])
			result.WriteString("{")
			result.WriteString(labels)
			result.WriteString("}")
			result.WriteString(line[i:])
		}
	} else {
		var prevIndex = index
		if line[index-1] == ',' {
			prevIndex -= 1
		}
		result.WriteString(line[:prevIndex])
		result.WriteString(",")
		result.WriteString(labels)
		result.WriteString(line[index:])
	}
	return result.String()
}

func ToLabels(settings *config.LabelInjectorConfig) string {
	var result strings.Builder
	for i, label := range settings.Labels {
		result.WriteString(label.Label)
		result.WriteString("=\"")
		result.WriteString(label.Value)
		result.WriteString("\"")
		if len(settings.Labels)-1 != i {
			result.WriteString(",")
		}
	}
	return result.String()
}
