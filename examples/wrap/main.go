package main

import (
	"fmt"

	"github.com/kapetndev/errors"
)

const (
	causeError       = customError("cause")
	parentError      = customError("parent")
	grandparentError = customError("grandparent")
)

type customError string

func (e customError) Error() string {
	return string(e)
}

func main() {
	err := generateError()
	fmt.Println(err)

	if errors.Is(err, causeError) {
		fmt.Println(causeError)
	}

	if errors.Is(err, parentError) {
		fmt.Println(parentError)
	}

	if errors.Is(err, grandparentError) {
		fmt.Println(grandparentError)
	}
}

func generateError() error {
	if err := generateParentError(); err != nil {
		return errors.Wrap(grandparentError, err)
	}

	return nil
}

func generateParentError() error {
	if err := generateCauseError(); err != nil {
		return errors.Wrap(parentError, err)
	}

	return nil
}

func generateCauseError() error {
	return errors.NewWithMessage(causeError, "something bad happened")
}
