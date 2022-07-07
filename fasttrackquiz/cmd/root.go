/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fasttrackquiz",
	Short: "cli for a quiz",
	Long:  "cli for a quiz with 3 answears per questiona and a couple of questions. Time from start to finish is also noted.",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}
//Get method for api request to get the questions.
func GetQuestion(questionId string) {

	resp, err := http.Get("http://localhost:9090/questions/" + questionId)

	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	log.Printf(sb)
}

//post method to post the name, score & time 
func Answer(answer [6]string, name string, time float64)int {

	var score = 0;
	//Not my proudest moment
	if answer[0] == "3" {
		score++
		fmt.Printf("Question 1 was correct \n")
	} else {
		fmt.Printf("Question 1 was inccorrect \n")
	}
	if answer[1] == "3" {
		score++
		fmt.Printf("Question 2 was correct \n")
	} else {
		fmt.Printf("Question 2 was inccorrect \n")
	}
	if answer[2] == "2" {
		score++
		fmt.Printf("Question 3 was correct \n")
	} else {
		fmt.Printf("Question 3 was inccorrect \n")
	}
	if answer[3] == "1" {
		score++
		fmt.Printf("Question 4 was correct \n")
	} else {
		fmt.Printf("Question 4 was inccorrect \n")
	}
	if answer[4] == "2" {
		score++
		fmt.Printf("Question 5 was correct \n")
	} else {
		fmt.Printf("Question 5 was inccorrect \n")
	}
	if answer[5] == "1" {
		score++
		fmt.Printf("Question 5 was correct \n")
	} else {
		fmt.Printf("Question 5 was inccorrect \n")
	}
	//--
	payload, err := json.Marshal(map[string]interface{}{
		"name":  name,
		"score": score,
		"time": time,
	})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	url := "http://localhost:9090/answers"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	return score

}
// get the scores and stores them in struct below for easier handling.
func GetScores() []Score{

	resp, err := http.Get("http://localhost:9090/answers/")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var newScore []Score

	if err := json.Unmarshal(body, &newScore); err != nil {
		fmt.Println(err)
	}
	
	return newScore
}

func ShowScore() {

	displayScore := GetScores()

	//Sorting by the highest score
	sort.Slice(displayScore, func(a, b int) bool {
		return displayScore[a].Score > displayScore[b].Score
	})
		
	var UsersScores [][]string
	for _, s := range displayScore {
		var scores []string
		t := fmt.Sprintf("%f", s.Time)
		scores = append(scores, s.Name, strconv.Itoa(s.Score), t)
		UsersScores = append(UsersScores, scores)
	}
		//Tabelwriter : https://github.com/olekukonko/tablewriter
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Score", "Time(Seconds)"})
		table.SetBorder(false)
		table.AppendBulk(UsersScores)
		table.Render()
}

func Performance(score int)string {

	Scores := GetScores()
	var nrOfpeopleBeaten int
	//Checks for the amount of people that you have a better score then.
	for _, s := range Scores{
		if s.Score < score {
			nrOfpeopleBeaten ++
		}
	}

	result := strconv.Itoa(nrOfpeopleBeaten)

	return result
}

type Score struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Time float64  `json:"time"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Time to check ur knowledge, please enter a name:\n")
	var name string
	fmt.Scanln(&name)

	//Sets time quiz was started
	startTime := time.Now()

	var answer [6]string

	var nr = 0
	//Loop for all questions, if amount fo questions is increased the value 6 must increase 
	for i := 1; i <= 6; i++ {
		id := strconv.Itoa(i)
		GetQuestion(id)

		fmt.Printf("Please pick a answer 1, 2 or 3\n")

		var tempAnswer string
		fmt.Scanln(&tempAnswer)

		answer[nr] = tempAnswer
		nr++
	}

	//Finish time and then changes int to seconds
	finishTime := time.Since(startTime)
	TimeSeconds := finishTime.Seconds()

	//Returns score so i can use it below
	var score = Answer(answer, name, TimeSeconds)

	fmt.Printf("Thanks "+name+" that was all the questions, you have a better result than this many: people" + Performance(score))
	fmt.Printf("\nPress enter to see your result compared to others")
	var stop string
	fmt.Scanln(&stop)

	ShowScore()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fasttrackquiz.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
