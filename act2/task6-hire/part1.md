Cymbal Supplements is experiencing explosive growth thanks to the newly discovered age-reversing properties of their hot sauces! To keep up with demand and continue their innovative product development, the R&D team needs to bring on a talented Food Scientist specializing in sauce development.

### Dataset

You have been provided with access to a dataset of resumes in Cloud Storage, accessible through Agentspace. This dataset contains PDF files for numerous applicants, each including information such as name, email, skills, experience, and education.

The Hiring Manager has also shared the following Job Description (JD) for the Food Scientist (Sauce Development) role:

```plaintext
R&D: Food Scientist (Sauce Development)
Title: Food Scientist (Sauce Development)
Department: Research & Development
Reports To: R&D Manager
Job Summary: Develop and refine new Cymbal Sauce recipes, ensuring quality, consistency, and innovation. Conduct sensory evaluations and shelf-life testing.
Responsibilities:

Formulate and test new sauce variations.
Optimize existing recipes for improved flavor and stability.
Conduct sensory analysis and maintain accurate records.
Collaborate with suppliers for ingredient sourcing.
Ensure compliance with food safety regulations. Qualifications:
Bachelor's degree in Food Science or related field.
Experience in sauce or condiment development.
Strong understanding of food chemistry and microbiology.
```

### Task

Using Agentspace, run a query that processes the dataset of resume PDFs and compares each candidate's qualifications against the provided Job Description. Your goal is to have Gemini analyze each resume and generate a "Suitability Score" reflecting how well the candidate's skills and experience align with the requirements of the Food Scientist role.

::info[Each candidate in the resume dataset must have a calculated "Suitability Score" generated by Gemini's analysis within Agentspace.]

1. Log in to the Retail Agentspace Demo ([see the earlier task for instructions on how to do this](/task/act1-task3))
2. Create a new Chat Conversation and select **retail resume datastore** as your Data Source
3. Formulate a query using Gemini within Agentspace. This query should:
   - Access the resume data.
   - Incorporate the provided Job Description (you can directly input the text into the query).
   - Instruct Gemini to analyze each resume against the JD.
   - Request Gemini to output a "Suitability Score" for each candidate based on their alignment with the JD.
4. Ensure you can view the results of your query, including the Suitability Scores for all candidates.

Provide a screenshot showing the output of your Agentspace query, clearly displaying the "Suitability Score" for at least three different candidates.
