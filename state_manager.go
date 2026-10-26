package main

func pushState(newState State, log *Log) {
	log.history = log.history[:log.currentState+1]
	if len(log.history) == maxHistoryLength {
		copy(log.history, log.history[:1])
		log.history[maxHistoryLength-1] = newState
	} else {
		log.history = append(log.history, newState)
	}
	log.currentState = len(log.history) - 1
}
