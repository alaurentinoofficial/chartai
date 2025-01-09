package charts_models

import (
	"fmt"
)

func StringToChart(value string) (ChartEntity, error) {
	if value, ok := registradedCharts[value]; ok {
		return value, nil
	}

	return nil, fmt.Errorf("unsupported chart type: %s", value)
}
