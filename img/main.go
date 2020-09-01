package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	expect "github.com/Netflix/go-expect"
)

func main() {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	cmd := exec.Command("zsh")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	go func() {
		c.ExpectEOF()
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	commands := []string{
		"woke --help",
		"echo 'This should not have whitelist' > test.txt",
		"woke test.txt",
		"sed -i '' 's/whitelist/allowlist/g' test.txt",
		"woke test.txt",
	}
	time.Sleep(time.Second * 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, command := range commands {
			for _, char := range command {
				c.Send(fmt.Sprintf("%c", char))
				time.Sleep(time.Millisecond * 60)
			}
			time.Sleep(time.Millisecond * 500)
			c.Send("\n")
			time.Sleep(time.Millisecond * 1500)
		}
		time.Sleep(time.Millisecond * 2000)
	}()

	wg.Wait()
	_ = cmd.Process.Signal(syscall.SIGTERM)
}
