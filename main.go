package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	questionFile = flag.String("f", "problems.csv", "A csv file containing questions and their solutions")
	timeLimit    = flag.Int("t", 30, "the time limit for the quiz in seconds")
)

type problem struct {
	question string
	answer   string
}

func main() {
	flag.Parse()

	// check file type
	if !strings.HasSuffix(*questionFile, "csv") {
		fmt.Printf("Question file '%s' is not a csv file", *questionFile)
	}
	// read csv file
	file, err := os.Open(*questionFile)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *questionFile))
	}
	// read and parse file
	data := csv.NewReader(file)
	lines, err := data.ReadAll()
	if err != nil {
		exit("Failed to parse a csv file: %s \n")
	}
	problems := parseLines(lines)

	correct := run(problems, *timeLimit)
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			// edge case: space between question and answer
			answer: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func run(problems []problem, timeLimit int) (correct int) {
	// create timer
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	correct = 0
	//iterative problem records
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.question)
		// run a goroutine to scan answers from users
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			// send user answer to the user answerChannel
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Time is out")
			fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
			// only break the select not for loop
			return correct
		case answer := <-answerChannel:
			if answer == p.answer {
				correct++
			}
		}
	}
	return correct
}
