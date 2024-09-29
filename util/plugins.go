package util

type Plugin interface {
	load() error
	restart() error
	kill() error
}

type ChatApp struct {
	id string
}

func (ca *ChatApp) load() error {
	return nil
}
