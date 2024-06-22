package src

type State struct {
	state map[string]interface{}
}

func (s *State) PutState(key string, value interface{}) {
	s.state[key] = value
}

func (s *State) GetState(key string) interface{} {
	return s.state[key]
}
