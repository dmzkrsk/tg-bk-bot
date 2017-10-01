package bk

import (
	"fmt"
	"regexp"
)

const (
	helpMessage    = "Тут будет справка"
	unknownCommand = "Неизвестная команда"
)

type InMessage struct {
	Uid     string
	Command []byte
}

func (m InMessage) is(s string) bool {
	return string(m.Command) == s
}

func (m InMessage) isCommand() bool {
	return m.Command[0] == '/'
}

func (m InMessage) match(r *regexp.Regexp) bool {
	return r.Match(m.Command)
}

type OutMessage struct {
	Uid  string
	Text string
}

type stateHandler interface {
	HandleStart(SecretGenerator) (string, stateHandler)
	HandleGuess([]byte) (string, stateHandler)
	HandleGiveup() (string, stateHandler)
	HandleDefault() (string, stateHandler)
}

type inactiveState struct{}
type activeState struct {
	*Secret
	n uint
}

func (st inactiveState) HandleStart(g SecretGenerator) (string, stateHandler) {
	activeState := activeState{Secret: g.CreateSecret()}
	return fmt.Sprintf("Я задумал новое число из %d цифр. Угадай его!", activeState.Size), &activeState
}

func (st inactiveState) HandleGuess(guess []byte) (string, stateHandler) {
	return "Игра еще не началась. Введите /start для запуска новой игры", st
}

func (st inactiveState) HandleGiveup() (string, stateHandler) {
	return "Сдаётесь даже не попробовав? Введите /start для запуска новой игры", st
}

func (st inactiveState) HandleDefault() (string, stateHandler) {
	return unknownCommand, st
}

func (st *activeState) HandleStart(_ SecretGenerator) (string, stateHandler) {
	return "Игра уже началась. Введите /giveup, если хотите сдаться и попробовать ещё раз", st
}

func (st *activeState) HandleGuess(guess []byte) (string, stateHandler) {
	st.n++

	if uint8(len(guess)) != st.Size {
		return fmt.Sprintf("Я загадал число из %d цифр", st.Size), st
	}

	match := st.Guess(guess)

	if match.IsWin() {
		if st.n == 1 {
			return "Вы угадали число с первой попытки! Мы точно не знакомы?", inactiveState{}
		}

		return fmt.Sprintf("Браво! Вы угадали число за %d попыток", st.n), inactiveState{}
	}

	return fmt.Sprintf("Неправильно. Держи подсказку: %s", match), st
}

func (st *activeState) HandleGiveup() (string, stateHandler) {
	var template string

	if st.n == 0 {
		template = "Даже не попробовали? Было задумано число %v. Введите /start для начала новой игры"
	} else {
		template = "Я победил! Было задумано число %v. Введите /start для начала новой игры"
	}

	return fmt.Sprintf(template, st.Secret), inactiveState{}
}

func (st *activeState) HandleDefault() (string, stateHandler) {
	return "Вводите только цифры", st
}

type Session struct {
	currentState stateHandler
	g            SecretGenerator
}

func NewSession(g SecretGenerator) Session {
	return Session{inactiveState{}, g}
}

var onlyNumbers = regexp.MustCompile(`^\d+$`)

func (s *Session) HandleMessage(msg InMessage) OutMessage {
	reply := unknownCommand

	switch true {
	case msg.is("/start"):
		reply, s.currentState = s.currentState.HandleStart(s.g)
	case msg.is("/giveup"):
		reply, s.currentState = s.currentState.HandleGiveup()
	case msg.is("/help"):
		reply = helpMessage
	case msg.isCommand():
		reply = unknownCommand
	case msg.match(onlyNumbers):
		reply, s.currentState = s.currentState.HandleGuess(msg.Command)
	default:
		reply, s.currentState = s.currentState.HandleDefault()
	}

	return OutMessage{msg.Uid, reply}
}
