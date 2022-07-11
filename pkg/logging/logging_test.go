package logging

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{
			"When level is Error then returns error",
			Error,
			"error",
		},
		{
			"When level is Info then returns info",
			Info,
			"info",
		},
		{
			"When level is Debug then returns debug",
			Debug,
			"debug",
		},
		{
			"When level is resolved value > 2 then returns nil string",
			Level(3),
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	if got := New(Info).(*logger); got.level != 1 {
		t.Errorf("New() = %v, want info", got.level)
	}
}

func Test_logger_Error(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stderr = w

	l := New(Error)
	l.Error("test error message")

	w.Close()

	out, _ := ioutil.ReadAll(r)
	s := string(out)

	if !strings.Contains(s, "test error message") {
		t.Errorf("Error() = %v, want test error message", s)
	}
}

func Test_logger_Info(t *testing.T) {
	t.Run("When logger level is Info then log message is printed to Stout", func(t *testing.T) {
		r, w, _ := os.Pipe()
		os.Stdout = w

		l := New(Info)
		l.Info("test info message")

		w.Close()

		out, _ := ioutil.ReadAll(r)
		s := string(out)

		if !strings.Contains(s, "test info message") {
			t.Errorf("Info() = %v, want test info message", s)
		}
	})

	t.Run("When log level is Error then no message is printed to Stout", func(t *testing.T) {
		r, w, _ := os.Pipe()
		os.Stdout = w

		l := New(Error)
		l.Info("test info message")

		w.Close()

		out, _ := ioutil.ReadAll(r)
		s := string(out)

		if s != "" {
			t.Errorf("Info() = %v, want \"\"", s)
		}
	})
}

func Test_logger_Debug(t *testing.T) {
	t.Run("When logger level is Debug then log message is printed to Stout", func(t *testing.T) {
		r, w, _ := os.Pipe()
		os.Stdout = w

		l := New(Debug)
		l.Debug("test debug message")

		w.Close()

		out, _ := ioutil.ReadAll(r)
		s := string(out)

		if !strings.Contains(s, "test debug message") {
			t.Errorf("Debug() = %v, want test debug message", s)
		}
	})

	t.Run("When log level is Info then no message is printed to Stout", func(t *testing.T) {
		r, w, _ := os.Pipe()
		os.Stdout = w

		l := New(Error)
		l.Debug("test debug message")

		w.Close()

		out, _ := ioutil.ReadAll(r)
		s := string(out)

		if s != "" {
			t.Errorf("Info() = %v, want \"\"", s)
		}
	})
}
