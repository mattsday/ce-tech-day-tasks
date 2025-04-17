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
This app manages the various lifecycle phases for the tech day.

Most of this is hardcoded for ease of use rather than repeatability
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"cloud.google.com/go/firestore"
)

var (
	hostPID string = "q2-25-tech-day-host"
	action  string
	db      *firestore.Client
)

func main() {
	flag.StringVar(&hostPID, "host", hostPID, "Optional: Host Project ID")
	flag.StringVar(&action, "action", "", "Action to enable")
	flag.Parse()

	if action == "" {
		log.Fatalf("Expected 'action' flag")
	}

	if hostPID == "" {
		hostPID = os.Getenv("PROJECT_ID")
	}

	if hostPID == "" {
		log.Fatalf("Expected 'host' flag or PROJECT_ID environment")
	}

	// Init Firestore
	ctx := context.Background()
	var err error
	// Use Firestore as a global variable
	db, err = firestore.NewClient(ctx, hostPID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer db.Close()

	switch action {
	case "lock":
		lock(ctx, false)
	case "unlock":
		lock(ctx, true)
	case "testing":
		testing(ctx)
	case "start":
		start(ctx)
	case "troubleshoot":
		troubleshoot(ctx)
	case "act1-end":
		act1End(ctx)
	case "act2":
		act2(ctx)
	case "act2-end":
		act2End(ctx)
	case "security-audit":
		log.Fatalf("TODO")
	case "act3-end":
		act3End(ctx)
	default:
		log.Fatalf("Unknown action: %s", action)
	}
}

func lock(ctx context.Context, lock bool) {
	tasksRef := db.Collection("tasks").Doc("tasks")

	err := db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		allTasks, err := tx.Get(tasksRef)
		if err != nil {
			return err
		}
		var t Tasks
		allTasks.DataTo(&t)

		t.Event.ScoringEnabled = lock

		err = tx.Set(tasksRef, t)

		return err
	})

	if err != nil {
		log.Fatalf("%v", err)
	}

}

func testing(ctx context.Context) {
	allowList := []string{}
	blockList := []string{}
	disableGroups := []string{}
	enableGroups := []string{"Act 1", "Act 2", "Act 3"}

	tasksRef := db.Collection("tasks").Doc("tasks")
	err := db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var err error
		allTasks, err := tx.Get(tasksRef)
		if err != nil {
			return err
		}
		var t Tasks
		allTasks.DataTo(&t)
		t.Event.ScoringEnabled = true

		taskWriteBack := disableTasks(t, tx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: false, lbHidden: false})
		for _, v := range taskWriteBack {
			err = tx.Set(db.Collection("tasks").Doc(v.Task.ID), v)
			if err != nil {
				return err
			}
		}
		err = tx.Set(tasksRef, t)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to retrieve tasks: %v", err)
	}
}

func start(ctx context.Context) {
	allowList := []string{}
	blockList := []string{"act1-task2"}
	disableGroups := []string{"Act 2", "Act 3", "The End"}
	enableGroups := []string{"Act 1"}
	updateList(ctx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: true, lbHidden: true})
}

func troubleshoot(ctx context.Context) {
	allowList := []string{}
	blockList := []string{}
	disableGroups := []string{"Act 2", "Act 3", "The End"}
	enableGroups := []string{"Act 1"}
	updateList(ctx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: true, lbHidden: true})
}

func act2(ctx context.Context) {
	allowList := []string{}
	blockList := []string{}
	disableGroups := []string{"Act 1", "Act 3", "The End"}
	enableGroups := []string{"Act 2"}
	updateList(ctx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: false, lbHidden: false})
}

// End Act 1 - disable all tasks and score the consequences task in act 2
func act1End(ctx context.Context) {
	// Allow list to keep
	allowList := []string{"act1-task2"}
	blockList := []string{"act2-task1"}
	disableGroups := []string{"Act 1"}
	enableGroups := []string{"Act 2"}

	scoreParts := []string{"basic", "fix", "database"}

	bonus := 500
	punishment := -500

	// Retrieve tasks from Firestore
	// Run as a transaction to ensure consistency
	tasksRef := db.Collection("tasks").Doc("tasks")
	scoreRef := db.Collection("scores")
	taskId := "act1-task2"
	err := db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		allTasks, err := tx.Get(tasksRef)
		if err != nil {
			return err
		}
		var t Tasks
		allTasks.DataTo(&t)

		// Disable tasks as required
		taskWriteBack := disableTasks(t, tx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: false, lbHidden: false})

		scoreWriteBack := make(map[string]ScoreSchema)

		// Parse teams and grant / deny points
		iter := tx.Documents(scoreRef)
		for {
			doc, err := iter.Next()
			if err != nil {
				break
			}
			var s ScoreSchema
			doc.DataTo(&s)
			successful := true
			for _, v := range scoreParts {
				key := fmt.Sprintf("%v_%v", taskId, v)
				if s.Tasks[key] < 350 {
					successful = false
					s.Tasks[key] = punishment
				}
				s = updateScoreTotal(s, taskId)
			}
			if successful {
				s.Tasks["act2-task1_part1"] = bonus
			} else {
				s.Tasks["act2-task1_part1"] = -1
			}
			s = updateScoreTotal(s, "act2-task1")
			scoreWriteBack[doc.Ref.ID] = s
		}

		// Now write back to Firestore
		// Update tasks
		for _, v := range taskWriteBack {
			err = tx.Set(db.Collection("tasks").Doc(v.Task.ID), v)
			if err != nil {
				return err
			}
		}
		err = tx.Set(tasksRef, t)
		if err != nil {
			return err
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

	if err != nil {
		log.Fatalf("Failed to retrieve tasks: %v", err)
	}
}

