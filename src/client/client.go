package client

import (
	"errors"
	"fmt"
	"github.com/deiwin/interact"
	"log"
	"os"
	"strconv"

	//ui "github.com/gizak/termui" // <- ui shortcut, optional
)

/*
	Pseudo terminal
	Draw squares with characterss
	Events
	- up
	- down
	- left
	- right
	- characters
	- enter

	notebook
	commands (Assignment?)

	widgets
	clicks
	focus
	opengl (cross platform)

	actions (need for loop for bootstrapping self containment)

*/

var (
	checkNotEmpty = func(input string) error {
		// note that the inputs provided to these checks are already trimmed
		if input == "" {
			return errors.New("Input should not be empty!")
		}
		return nil
	}
	checkIsAPositiveNumber = func(input string) error {
		if n, err := strconv.Atoi(input); err != nil {
			return err
		} else if n < 0 {
			return errors.New("The number can not be negative!")
		}
		return nil
	}
)

func Run() {
	log.Printf("Works")

	actor := interact.NewActor(os.Stdin, os.Stdout)

	message := "type help for help"
	notEmpty, err := actor.Prompt(message, checkNotEmpty)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("msg= %s\n", notEmpty)

	message = "Please enter a positive number"
	n1, err := actor.PromptAndRetry(message, checkNotEmpty, checkIsAPositiveNumber)
	if err != nil {
		log.Fatal(err)
	}
	message = "Please enter another positive number"
	n2, err := actor.PromptOptionalAndRetry(message, "7", checkNotEmpty, checkIsAPositiveNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Thanks! (%s, %s, %s)\n", notEmpty, n1, n2)

}
