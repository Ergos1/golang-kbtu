package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"time"
)

const (
	timeLimit = 10
	maxVal    = 10
	minVal    = 1
)

var (
	operations  = [3]string{"*", "+", "-"}
	user_answer = 0
	score       = 0
	in          = make(chan int)
	username    = ""
)

type Question struct {
	x, y      int
	operation string
}

func (q Question) String() string {
	return fmt.Sprintf("%d %s %d", q.x, q.operation, q.y)
}

type User struct {
	Name  string
	Score int
}

func (u User) String() string {
	return fmt.Sprintf("%s %d", u.Name, u.Score)
}

type ByScore []*User

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

func swap(x, y *int) {
	*x = *x + *y // a = a + b
	*y = *x - *y // a + b - b
	*x = *x - *y // a + b - a
}

func generateNewQuetion() Question {
	question := Question{}
	question.x = rand.Intn(maxVal-minVal) + minVal
	question.y = rand.Intn(maxVal-minVal) + minVal
	if question.x < question.y {
		swap(&question.x, &question.y)
	}
	question.operation = operations[rand.Intn(len(operations))]
	return question
}

func getAnswer(q *Question) (int, error) {
	switch q.operation {
	case "+":
		return q.x + q.y, nil
	case "-":
		return q.x - q.y, nil
	case "*":
		return q.x * q.y, nil
	default:
		return -1, fmt.Errorf("Invalid operation")
	}
}

func checkAnswer(question *Question, userAnswer int) (bool, error) {
	questionAnswer, err := getAnswer(question)
	if err != nil {
		return false, err
	}
	if questionAnswer == userAnswer {
		return true, nil
	}
	return false, nil
}

func startGame() {
	timer := time.After(timeLimit * time.Second)
GAME:
	for {
		question := generateNewQuetion()

		fmt.Print(question, " = ")
		go func() {
			fmt.Scanf("%d\n", &user_answer)
			in <- user_answer
		}()

		select {
		case <-in:
			if ok, _ := checkAnswer(&question, user_answer); ok {
				score += 1
			}
			continue GAME
		case <-timer:
			break GAME
		}
	}
	fmt.Printf("\nYou(%s) scored: %d", username, score)
	user := &User{username, score}
	addToScoreBoard(user)
}

func readScoreBoard() {
	file, err := ioutil.ReadFile("scoreboard.json")
	if err != nil {
		ioutil.WriteFile("scoreboard.json", []byte("[]"), 0755)
	}
	users := []*User{}
	json.Unmarshal(file, &users)
	for _, user := range users {
		fmt.Println(user)
	}
}

func addToScoreBoard(user *User) {
	file, err := ioutil.ReadFile("scoreboard.json")
	if err != nil {
		ioutil.WriteFile("scoreboard.json", []byte("[]"), 0755)
	}
	users := []*User{}
	json.Unmarshal(file, &users)
	users = append(users, user)
	sort.Sort(ByScore(users))
	js, _ := json.Marshal(users)
	ioutil.WriteFile("scoreboard.json", js, 0755)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer close(in)

	fmt.Print("Enter your username: ")
	fmt.Scanf("%s\n", &username)

	startGame()

	fmt.Println("\nSCOREBOARD:")
	readScoreBoard()
	fmt.Scan()
}
