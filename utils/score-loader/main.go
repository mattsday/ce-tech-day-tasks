/*
Load a sheet of scores and apply it all at once to give some suspense on the day

Expect it with the headers:
team_id, task_id, part_id, points
*/
package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

var (
	sheet      string
	ignoreRows int
	hostPID    string
)

func main() {
	flag.StringVar(&sheet, "sheet", "scores.csv", "CSV File containing points")
	flag.IntVar(&ignoreRows, "ignore-rows", 1, "Number of top rows to ignore (e.g. CSV header row)")
	flag.StringVar(&hostPID, "host-pid", "q2-25-tech-day-host", "Host Project ID")

	flag.Parse()

	ctx := context.Background()
	db, err := firestore.NewClient(ctx, hostPID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer db.Close()

	// Iterate through CSV file
	file, err := os.Open(sheet)
	if err != nil {
		log.Fatalf("File %v not found\n", sheet)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	err = db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		scoreWriteBack := make(map[string]ScoreSchema)

		// Parse CSV into map
		for r, rv := range data {
			if r == (ignoreRows - 1) {
				continue
			}
			if len(rv) < 4 {
				continue
			}
			team := rv[0]
			task := rv[1]
			part := rv[2]
			points, err := strconv.Atoi(rv[3])
			if err != nil {
				points = 0
				continue
			}

			var s ScoreSchema
			// Do we already have the score in our writeback?
			if _, ok := scoreWriteBack[team]; ok {
				s = scoreWriteBack[team]
			} else {
				docRef := db.Collection("scores").Doc(team)
				doc, err := tx.Get(docRef)
				if err != nil {
					log.Printf("Error: %v", err)
					continue
				}
				doc.DataTo(&s)
			}
			if task == "bonus" {
				s.BonusScores[part] = points
			} else {
				taskKey := fmt.Sprintf("%v_%v", task, part)
				s.Tasks[taskKey] = points
			}
			s = updateScoreTotal(s, task)
			scoreWriteBack[team] = s
		}

		// Update scores
		for k, v := range scoreWriteBack {
			err = tx.Set(db.Collection("scores").Doc(k), v)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func updateScoreTotal(data ScoreSchema, taskId string) ScoreSchema {
	if taskId == "bonus" {
		total := 0
		for _, v := range data.BonusScores {
			total += v
		}
		data.Totals.BonusTotal = total
		data.Totals.Total = data.Totals.RegularTotal + total
		return data
	}

	taskTotal := 0
	for k, v := range data.Tasks {
		if strings.HasPrefix(k, taskId) {
			taskTotal += v
		}
	}
	data.TaskTotals[taskId] = taskTotal

	total := 0
	for _, v := range data.TaskTotals {
		total += v
	}
	data.Totals.RegularTotal = total
	data.Totals.Total = total + data.Totals.BonusTotal

	return data
}

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
