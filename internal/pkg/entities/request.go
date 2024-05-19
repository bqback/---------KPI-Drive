package entities

type SaveFact struct {
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

type GetFacts struct {
	PeriodStart string
	PeriodEnd   string
	PeriodKey   string
	MoID        string
}
