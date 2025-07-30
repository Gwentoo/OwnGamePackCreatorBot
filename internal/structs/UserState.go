package structs

import "sync"

type UserState struct {
	mu       sync.Mutex
	state    string
	position [3]int
}

func NewUserState() UserState {
	return UserState{
		position: [3]int{-1, -1, -1},
		state:    "",
	}
}

func (us *UserState) SetState(state string) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.state = state
}

func (us *UserState) GetState() string {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.state
}

func (us *UserState) SetPos(pos int, value int) {
	us.mu.Lock()
	defer us.mu.Unlock()
	us.position[pos] = value
}

func (us *UserState) GetPos(pos int) int {
	us.mu.Lock()
	defer us.mu.Unlock()
	return us.position[pos]
}
