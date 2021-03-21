package errors_test

import (
	"testing"

	"github.com/kapetndev/errors"
)

var (
	causeError = &testError{"cause"}
	newError   = &testError{"new"}
	otherError = &missingError{"other"}
)

type testError struct {
	message string
}

func (e testError) Error() string {
	return e.message
}

type missingError testError

func (e missingError) Error() string {
	return e.message
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error", func(t *testing.T) {
		if err := errors.New(newError); err == nil {
			t.Error("expected: error, got: <nil>")
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.New(nil); err != nil {
			t.Errorf("expected: <nil>, got: %v", err)
		}
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error and cause", func(t *testing.T) {
		if err := errors.Wrap(newError, causeError); err == nil {
			t.Error("expected: error, got: <nil>")
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.Wrap(nil, causeError); err != nil {
			t.Errorf("expected: <nil>, got: %v", err)
		}
	})
}

func TestWrap_Error(t *testing.T) {
	t.Parallel()

	t.Run("returns the full error trace", func(t *testing.T) {
		err := errors.Wrap(newError, causeError)

		expected := "new: cause"
		if err.Error() != expected {
			t.Errorf("expected %s, got: %s", expected, err.Error())
		}
	})

	t.Run("returns the full error chase when there is no root cause", func(t *testing.T) {
		err := errors.New(newError)

		expected := "new"
		if err.Error() != expected {
			t.Errorf("expected %s, got: %s", expected, err.Error())
		}
	})
}

func TestNewWithMessage(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error and contextual message", func(t *testing.T) {
		if err := errors.NewWithMessage(newError, "error"); err == nil {
			t.Error("expected: error, got: <nil>")
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.NewWithMessage(nil, "error"); err != nil {
			t.Errorf("expected: <nil>, got: %v", err)
		}
	})
}

func TestWrapWithMessage(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error from a single error and cause and contextual message", func(t *testing.T) {
		if err := errors.WrapWithMessage(newError, causeError, "error"); err == nil {
			t.Error("expected: error, got: <nil>")
		}
	})

	t.Run("returns nil when the error is nil", func(t *testing.T) {
		if err := errors.WrapWithMessage(nil, causeError, "error"); err != nil {
			t.Errorf("expected: <nil>, got: %v", err)
		}
	})
}

func TestWrapWithMessage_Error(t *testing.T) {
	t.Parallel()

	t.Run("returns the full error trace", func(t *testing.T) {
		err := errors.WrapWithMessage(newError, causeError, "context")

		expected := "context: new: cause"
		if err.Error() != expected {
			t.Errorf("expected %s, got: %s", expected, err.Error())
		}
	})

	t.Run("returns the full error chase when there is no root cause", func(t *testing.T) {
		err := errors.NewWithMessage(newError, "context")

		expected := "context: new"
		if err.Error() != expected {
			t.Errorf("expected %s, got: %s", expected, err.Error())
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("returns a wrapped error", func(t *testing.T) {
		err := errors.Wrap(newError, causeError)

		if err = errors.Unwrap(err); err != causeError {
			t.Errorf("expected: %v, got: %v", causeError, err)
		}
	})

	t.Run("returns nil when there are no more errors to unwrap", func(t *testing.T) {
		err := errors.Wrap(newError, nil)

		if err = errors.Unwrap(err); err != nil {
			t.Errorf("expected: <nil>, got: %v", err)
		}
	})
}

func TestIs(t *testing.T) {
	t.Parallel()

	t.Run("returns true if an error matches target", func(t *testing.T) {
		err := errors.Wrap(newError, causeError)

		if m := errors.Is(err, newError); !m {
			t.Error("expected: true, got: false")
		}

		if m := errors.Is(err, causeError); !m {
			t.Error("expected: true, got: false")
		}
	})

	t.Run("returns false if an error cannot be found that matches target", func(t *testing.T) {
		err := errors.Wrap(newError, causeError)

		if m := errors.Is(err, otherError); m {
			t.Error("expected: false, got: true")
		}
	})
}

func TestAs(t *testing.T) {
	t.Parallel()

	t.Run("returns true if an error matches target, and sets target to that error value", func(t *testing.T) {
		test := func(err, expected error) {
			var targetError *testError
			if m := errors.As(err, &targetError); !m || targetError != expected {
				t.Errorf("expected: (true, %v), got: (%t, %v)", expected, m, targetError)
			}
		}

		test(errors.Wrap(newError, otherError), newError)
		test(errors.Wrap(otherError, causeError), causeError)
	})

	t.Run("returns false if an error cannot be found that matches target", func(t *testing.T) {
		err := errors.Wrap(newError, causeError)

		var target *missingError
		if m := errors.As(err, &target); m || target != nil {
			t.Errorf("expected: (false, <nil>), got: (%t, %v)", m, target)
		}
	})
}
