package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
    fmt.Println("Hello, Go!")
    csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
    timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

    flag.Parse()

    file, err := os.Open(*csvFilename)
    if err != nil {
        exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
    }

    r := csv.NewReader(file)
    lines, err := r.ReadAll()
    if err != nil {
        exit("Failed to parse the provided CSV file.")
    }

    problems := parseLines(lines)
   
    timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

    problemloop:
        for i, p := range problems {
            fmt.Printf("Problem #%d: %s = ", i+1, p.question)
            answeCh := make(chan string) /// create a channel to receive the answer -> goroutine

            go func() { // anonymous function -> goroutine
                var answer string
                fmt.Scanf("%s\n", &answer)
                answeCh <- answer
            }()
            
            select {
            case <-timer.C:
                fmt.Println("Time's up!")
                break problemloop
            case answer := <- answeCh:
                if answer == p.answer {
                    fmt.Println("Correct!")
                } else {
                    fmt.Println("Wrong!")
                }
            }
        }  
    }

func parseLines(lines [][]string) []problem {
    ret := make([]problem, len(lines))
    for i, line := range lines {
        ret[i] = problem{
            question: line[0],
            answer: strings.TrimSpace(line[1]),
        }
    }
    return ret
}

func exit(s string) {
	fmt.Println(s)
    os.Exit(1)
}

type problem struct {
    question string
    answer   string
}