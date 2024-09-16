package task

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

var stateTtansiitonMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Running, Failed},
	Running:   {Completed, Failed},
	Completed: {},
	Failed:    {},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func ValidStateTransition(src State, dst State) bool {
	if src == dst {
		return true
	}
	return Contains(stateTtansiitonMap[src], dst)
}
