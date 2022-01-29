package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

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
	var score int
	for i, item := range quizItems {
		fmt.Printf("\nQuestion #%d is: %s = ?\n", i+1, item.q)
		answer := getInput(inputReader)
		if answer != strings.TrimSpace(item.a) {
			fmt.Printf("Wrong! The answer is %s.\n.", item.a)
		} else {
			score += 1
			fmt.Printf("Correct! Your score is now %d/%d.\n", score, len(quizItems))
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

func getInput(reader *bufio.Reader) string {
	input, err := reader.ReadString('\n')

	if err != nil {
		exit(fmt.Sprintf("error reading user input: %v", err))
	}

	return strings.TrimSpace(input)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
