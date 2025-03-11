Now, let's turn up the heat! We're using Gemini to analyze each resume and give us a 'Suitability Score'. Think of it as a spice rating for each candidate – how well do they match our job description?

### Task

1. Create a Gemini Agent: "Let's create a new agent in Agentspace. We'll call it 'The Resume Rater'."
2. Configure the Agent:
    * "Give it a clear prompt. We need to tell Gemini exactly what we want. Something like: 'You're a recruitment expert, and you're going to score resumes based on this job description.'"
    * "Instructions: 'Read the Job Description. Then, read each resume. Compare the skills and experience to the job requirements. Give each candidate a score from 0 to 100, where 100 is a perfect match. Tell us why you gave them that score. Consider experience, skills, and education. Give me the Name and Score and reasoning.'"
    * "Connect the 'resumes.csv' and the Job Description as input. We need to feed Gemini the data!"
3. Run the Agent: "Let Gemini do its thing. Let's see who's got the right blend of skills."
4. Review Results: "Check the scores and reasoning. Did Gemini nail it? Or do we need to tweak the recipe?"
5. TODO - how do we verify this?
