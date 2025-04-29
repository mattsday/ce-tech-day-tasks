package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
	"google.golang.org/api/iterator"
)

// Maximum retries to score a team
const maxAttempts = 5

func main() {
	// Get the task key from the environment variable.
	taskId := os.Getenv("TASK_ID")
	if taskId == "" {
		log.Fatal("TASK_ID environment variable not set")
	}
	partId := os.Getenv("PART_ID")
	if partId == "" {
		log.Fatal("PART_ID environment variable not set")
	}
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("PROJECT_ID environment variable not set")
	}

	log.Printf("Evaluating asset submissions for %v", taskId)

	taskKey := fmt.Sprintf("%s_%s", taskId, partId)

	// Use the application default credentials.
	ctx := context.Background()
	db, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer db.Close()

	// Get the task document from Firestore.
	doc, err := db.Collection("tasks").Doc(taskId).Get(ctx)
	if err != nil {
		log.Fatalf("Error getting task document: %v", err)
	}

	// Unmarshal the document data into a TaskSchema struct.
	var taskSchema TaskSchema
	if err := doc.DataTo(&taskSchema); err != nil {
		log.Fatalf("Error unmarshalling task document: %v", err)
	}

	// Find the part with the matching partId.
	var foundPart Part
	found := false
	for _, part := range taskSchema.Parts {
		if part.ID == partId {
			foundPart = part
			found = true
			break
		}
	}

	if !found {
		log.Fatalf("Part with ID '%s' not found in task '%s'", partId, taskId)
	}

	// Construct the GCS path
	gcsPath := fmt.Sprintf("instructions/%s/%s", taskId, foundPart.InstructionsLink)
	bucketName := fmt.Sprintf("%s-score-assets", projectID)

	// Initialize GCS client
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer gcsClient.Close()

	aiClient, err := genai.NewClient(ctx, projectID, "us-central1")
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}
	gemini := aiClient.GenerativeModel("gemini-2.5-flash-preview-04-17")
	gemini.GenerationConfig.ResponseMIMEType = "application/json"
	// Set a low temperature for consistent judgements
	gemini.Temperature = genai.Ptr[float32](0.1)

	// Get the GCS object
	rc, err := gcsClient.Bucket(bucketName).Object(gcsPath).NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS reader: %v", err)
	}
	defer rc.Close()

	// Read the contents of the file into a string
	fileContents, err := io.ReadAll(rc)
	if err != nil {
		log.Fatalf("Failed to read GCS file: %v", err)
	}

	// Store the file contents in a string variable
	instructions := string(fileContents)

	scoreWriteBack := make(map[string]ScoreSchema)
	feedbackWriteBack := make(map[string]Feedback)

	// Get overall event theme

	// Get the task document from Firestore.
	tasksDoc, err := db.Collection("tasks").Doc("tasks").Get(ctx)
	if err != nil {
		log.Fatalf("Error getting task document: %v", err)
	}

	// Unmarshal the document data into a TaskSchema struct.
	var tasks Tasks
	if err := tasksDoc.DataTo(&tasks); err != nil {
		log.Fatalf("Error unmarshalling tasks document: %v", err)
	}

	mimeType := "application/pdf"

	if val, ok := taskSchema.Task.Metadata["mimeType"]; ok {
		log.Printf("Setting MIME type to %v\n", val)
		mimeType = val
	}

	// Get overall event theme
	theme := tasks.Event.Theme
	name := tasks.Event.Name
	taskName := taskSchema.Task.Name
	partName := foundPart.Name
	maxPoints := foundPart.MaxPoints
	extraInstructions := "None"
	if foundPart.LLMInstructions != "" {
		extraInstructions += foundPart.LLMInstructions
	}

	gemini.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(fmt.Sprintf(`
		Your name is "Judge Gemini 🧑‍⚖". You are a judge for an event called %v.
	
	Your role is to judge uploaded content from participants and determine if they have adequately completed the task and review if they can improve. The majority of submissions will be screenshots showing various tasks completed in Google Cloud and other related services.
	
	Provide each submission with a score and feedback. Award zero points if the task is obviously not completed. Award negative points if you suspect cheating or gaming the system. Provide your rationale in your feedback by showing a fair breakdown of your score.
	
	Participants will receive your feedback directly, so also include tips or hints that can help them improve their score.
	
	Explain your score by breaking it down in a bullet-pointed list.
	
	Provide your response in JSON format with two fields:
	score - an integer
	verdict - markdown formatted response
	
	Event Theme:
	%v
	
	Task Name:
	"%v"
	
	Part Name:
	"%v"
	
	Maximum Points:
	"%v"
	
	Instructions given to users for this task:
	%v
	
	Extra Instructions:
	%v`, name, theme, taskName, partName, maxPoints, instructions, extraInstructions))},
	}

	prompt := genai.Text("The attached document is this team's submission")

	// Write scores and feedback in a transaction
	err = db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// iter := db.Collection("scores").Documents(ctx)
		iter := tx.Documents(db.Collection("scores"))

		// Iterate over the "scores" collection.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("iter.Next: %v", err)
			}

			var s ScoreSchema
			doc.DataTo(&s)

			// Extract the value associated with the task key.
			taskValue, ok := s.Screenshots[taskKey]
			if !ok {
				log.Printf("Document %s 'screenshots' map does not have key '%s'\n", doc.Ref.ID, taskKey)
				continue
			}

			gcsURL := fmt.Sprintf("gs://%s/misc/%s/%s/%s/%v", bucketName, doc.Ref.ID, taskId, partId, taskValue)

			file := genai.FileData{
				MIMEType: mimeType,
				FileURI:  gcsURL,
			}

			attempt := 0
			complete := false
			for attempt < maxAttempts {
				log.Printf("Attempt %v of %v for team %v\n", attempt+1, maxAttempts, doc.Ref.ID)
				attempt++

				res, err := gemini.GenerateContent(ctx, file, prompt)
				if err != nil {
					log.Printf("File: %v, error generating content: %v", gcsURL, err)
					continue
				}
				if len(res.Candidates) == 0 ||
					len(res.Candidates[0].Content.Parts) == 0 {
					log.Printf("empty response from model")
					continue
				}
				// Unmarshal the JSON response into the Verdict struct
				var verdict Verdict
				jsonString, ok := res.Candidates[0].Content.Parts[0].(genai.Text)
				if !ok {
					log.Printf("Could not convert response to string")
					continue
				}
				err = json.Unmarshal([]byte(jsonString), &verdict)
				if err != nil {
					log.Printf("Error unmarshalling JSON: %v, JSON: %v", err, jsonString)
					continue
				}
				complete = true

				s.Tasks[taskKey] = verdict.Score
				s = updateScoreTotal(s, taskId)
				scoreWriteBack[doc.Ref.ID] = s

				feedbackWriteBack[fmt.Sprintf("%v-%v-%v", doc.Ref.ID, taskId, partId)] = Feedback{AIFeedback: verdict.Verdict, AIScore: verdict.Score}
				break
			}
			if attempt+1 >= maxAttempts && !complete {
				log.Printf("Warning: %v took more than %v attempts - content NOT evaluated\n", doc.Ref.ID, attempt)
			}
		}
		log.Printf("Analysis complete, writing scores & feedback\n")
		// Update scores
		for k, v := range scoreWriteBack {
			err = tx.Set(db.Collection("scores").Doc(k), v)
			if err != nil {
				return err
			}
		}
		// Update feedback
		for k, v := range feedbackWriteBack {
			err = tx.Set(db.Collection("feedback").Doc(k), v)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error running transaction: %v", err)
	}

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

type Verdict struct {
	Score   int    `json:"score"`
	Verdict string `json:"verdict"`
}

type Feedback struct {
	AIFeedback string `json:"ai-feedback" firestore:"ai-feedback"`
	AIScore    int    `json:"ai-score" firestore:"ai-score"`
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

type Metadata struct {
	Version int      `yaml:"version"`
	Authors []Author `yaml:"authors" firestore:"authors"`
}

type Author struct {
	Name  string `yaml:"name" firestore:"name"`
	Email string `yaml:"email" firestore:"email"`
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
