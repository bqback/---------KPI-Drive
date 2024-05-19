package entities

import "time"

type Messages struct {
	Error   *string   `json:"error"`
	Warning *string   `json:"warning"`
	Info    []*string `json:"info"`
}

type Value struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

type TagInfo struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Key           string  `json:"key"`
	Values        []Value `json:"values"`
	ValuesSource  int     `json:"values_source"`
	AuthorUserID  int     `json:"AutorUserID"`
	AuthorUserUID string  `json:"AuthorUserUid"`
}

type Tag struct {
	Tag        TagInfo `json:"tag"`
	Value      string  `json:"value"`
	ValueTitle string  `json:"value_title"`
}

type GetDataRow struct {
	MoFactID       int           `json:"indicator_to_mo_fact_id"`
	MoID           int           `json:"indicator_to_mo_id"`
	UserID         int           `json:"user_id"`
	IsPlan         bool          `json:"is_plan"`
	Value          int           `json:"value"`
	Comment        string        `json:"comment"`
	Author         string        `json:"author"`
	MeasureID      int           `json:"measure_id"`
	MeasureValue   int           `json:"in_measure_value"`
	Complexity     string        `json:"complexity"`
	Mark           int           `json:"mark"`
	FactTime       time.Time     `json:"fact_time"`
	FactExactTime  time.Time     `json:"fact_exact_time"`
	PostTime       time.Time     `json:"post_time"`
	Common         bool          `json:"common"`
	MarkIgnored    bool          `json:"mark_ignored"`
	AllowEdit      bool          `json:"allow_edit"`
	IndicatorSeq   int           `json:"indicator_sequence"`
	TaskStatus     int           `json:"task_status"`
	TaskStatusFmt  string        `json:"task_status_fmt"`
	TaskStatusDesc string        `json:"task_status_desc"`
	Tags           []Tag         `json:"tags"`
	Files          []interface{} `json:"files"`
	SuperTags      []Tag         `json:"supertags"`
	PeriodStart    time.Time     `json:"period_start"`
	PeriodEnd      time.Time     `json:"period_end"`
}

type GetDataEntry struct {
	Page      int          `json:"page"`
	PageCount int          `json:"pages_count"`
	RowCount  int          `json:"rows_count"`
	Rows      []GetDataRow `json:"rows"`
}

type SaveResponse struct {
	Messages Messages       `json:"MESSAGES"`
	Data     map[string]int `json:"DATA"`
	Status   string         `json:"STATUS"`
}

type GetResponse struct {
	Messages Messages `json:"MESSAGES"`
	Status   string   `json:"STATUS"`
}
