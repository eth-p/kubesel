package parallel

import (
	"os"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type simulatedDoc struct {
	Foo *string `yaml:"foo"`
	Baz struct {
		Nested struct {
			Value *int `yaml:"value"`
		} `yaml:"nested"`
	} `yaml:"baz"`
}

func (d *simulatedDoc) Merge(other *simulatedDoc) *simulatedDoc {
	var merged simulatedDoc
	merged.Foo = d.Foo
	merged.Baz.Nested.Value = other.Baz.Nested.Value
	return &merged
}

var simulatedYAML = []byte(dedent.Dedent(`
foo: bar
baz:
  nested:
    value: 3
`))

func simulateReadYAML(t testing.TB) *simulatedDoc {

	// Pretend we read the file.
	file, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		require.NoError(t, err, "opening /dev/null")
	}

	_, _ = file.Write(simulatedYAML)
	_ = file.Close()

	// Unmarshal it.
	var data simulatedDoc
	err = yaml.Unmarshal(simulatedYAML, &data)
	if err != nil {
		require.NoError(t, err, "unmarshalling data")
	}

	return &data
}
