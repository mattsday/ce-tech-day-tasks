package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/vertexai/genai"
)

const defaultProject = "q2-25-tech-day-host"
const maxAttempts = 5

type ScoreRecord struct {
	TeamID   string `json:"team_id"`
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

func main() {
	// Get total dragons den scores
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Printf("Setting project ID to %v", defaultProject)
		projectID = defaultProject
	}

	tasks := []string{"act1-task4_part1", "act2-task2_audition"}
	// Tech debt...
	fbTasks := []string{"act1-task4-part1", "act2-task2-audition"}

	log.Printf("Looping through teams and obtaining scores")

	// Use the application default credentials.
	ctx := context.Background()
	db, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}
	defer db.Close()

	aiClient, err := genai.NewClient(ctx, projectID, "us-central1")
	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	// Loop through scores collection
	scoreRef := db.Collection("scores")

	scores := make(map[string]ScoreRecord)
	err = db.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		iter := tx.Documents(scoreRef)
		for {
			doc, err := iter.Next()
			if err != nil {
				break
			}

			sr := ScoreRecord{
				TeamID:   doc.Ref.ID,
				Score:    0,
				Feedback: "",
			}

			var s ScoreSchema
			doc.DataTo(&s)
			// Check tasks
			for i, v := range tasks {
				if val, ok := s.Tasks[v]; ok {
					sr.Score += val
					// Check if feedback exists too
					fbDoc := db.Collection("feedback").Doc(doc.Ref.ID + "-" + fbTasks[i])
					fb, err := tx.Get(fbDoc)
					if err != nil {
						log.Printf("Error: %v", err)
					} else {
						if fb.Exists() {
							var fbs FeedbackSchema
							fb.DataTo(&fbs)
							sr.Feedback += fbs.AIFeedback
						}
					}
				}
			}
			scores[doc.Ref.ID] = sr

		}
		// log.Printf("%+v", scores)
		// Mashall scores as JSON
		j, err := json.Marshal(scores)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		// Construct Gemini prompt to sort and categorise entries and pick top 5

		prompt := genai.Text(fmt.Sprintf("Pick the top 5 entries from the following JSON payload and set the success field to true. For the rest set it to false. The most important criteria is score. For those entries where the score is equal, use the feedback to determine the best. Return the list in JSON format with the team_id, ordered by ranking (best ranked first).\n\nGive a short reason (1-2 sentences) for your selection of this team and why it's ranked where it has\n\nThe schema should include 'team_id', 'score', 'success', and 'reason' fields.\n\nJSON payload:\n%v", string(j)))

		gemini := aiClient.GenerativeModel("gemini-2.5-pro-preview-03-25")
		gemini.GenerationConfig.ResponseMIMEType = "application/json"
		// Set a low temperature for consistent judgements
		gemini.Temperature = genai.Ptr[float32](0.1)

		attempt := 0
		complete := false
		var verdict []Verdict
		log.Printf("Gemini is calculating scores and ranking")
		for attempt < maxAttempts {
			log.Printf("Attempt %v of %v", attempt+1, maxAttempts)
			attempt++
			res, err := gemini.GenerateContent(ctx, prompt)
			if err != nil {
				log.Printf("Error generating content: %v", err)
				continue
			}
			if len(res.Candidates) == 0 ||
				len(res.Candidates[0].Content.Parts) == 0 {
				log.Printf("empty response from model")
				continue
			}

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
			break
		}
		if attempt+1 >= maxAttempts && !complete {
			log.Fatalf("Warning: took more than %v attempts - content NOT evaluated\n", attempt)
		}
		log.Printf("Scores calculated - writing to database")
		// Write feedback to the database
		for _, v := range verdict {
			// Write feedback to output task
			fbRef := db.Collection("feedback").Doc(fmt.Sprintf("%v-act3-task1-finals", v.TeamID))
			fb := FeedbackSchema{
				AIFeedback: v.Reason,
			}
			err = tx.Set(fbRef, fb)
			if err != nil {
				log.Printf("Error: %v", err)
			}
		}
		// Write output to dragons collection
		// Structure in a way Firestore likes it
		fbVerdict := make(map[string]Verdict)
		apVerdict := make(map[string]Appointments)

		appointments := []string{"15:55", "16:05", "16:15", "16:25", "16:35"}

		a := 0

		for _, v := range verdict {
			fbVerdict[v.TeamID] = v
			if v.Success {
				apVerdict[v.TeamID] = Appointments{
					Time:   appointments[a],
					Score:  v.Score,
					Reason: v.Reason,
				}
				a++
			}
		}
		err = tx.Set(db.Collection("dragons").Doc("scores"), fbVerdict)
		if err != nil {
			log.Printf("Error saving scores: %v", err)
		}
		err = tx.Set(db.Collection("dragons").Doc("appointments"), apVerdict)
		if err != nil {
			log.Printf("Error saving scores: %v", err)
		}
		log.Printf("Complete")

		return nil
	})

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

type FeedbackSchema struct {
	AIFeedback string `json:"ai-feedback" firestore:"ai-feedback"`
}

type Verdict struct {
	TeamID  string `json:"team_id" firestore:"team_id"`
	Score   int    `json:"score" firestore:"score"`
	Reason  string `json:"reason" firestore:"reason"`
	Success bool   `json:"success" firestore:"success"`
}

type Appointments struct {
	Time   string `firestore:"time"`
	Score  int    `firestore:"score"`
	Reason string `firestore:"reason"`
}
