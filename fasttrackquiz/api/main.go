package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Question struct {
	Id       int       `json:"id"`
	Question string    `json:"question"`
	Options  [3]string `json:"options"`
	Correct  string    `json:"correct"`
}
type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Time float64  `json:"time"`
}

var Questions = [6]Question{

	{
		Id:       1,
		Question: "What year was Go announced?",
		Options:  [3]string{"1: 2012", "2: 2011", "3: 2009"},
		Correct:  "3",
	},
	{
		Id:       2,
		Question: "What year was c# released?",
		Options:  [3]string{"1: 2000", "2: 1998", "3: 2002"},
		Correct:  "3",
	},
	{
		Id:       3,
		Question: "Who is one of the founders of Apple?",
		Options:  [3]string{"1: Elon Musk", "2: Steve Jobs", "3: Jeff Bezos"},
		Correct:  "2",
	},
	{
		Id:       4,
		Question: "How long does it take for an apple to rot?",
		Options:  [3]string{"1: 1 month", "2: 2 months", "3: 3 months"},
		Correct:  "1",
	},
	{
		Id:       5,
		Question: "What country has the highest life expectancy?",
		Options:  [3]string{"1: Sweden", "2: Hongkong", "3: Portugal"},
		Correct:  "2",
	},
	{
		Id:       6,
		Question: "What is the most common surname in the United States? ",
		Options:  [3]string{"1: Smith", "2: 2 Johnson", "3: Williams"},
		Correct:  "1",
	},
}

var Scores = []Score{
	{
        Name: "Elias",
        Score: 6,
        Time: 4.3332563,
    },
	{
		Name: "Anas",
		Score: 2,
		Time: 7.0841213,
	},
}

func GetQuestions(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Questions)
}

func getQuestion(context *gin.Context) {
	id := context.Param("id")
	question, err := GetQuestionById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"messege": "Question not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, question.Question)
	context.IndentedJSON(http.StatusOK, question.Options)
}

func GetQuestionById(id string) (*Question, error) {
	newid, _ := strconv.Atoi(id)

	for i, q := range Questions {
		if q.Id == newid {
			return &Questions[i], nil
		}
	}
	return nil, errors.New("Question not found")
}

func AddAnswer(context *gin.Context) {
	var newScore Score
	if err := context.BindJSON(&newScore); err != nil {
		return
	}

	Scores = append(Scores, newScore)

	context.IndentedJSON(http.StatusCreated, newScore)
}

func GetScore(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, Scores)
}

//Score calc

func main() {
	router := gin.Default()
	router.GET("/questions", GetQuestions)
	router.GET("/questions/:id", getQuestion)

	router.GET("/answers", GetScore)
	router.POST("/answers", AddAnswer)

	router.Run("localhost:9090")
}
