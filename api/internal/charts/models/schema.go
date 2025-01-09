package charts_models

type ChartColumn struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Optional    bool   `json:"optional"`
}

var categorical_schema = []ChartColumn{
	{Name: "Category", Type: "string", Description: "Used as the X-Basis for the chart"},
	{Name: "Value", Type: "numerical", Description: "Used to display the amount on the chart"},
}
