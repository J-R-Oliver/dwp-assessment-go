package logging

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
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

func TestLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  Level
		err   error
	}{
		{
			"When passed valid Level string as byte array then sets level",
			[]byte("info"),
			Info,
			nil,
		},
		{
			"When passed invalid Level string as byte array then returns error",
			[]byte("Not a valid level"),
			0,
			errors.New("Level.UnmarshalJSON: failed to unmarshal: Not a valid level is not a valid log level - valid options are error, info or debug"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l Level

			err := l.UnmarshalJSON(tt.bytes)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Level.UnmarshalJSON() err = %v, want %v", err, tt.err)
			}

			if l != tt.want {
				t.Errorf("Level.UnmarshalJSON() = %v, want %v", l, tt.want)
			}
		})
	}
}

func TestLevel_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		node *yaml.Node
		want Level
		err  error
	}{
		{
			"When passed valid Level string as *yaml.Node then sets level",
			&yaml.Node{Value: "info"},
			Info,
			nil,
		},
		{
			"When passed invalid Level string as *yaml.Node then returns error",
			&yaml.Node{Value: "Not a valid level"},
			0,
			errors.New("Level.UnmarshalYAML: failed to unmarshal: Not a valid level is not a valid log level - valid options are error, info or debug"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l Level

			err := l.UnmarshalYAML(tt.node)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Level.UnmarshalJSON() err = %v, want %v", err, tt.err)
			}

			if l != tt.want {
				t.Errorf("Level.UnmarshalJSON() = %v, want %v", l, tt.want)
			}
		})
	}
}

func Test_stringToLevel(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name    string
		args    args
		want    Level
		wantErr bool
	}{
		{
			"When passed error then returns Error level",
			args{"error"},
			Error,
			false,
		},
		{
			"When passed info then returns Info level",
			args{"info"},
			Info,
			false,
		},
		{
			"When passed debug then returns Debug level",
			args{"debug"},
			Debug,
			false,
		},
		{
			"When passed invalid level then returns error",
			args{"Not a valid level"},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToLevel(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("stringToLevel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	if got := New(Info).(*logger); got.level != 1 {
		t.Errorf("New() = %v, want info", got.level)
	}
}

func captureLogs(writeToLogsFunc func()) string {
	r, w, _ := os.Pipe()

	os.Stdout = w
	os.Stderr = w

	writeToLogsFunc()

	w.Close()

	s, _ := ioutil.ReadAll(r)

	return string(s)
}

func Test_logger_Error(t *testing.T) {
	s := captureLogs(func() {
		l := New(Error)
		l.Error("test error message")
	})

	if !strings.Contains(s, "test error message") {
		t.Errorf("Error() = %v, want test error message", s)
	}
}

func Test_logger_Info(t *testing.T) {
	t.Run("When logger level is Info then log message is printed to Stout", func(t *testing.T) {
		s := captureLogs(func() {
			l := New(Info)
			l.Info("test info message")
		})

		if !strings.Contains(s, "test info message") {
			t.Errorf("Info() = %v, want test info message", s)
		}
	})

	t.Run("When log level is Error then no message is printed to Stout", func(t *testing.T) {
		s := captureLogs(func() {
			l := New(Error)
			l.Info("test info message")
		})

		if s != "" {
			t.Errorf("Info() = %v, want \"\"", s)
		}
	})
}

func Test_logger_Debug(t *testing.T) {
	t.Run("When logger level is Debug then log message is printed to Stout", func(t *testing.T) {
		s := captureLogs(func() {
			l := New(Debug)
			l.Debug("test debug message")
		})

		if !strings.Contains(s, "test debug message") {
			t.Errorf("Debug() = %v, want test debug message", s)
		}
	})

	t.Run("When log level is Info then no message is printed to Stout", func(t *testing.T) {
		s := captureLogs(func() {
			l := New(Error)
			l.Debug("test debug message")
		})

		if s != "" {
			t.Errorf("Info() = %v, want \"\"", s)
		}
	})
}