func act2End(ctx context.Context) {
	// Allow list to keep
	allowList := []string{"act1-task2"}
	blockList := []string{}
	disableGroups := []string{"Act 1", "Act 2"}
	enableGroups := []string{"Act 3"}

	// Retrieve tasks from Firestore
	// Run as a transaction to ensure consistency
	tasksRef := db.Collection("tasks").Doc("tasks")

	err := db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		allTasks, err := tx.Get(tasksRef)
		if err != nil {
			return err
		}
		var t Tasks
		allTasks.DataTo(&t)

		// t.Event.ScoringEnabled = false

		// Disable tasks as required
		taskWriteBack := disableTasks(t, tx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: false, lbHidden: false})

		// Now write back to Firestore
		// Update tasks
		for _, v := range taskWriteBack {
			err = tx.Set(db.Collection("tasks").Doc(v.Task.ID), v)
			if err != nil {
				return err
			}
		}
		err = tx.Set(tasksRef, t)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to retrieve tasks: %v", err)
	}
}

func act3End(ctx context.Context) {
	// Allow list to keep
	allowList := []string{}
	blockList := []string{}
	disableGroups := []string{"Act 1", "Act 2", "Act 3"}
	enableGroups := []string{}

	// Retrieve tasks from Firestore
	// Run as a transaction to ensure consistency
	tasksRef := db.Collection("tasks").Doc("tasks")

	err := db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		allTasks, err := tx.Get(tasksRef)
		if err != nil {
			return err
		}
		var t Tasks
		allTasks.DataTo(&t)

		// Disable tasks as required
		taskWriteBack := disableTasks(t, tx, TaskModification{disableGroups: disableGroups, enableGroups: enableGroups, allowList: allowList, blockList: blockList, hidden: false, lbHidden: false})

		// Now write back to Firestore
		// Update tasks
		for _, v := range taskWriteBack {
			err = tx.Set(db.Collection("tasks").Doc(v.Task.ID), v)
			if err != nil {
				return err
			}
		}
		err = tx.Set(tasksRef, t)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to retrieve tasks: %v", err)
	}
}

type TaskModification struct {
	disableGroups []string
	enableGroups  []string
	allowList     []string
	blockList     []string
	hidden        bool
	lbHidden      bool
}

// Update a list of tasks without changing points
func updateList(ctx context.Context, mod TaskModification) {
	tasksRef := db.Collection("tasks").Doc("tasks")
	err := db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var err error
		allTasks, err := tx.Get(tasksRef)
		if err != nil {
			return err
		}
		var t Tasks
		allTasks.DataTo(&t)
		t.Event.ScoringEnabled = true
		taskWriteBack := disableTasks(t, tx, mod)
		for _, v := range taskWriteBack {
			err = tx.Set(db.Collection("tasks").Doc(v.Task.ID), v)
			if err != nil {
				return err
			}
		}
		err = tx.Set(tasksRef, t)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to retrieve tasks: %v", err)
	}
}

func disableTasks(t Tasks, tx *firestore.Transaction, mod TaskModification) map[string]TaskSchema {
	writeBack := make(map[string]TaskSchema)

	for k, v := range t.Tasks {
		// Get the full task to update it
		taskDoc, err := tx.Get(db.Collection("tasks").Doc(v.ID))

		if err != nil {
			log.Printf("Warning: %v", err)
			continue
		}

		var task TaskSchema
		taskDoc.DataTo(&task)

		if slices.Contains(mod.disableGroups, v.Group) && !slices.Contains(mod.allowList, v.ID) {
			// Disable the task
			task.Task.Enabled = false
			// Hide tasks from the UI / Leaderboard if configured
			if mod.hidden {
				task.Task.Hidden = true
			}
			if mod.lbHidden {
				task.Task.LBHidden = true
			}
		} else if slices.Contains(mod.blockList, v.ID) {
			log.Printf("%v is blocklisted", v.ID)
			if mod.hidden {
				task.Task.Hidden = true
			}
			if mod.lbHidden {
				task.Task.LBHidden = true
			}

		} else if slices.Contains(mod.enableGroups, v.Group) && !slices.Contains(mod.blockList, v.ID) {
			// Enable the task and unhide it
			task.Task.Hidden = false
			task.Task.LBHidden = false
			task.Task.Enabled = true
		}
		v = task.Task
		writeBack[task.Task.ID] = task
		t.Tasks[k] = task.Task
	}

	return writeBack

}

func updateScoreTotal(data ScoreSchema, taskId string) ScoreSchema {
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
