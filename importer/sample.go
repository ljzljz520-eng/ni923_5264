package importer

import "example.com/toolnav/model"

func SupportedCategories() []model.Category { return model.Categories() }

func ExampleRows() []string {
	return []string{
		"deployctl|DeployCtl|deployment|https://deployctl.example|active|0|release,ssh",
		"pulsewatch|PulseWatch|monitoring|https://pulsewatch.example|beta|1|metrics,alerts",
		"payrail|PayRail|payments|https://payrail.example|active|2|billing,webhooks",
		"docforge|DocForge|documentation|https://docforge.example|active|3|docs,markdown",
		"policykit|PolicyKit|compliance|https://policykit.example|archived|4|audit,policy",
	}
}
