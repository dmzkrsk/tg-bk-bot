package main

import "../bk"
import "fmt"
import "bufio"
import "os"

import "bytes"

func main() {
	reader := bufio.NewReader(os.Stdin)
	secretGenerator := bk.NewRandomSecretGenerator(4)
	session := bk.NewSession(secretGenerator)

	fmt.Println("Игровая сессия начата")

	for {
		line, err := reader.ReadBytes('\n')
		line = bytes.TrimSpace(line)

		if err != nil {
			fmt.Println(err) // todo stderr
			break
		}

		if string(line) == "/quit" {
			fmt.Println("Игровая сессия завершена")
			break
		}

		msg := bk.InMessage{Command: line}
		reply := session.HandleMessage(msg)

		fmt.Println(reply.Text)
	}
}
