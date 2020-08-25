package k8s

import (
	"regexp"

	"github.com/spf13/afero"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"
)

type JobParser interface {
	Parse(filename string) (*batchv1.Job, error)
}

type jobParser struct {
	JobParser
	Fs afero.Fs
}

func NewJobParser(fs afero.Fs) JobParser {
	return &jobParser{Fs: fs}
}

func (p *jobParser) Parse(filename string) (*batchv1.Job, error) {
	b, err := afero.ReadFile(p.Fs, filename)
	if err != nil {
		return nil, err
	}

	var job *batchv1.Job
	re := regexp.MustCompile(`apiVersion:`)
	if re.Match(b) {
		job = &batchv1.Job{}
		if err := yaml.UnmarshalStrict(b, job); err != nil {
			return nil, err
		}
	} else {
		simple := &SimpleJob{}
		if err := yaml.UnmarshalStrict(b, simple); err != nil {
			return nil, err
		}
		job = simple.Expand()
	}

	return job, nil
}
