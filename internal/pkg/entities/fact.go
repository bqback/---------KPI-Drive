package entities

type Fact struct {
	PeriodStart string
	PeriodEnd   string
	PeriodKey   string
	MoID        string
	MoFactID    string
	Value       string
	FactTime    string
	IsPlan      string
	AuthUserID  string
	Comment     string
}

type GeneratorPreset struct {
	Facts       int    `yaml:"facts"`
	PeriodStart string `yaml:"period_start"`
	PeriodEnd   string `yaml:"period_end"`
	PeriodKey   string `yaml:"period_key"`
	MoID        string `yaml:"mo_id"`
	MoFactID    string `yaml:"mo_fact_id"`
	IsPlan      string `yaml:"is_plan"`
	AuthUserID  string `yaml:"-"`
}
