package main

import vegeta "github.com/tsenart/vegeta/v12/lib"

type Stage struct {
	Name    string           `json:"name"`
	Target  string           `json:"target"`
	Profile map[string]int64 `json:"profile"`
	Steps   []Step           `json:"steps"`
}

type Step struct {
	Duration uint64 `json:"duration"`
	Rate     int    `json:"rate"`
}

type StageResult struct {
	Name         string         `json:"name"`
	StageSummary vegeta.Metrics `json:"stage_summary"`
	Steps        []StepResult   `json:"steps"`
}

type StepResult struct {
	StepSummary vegeta.Metrics `json:"step_summary"`
}
