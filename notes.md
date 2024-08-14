## Notes and Todos

- so for each practice question, I want to have a file (I think JSON) that will list all the questions for that concept.
- This file should be downloaded with all the questions for a concept
- "questions.json" (I think that makes sense for the name)
- possible question types
  - MCQ: these are multiple choice questions. User can select one answer from a list of answers
  - edit: these questions allow the user to edit a file and add code and then run tests
- something like this

```json
[
  {
    "title": "question 1",
    "type": "mcq",
    "question": "what is 1+1",
    "answers": ["1", "2", "11"],
    "answer": "2"
  },
  {
    "title": "question 2",
    "type": "edit",
    "file": "main.go",
    "testFile": "hello-world_test.go"
  }
]
```

- We pull the folder for a concept (if it does not exist), add the questions to a DB with complete set to false, allow the user to answer each question, and update the DB if they get it right
