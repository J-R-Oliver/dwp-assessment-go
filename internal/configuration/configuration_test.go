package configuration

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/J-R-Oliver/dwp-assessment-go/pkg/logging"
)

func TestLoadConfiguration(t *testing.T) {
	expected := Configuration{
		Port:         "8080",
		LoggingLevel: logging.Info,
		PeopleConfiguration: peopleConfiguration{
			BaseURL:  "https://dwp-techtest.herokuapp.com",
			Distance: 50,
		},
		Cities: map[string]City{
			"London": {
				Latitude:  "51.514248",
				Longitude: "-0.093145",
			},
		},
	}

	os.Setenv("PORT", "8080")
	os.Setenv("LOGGING_LEVEL", "info")
	os.Setenv("PEOPLE_ENDPOINT", "https://dwp-techtest.herokuapp.com")

	tests := []struct {
		name     string
		filename string
		want     Configuration
		wantErr  bool
	}{
		{
			"When configuration file contains no env vars then uses defaults",
			"./testdata/test-configuration.yaml",
			expected,
			false,
		},
		{
			"When configuration file contains only env vars then uses env vars",
			"./testdata/test-configuration-env-var.yaml",
			expected,
			false,
		},
		{
			"When configuration file contains mix of defaults and env vars then prioritises env vars",
			"./testdata/test-configuration-mixed.yaml",
			expected,
			false,
		},
		{
			"When configuration file does not exist then returns error",
			"./testdata/test-configuration-missing.yaml",
			Configuration{},
			true,
		},
		{
			"When configuration is invalid then returns error",
			"./testdata/test-configuration-invalid.yaml",
			Configuration{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfiguration(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Example_replaceSubstitutions_environment_variable_present() {
	o := replaceSubstitutions("8080:-9090")
	fmt.Println(o)
	// Output: 8080
}

func Example_replaceSubstitutions_environment_variable_not_present() {
	o := replaceSubstitutions(":-9090")
	fmt.Println(o)
	// Output: 9090
}

func Test_replaceSubstitutions(t *testing.T) {
	input := `port: 8080:-9090
logging-level: info:-debug

people:
  base-url: https://dwp-techtest.herokuapp.com:-https://dwp-techtest.herokuapp.com
  default-distance: :-50

cities:
  London:
    lat: 51.514248
    lon: -0.093145
`

	expected := `port: 8080
logging-level: info

people:
  base-url: https://dwp-techtest.herokuapp.com
  default-distance: 50

cities:
  London:
    lat: 51.514248
    lon: -0.093145
`
	if got := replaceSubstitutions(input); got != expected {
		t.Errorf("replaceSubstitutions() = %v, want %v", got, expected)
	}
}
