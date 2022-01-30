package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const TOTAL_TIME = 30

func main() {
	fmt.Println("Here we go. ")
	filename := flag.String("csvfile", "problems.csv", "The filename.extn of the csv file that has question-answer string tuples")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file with name %q due to err : %s", *filename, err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read records from the file %q", *filename))
	}

	quizItems := parseRecords(records)
	inputReader := bufio.NewReader(os.Stdin)
	timer := time.NewTimer(TOTAL_TIME * time.Second)

	var score int
	for i, item := range quizItems {
		fmt.Printf("\nQuestion #%d is: %s = ?\n", i+1, item.q)

		// Spin up goroutine to collect answer and pass it to main routine.
		ansCh := make(chan string)
		go func() {
			input, err := inputReader.ReadString('\n') // blocks until input
			if err != nil {
				exit(fmt.Sprintf("error reading user input: %v", err))
			}
			ansCh <- strings.TrimSpace(input)
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime's up!! You scored %d out of %d.\n", score, len(quizItems))
			return
		case answer := <-ansCh:
			if answer != strings.TrimSpace(item.a) {
				fmt.Printf("Wrong! The answer is %s.\n.", item.a)
			} else {
				score += 1
				fmt.Printf("Correct! Your score is now %d/%d.\n", score, len(quizItems))
			}
			// No default case.  Either timer messages or the goRoutine do
		}
	}
	fmt.Printf("Quiz over. Your final score is  %d/%d.\n", score, len(quizItems))
}

type QuizItem struct {
	q string
	a string
}

func parseRecords(records [][]string) []QuizItem {
	items := make([]QuizItem, len(records))
	for i, rec := range records {
		items[i] = QuizItem{
			q: rec[0],
			a: rec[1],
		}
	}
	return items
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
