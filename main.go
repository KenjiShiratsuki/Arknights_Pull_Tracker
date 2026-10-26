package main

const maxHistoryLength = 10

type State struct {
	threeStarCount int
	fourStarCount  int
	fiveStarCount  int
	sixStarCount   int
	pityCount      int
	totalPullCount int
	threeStarNew   int
	fourStarNew    int
	fiveStarNew    int
	sixStarNew     int
	totalNew       int
}

type Log struct {
	history      []State
	currentState int
}

func main() {
	log := initializeState()
	a := initializeApp(log)
	a.Run()
}
