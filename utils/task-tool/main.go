/**
 * Copyright 2024 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
This app parses tasks and uploads them into Firestore.

It's very, very hastily written and a mess.
*/

package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"flag"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"gopkg.in/yaml.v3"
)

var (
	hostPID      string
	baseFolder   string
	bucket       string
	upload       bool
	uploadImages bool
	tasktf       bool
	tfFile       string
)

type Metadata struct {
	Version int      `yaml:"version"`
	Authors []Author `yaml:"authors" firestore:"authors"`
}

type Author struct {
	Name  string `yaml:"name" firestore:"name"`
	Email string `yaml:"email" firestore:"email"`
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
	UserFiles   []UserFiles `yaml:"user_files" firestore:"-"`
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

type UserFiles struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
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
	MaxFiles         int      `yaml:"max_files,omitempty" firestore:"max_files"`
}

type TaskSchema struct {
	Metadata Metadata `yaml:"metadata" firestore:"metadata"`
	Task     Task     `yaml:"task" firestore:"task"`
	Parts    []Part   `yaml:"parts" firestore:"parts"`
}

func main() {
	flag.StringVar(&hostPID, "host-pid", "", "Host Project ID")
	flag.StringVar(&baseFolder, "base-folder", "", "Root folder for tasks")
	flag.StringVar(&bucket, "bucket", "", "Storage bucket to use")
	flag.BoolVar(&upload, "upload", true, "Upload assets to GCS")
	flag.BoolVar(&uploadImages, "upload-images", true, "Upload images (can be slow)")
	flag.BoolVar(&tasktf, "tf-only", false, "Whether to only generate task TF config")
	flag.StringVar(&tfFile, "tf-file", "terraform.tfvars", "Terraform output file")

	flag.Parse()

	if tasktf {
		tfOnly()
		return
	}

	if bucket == "" {
		log.Fatalf("No bucket specified\n")
	}

	// Check if the base folder exists
	if _, err := os.Stat(baseFolder); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("%v does not exist", baseFolder)
	}

	// Init Firestore
	ctx := context.Background()
	db, err := firestore.NewClient(ctx, hostPID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer db.Close()

	// Init GCS
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer gcs.Close()

	// Get current tasks from Firestore and track if they should be kept
	keepTasks := make(map[string]bool)
	log.Printf("Reading current tasks\n")
	iter := db.Collection("tasks").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error %v\n", err)
			continue
		}
		keepTasks[doc.Ref.ID] = false
	}
	// Always keep the tasks doc which is an overview
	keepTasks["tasks"] = true

	// Parse base tasks.yaml file
	tasksFile := fmt.Sprintf("%v/tasks.yaml", baseFolder)

	t := Tasks{}

	// Parse it if it exists, otherwise fall-back to defaults
	if _, err := os.Stat(tasksFile); err == nil {
		log.Printf("Parsing tasks.yaml\n")
		data, err := os.ReadFile(tasksFile)
		if err != nil {
			log.Fatal(err)
		}
		err = yaml.Unmarshal([]byte(data), &t)
		if err != nil {
			log.Fatal(err)
		}
		if upload {
			if uploadImages {
				if t.Event.Logo.Path != "" {
					folder := filepath.Dir(tasksFile)
					err = uploadFile(ctx, gcs, t.Event.Logo.Path, "assets", folder)
					if err != nil {
						log.Printf("Warning: %v\n", err)
					}
				}
			}
			if t.Event.Intro != "" {
				folder := filepath.Dir(tasksFile)
				err = uploadFile(ctx, gcs, t.Event.Intro, "assets", folder)
				if err != nil {
					log.Printf("Warning: %v\n", err)
				}
			}
		}
	} else {
		log.Printf("Warning: no tasks.yaml file - using defaults\n")
		t.Event.Name = "hacksday"
		t.Event.ScoringEnabled = true
		t.Event.Theme = "The theme of the event is using AI tools in Google Cloud. Participants also have access to a number of food and drink items in the room and nutritional datasets."
	}

	// Keep track of points available
	total := 0

	// Keep a tally of all tasks
	taskRecord := make(map[string]Task)

	// Traverse base folder for tasks
	err = filepath.WalkDir(baseFolder, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "disabled" {
			return filepath.SkipDir
		}
		if d.Name() == "task.yaml" && !d.IsDir() {
			t, err := parseTask(ctx, gcs, path, db, taskRecord, keepTasks)
			if err != nil {
				log.Printf("Warning: %v\n", err)
			} else {
				total += t
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range keepTasks {
		if !v {
			log.Printf("Removing dangling task '%v'", k)
			dr := db.Collection("tasks").Doc(k)
			_, err = dr.Delete(ctx)
			if err != nil {
				log.Printf("Warning: %v\n", err)
			}
		}
	}

	t.Tasks = taskRecord

	t.Event.MaxPoints = total

	// Write overall task list into Firestore
	dr := db.Collection("tasks").Doc("tasks")
	_, err = dr.Set(ctx, t)
	if err != nil {
		log.Fatal(err)
	}
}

func parseTask(ctx context.Context, gcs *storage.Client, taskFile string, db *firestore.Client, taskRecord map[string]Task, keepTasks map[string]bool) (int, error) {
	// log.Printf("Parsing %v\n", taskFile)
	folder := filepath.Dir(taskFile)

	if _, err := os.Stat(taskFile); errors.Is(err, os.ErrNotExist) {
		return 0, err
	}

	data, err := os.ReadFile(taskFile)
	if err != nil {
		return 0, err
	}

	t := TaskSchema{}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return 0, err
	}
	log.Printf("Creating / Updating Task: %v\n", t.Task.Name)

	total := 0

	if t.Task.ID == "tasks" {
		return 0, fmt.Errorf("invalid task name 'tasks' in %v", taskFile)
	}

	// Upload any additional assets
	for _, asset := range t.Task.UploadFiles {
		if upload {
			err = uploadFile(ctx, gcs, asset, fmt.Sprintf("instructions/%v", t.Task.ID), folder)
			if err != nil {
				log.Printf("Warning: %v\n", err)
			}
		}
	}

	keepTasks[t.Task.ID] = true

	// Loop through parts and upload assets
	for _, part := range t.Parts {
		total += part.MaxPoints
		if upload {
			if part.InstructionsLink != "" {
				err = uploadFile(ctx, gcs, part.InstructionsLink, fmt.Sprintf("instructions/%v", t.Task.ID), folder)
				if err != nil {
					log.Printf("Warning: %v\n", err)
				}
			}
		}
		if upload && uploadImages {
			for _, sample := range part.GoodExamples {
				if sample != "" {
					err = uploadFile(ctx, gcs, sample, fmt.Sprintf("instructions/%v", t.Task.ID), folder)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}
			}
		}
	}

	t.Task.MaxPoints = total

	taskRecord[t.Task.ID] = t.Task

	// Bang task into Firestore
	dr := db.Collection("tasks").Doc(t.Task.ID)
	_, err = dr.Set(ctx, t)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func uploadFile(ctx context.Context, gcs *storage.Client, filename string, destination string, folder string) error {
	assetPath := fmt.Sprintf("%v/%v", folder, filename)
	uploadedAssetPath := fmt.Sprintf("%v/%v", destination, filename)
	// fmt.Printf("Uploading file: '%v' to '%v'\n", assetPath, uploadedAssetPath)
	if _, err := os.Stat(assetPath); errors.Is(err, os.ErrNotExist) {
		return err
	}
	f, err := os.Open(assetPath)
	if err != nil {
		return err
	}
	defer f.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	o := gcs.Bucket(bucket).Object(uploadedAssetPath)
	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func tfOnly() {
	if hostPID == "" || baseFolder == "" {
		fmt.Printf("Expected host-pid and base-folder variables\n")
		os.Exit(1)
	}

	var tfFlags = make(map[string]bool)
	var bqDatasets = make(map[string]BQDataset)
	var userFiles = make(map[string]UserFiles)

	// Traverse base folder for tasks
	err := filepath.WalkDir(baseFolder, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == "disabled" {
			return filepath.SkipDir
		}
		if d.Name() == "task.yaml" && !d.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			} else {
				t := TaskSchema{}
				err = yaml.Unmarshal([]byte(data), &t)
				if len(t.Task.TFVars) != 0 {
					for _, v := range t.Task.TFVars {
						// If a task is enabled anywhere that takes priority
						if !tfFlags[v] {
							tfFlags[v] = t.Task.TFEnabled
						}
					}
				}
				if len(t.Task.BQDatasets) != 0 {
					for _, v := range t.Task.BQDatasets {
						bqDatasets[v.Name] = v
					}
				}
				if len(t.Task.UserFiles) != 0 {
					for _, v := range t.Task.UserFiles {
						userFiles[fmt.Sprintf("%v --> %v", v.Source, v.Destination)] = v
					}
				}
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(tfFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	w := bufio.NewWriter(f)

	for k, v := range tfFlags {
		w.WriteString(fmt.Sprintf("%v = %v\n", k, v))
	}

	w.WriteString("bq_datasets = [\n")
	for _, v := range bqDatasets {
		w.WriteString("\t{\n")
		w.WriteString(fmt.Sprintf("\t\tname = \"%v\"\n", v.Name))
		w.WriteString(fmt.Sprintf("\t\tdescription = \"%v\"\n", v.Description))
		w.WriteString("\t\ttables = [\n")
		for _, v1 := range v.Tables {
			w.WriteString("\t\t\t{\n")
			w.WriteString(fmt.Sprintf("\t\t\t\tname = \"%v\"\n", v1.Name))
			w.WriteString(fmt.Sprintf("\t\t\t\tsource = \"%v\"\n", v1.Source))
			w.WriteString(fmt.Sprintf("\t\t\t\tschema = \"%v\"\n", v1.Schema))
			w.WriteString(fmt.Sprintf("\t\t\t\tdescription = \"%v\"\n", v1.Description))
			w.WriteString("\t\t\t},\n")
		}
		w.WriteString("\t\t]\n")
		w.WriteString("\t},\n")
	}
	w.WriteString("]\n")
	w.WriteString("user_files = [\n")
	for _, v := range userFiles {
		w.WriteString("\t{\n")
		w.WriteString(fmt.Sprintf("\t\tsource = \"%v\"\n", v.Source))
		w.WriteString(fmt.Sprintf("\t\tdestination = \"%v\"\n", v.Destination))
		w.WriteString("\t},\n")

	}
	w.WriteString("]\n")
	w.Flush()
}
