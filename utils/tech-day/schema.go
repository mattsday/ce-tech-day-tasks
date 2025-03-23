package main

import "time"

type Metadata struct {
	Version int `yaml:"version"`
}

type Tasks struct {
	Metadata Metadata        `yaml:"metadata" firestore:"metadata"`
	Event    Event           `yaml:"event" firestore:"event"`
	Tasks    map[string]Task `yaml:"tasks" firestore:"tasks"`
}

type Event struct {
	Name           string `yaml:"name" firestore:"name"`
	Logo           Logo   `yaml:"logo" firestore:"logo"`
	ScoringEnabled bool   `yaml:"scoring_enabled" firestore:"scoring_enabled"`
	MaxPoints      int    `yaml:"max_points" firestore:"max_points"`
	Theme          string `yaml:"theme" firestore:"theme"`
	Intro          string `yaml:"intro" firestore:"intro"`
}

type Logo struct {
	Path   string `yaml:"path" firestore:"path"`
	Width  int    `yaml:"width" firestore:"width"`
	Height int    `yaml:"height" firestore:"height"`
}

type Task struct {
	Name        string            `yaml:"name" firestore:"name"`
	ID          string            `yaml:"id" firestore:"id"`
	Description string            `yaml:"description" firestore:"description"`
	Overview    []string          `yaml:"overview" firestore:"overview"`
	Enabled     bool              `yaml:"enabled" firestore:"enabled"`
	Hidden      bool              `yaml:"hidden" firestore:"hidden"`
	LBHidden    bool              `yaml:"lb_hidden" firestore:"lb_hidden"`
	MaxPoints   int               `yaml:"max_points" firestore:"max_points"`
	Group       string            `yaml:"group" firestore:"group"`
	Metadata    map[string]string `yaml:"metadata" firestore:"metadata"`

	// Non-persisted fields
	TFEnabled   bool        `yaml:"tf_enabled" firestore:"-"`
	TFVars      []string    `yaml:"tf_vars" firestore:"-"`
	BQDatasets  []BQDataset `yaml:"bq_datasets" firestore:"-"`
	UploadFiles []string    `yaml:"upload_files" firestore:"-"`
}

type BQDataset struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Tables      []BQTable `yaml:"tables"`
}

type BQTable struct {
	Name        string `yaml:"name"`
	Source      string `yaml:"source"`
	Schema      string `yaml:"schema"`
	Description string `yaml:"description"`
}

type Part struct {
	Name             string   `yaml:"name" firestore:"name"`
	ID               string   `yaml:"id" firestore:"id"`
	Type             string   `yaml:"type" firestore:"type"`
	UploadText       string   `yaml:"upload_text" firestore:"upload_text"`
	MaxPoints        int      `yaml:"max_points" firestore:"max_points"`
	InstructionsLink string   `yaml:"instructions_link" firestore:"instructions_link"`
	Challenge        bool     `yaml:"challenge,omitempty" firestore:"challenge"`
	GoodExamples     []string `yaml:"good_examples,omitempty" firestore:"good_examples"`
	LLMInstructions  string   `yaml:"llm_instructions,omitempty" firestore:"llm_instructions"`
	Open             bool     `yaml:"open,omitempty" firestore:"open"`
	Component        string   `yaml:"component,omitempty" firestore:"component"`
	Hidden           bool     `yaml:"hidden" firestore:"hidden"`
	DependsOn        []string `yaml:"depends_on,omitempty" firestore:"depends_on"`
}

type TaskSchema struct {
	Metadata Metadata `yaml:"metadata" firestore:"metadata"`
	Task     Task     `yaml:"task" firestore:"task"`
	Parts    []Part   `yaml:"parts" firestore:"parts"`
}

type ScoreSchema struct {
	Totals struct {
		Total        int `firestore:"total"`
		BonusTotal   int `firestore:"bonus_total"`
		RegularTotal int `firestore:"regular_total"`
		OldTotal     int `firestore:"old_total"`
	} `firestore:"totals"`
	Screenshots  map[string]string `firestore:"screenshots"`
	BonusScores  map[string]int    `firestore:"bonus_scores"`
	TaskTotals   map[string]int    `firestore:"task_totals"`
	Reviewed     map[string]bool   `firestore:"reviewed"`
	Tasks        map[string]int    `firestore:"tasks"`
	ScoreUpdated time.Time         `firestore:"score_updated"`
}
