package errors_test

import (
	"testing"

	"github.com/kapetndev/errors"
)

var (
	causeError      = &customError{"cause"}
	parentError     = &customError{"parent"}
	unexpectedError = &customError{"unexpected"}
)

type customError struct {
	Message string
}

func (e *customError) Error() string {
	return e.Message
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error", func(t *testing.T) {
		err := errors.New(causeError)
		if err == nil {
			t.Error("error was <nil>")
		}

		expectedErrorMessage := "cause"
		if err.Error() != expectedErrorMessage {
			t.Errorf("error messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.New(nil); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error with cause", func(t *testing.T) {
		err := errors.Wrap(parentError, causeError)
		if err == nil {
			t.Error("error was <nil>")
		}

		expectedErrorMessage := "parent: cause"
		if err.Error() != expectedErrorMessage {
			t.Errorf("error messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.Wrap(nil, causeError); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})
}

func TestNewWithMessage(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error with contextual message", func(t *testing.T) {
		err := errors.NewWithMessage(causeError, "error")
		if err == nil {
			t.Error("error was <nil>")
		}

		expectedErrorMessage := "error: cause"
		if err.Error() != expectedErrorMessage {
			t.Errorf("error messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.NewWithMessage(nil, "error"); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})
}

func TestWrapWithMessage(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error with cause and contextual message", func(t *testing.T) {
		err := errors.WrapWithMessage(parentError, causeError, "error")
		if err == nil {
			t.Error("error was <nil>")
		}

		expectedErrorMessage := "error: parent: cause"
		if err.Error() != expectedErrorMessage {
			t.Errorf("error messages are not equal: %s != %s", err, expectedErrorMessage)
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.WrapWithMessage(nil, causeError, "error"); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error", func(t *testing.T) {
		err := errors.Wrap(parentError, causeError)

		if err = errors.Unwrap(err); err != causeError {
			t.Errorf("errors are not equal: %s != %s", err, causeError)
		}
	})

	t.Run("returns nil when there are no more errors to unwrap", func(t *testing.T) {
		err := errors.Wrap(parentError, nil)

		if err = errors.Unwrap(err); err != nil {
			t.Errorf("error was not <nil>: %s", err)
		}
	})
}

func TestIs(t *testing.T) {
	t.Parallel()

	t.Run("returns true if an error matches target", func(t *testing.T) {
		err := errors.Wrap(parentError, causeError)

		if m := errors.Is(err, parentError); !m {
			t.Errorf("error was not expected type: %T != %T", err, parentError)
		}

		if m := errors.Is(err, causeError); !m {
			t.Errorf("error was not expected type: %T != %T", err, causeError)
		}
	})

	t.Run("returns false if an error cannot be found that matches target", func(t *testing.T) {
		err := errors.Wrap(parentError, causeError)

		if m := errors.Is(err, unexpectedError); m {
			t.Errorf("error was an unexpected type: %T", unexpectedError)
		}
	})
}

type missingError customError

func (e *missingError) Error() string {
	return e.Message
}

func TestAs(t *testing.T) {
	t.Parallel()

	t.Run("returns true if an error matches target, and sets target to that error value", func(t *testing.T) {
		err := errors.Wrap(parentError, causeError)

		var targetError *customError
		_ = errors.As(err, &targetError)

		// errors.As will set targetError to the first error in the chain that
		// matches the target type.
		if targetError != parentError {
			t.Errorf("error was not expected type: %T != %T", targetError, parentError)
		}

		// targetError should not be set to the causeError, as it is not of the
		// first error in the chain that matches the target type.
		if targetError == causeError {
			t.Errorf("error was an unexpected type: %T", targetError)
		}
	})

	t.Run("returns false if an error cannot be found that matches target", func(t *testing.T) {
		err := errors.Wrap(parentError, causeError)

		var targetError *missingError
		if _ = errors.As(err, &targetError); targetError != nil {
			t.Errorf("error was an unexpected type: %T", targetError)
		}
	})
}
