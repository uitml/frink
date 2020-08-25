package k8s

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestParseValidJobSpec(t *testing.T) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := NewJobParser(fs)

	job, err := parser.Parse("job/basic.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "foo", job.Name)
}

func TestParseInvalidJobSpec(t *testing.T) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := NewJobParser(fs)

	job, err := parser.Parse("job/invalid.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown field \"naem\"")
	assert.Nil(t, job)
}

func TestParseValidSimpleJobSpec(t *testing.T) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := NewJobParser(fs)

	job, err := parser.Parse("simplejob/basic.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "foo", job.Name)
}

func TestParseInvalidSimpleJobSpec(t *testing.T) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := NewJobParser(fs)

	job, err := parser.Parse("simplejob/invalid.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown field \"naem\"")
	assert.Nil(t, job)
}

func TestParseMissingFile(t *testing.T) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), "testdata")
	parser := NewJobParser(fs)

	job, err := parser.Parse("missing.yaml")
	assert.Error(t, err)
	assert.IsType(t, &os.PathError{}, err)
	assert.Nil(t, job)
}
