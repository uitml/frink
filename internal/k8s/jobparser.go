package k8s

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"
)

type JobParser interface {
	Parse(filename string) (*batchv1.Job, error)
}

type jobParser struct {
	JobParser
}

func NewJobParser() JobParser {
	return &jobParser{}
}

func (p *jobParser) Parse(filename string) (*batchv1.Job, error) {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("specified file does not exist: %v", filename)
		}

		return nil, fmt.Errorf("unable to access file: %w", err)
	}

	return parseJob(filename)
}

func parseJob(file string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// TODO: Refactor this by extracting functions, etc.
	var job *batchv1.Job
	re := regexp.MustCompile(`apiVersion:`)
	if re.Match(data) {
		job = &batchv1.Job{}
		if err := yaml.UnmarshalStrict(data, job); err != nil {
			return nil, err
		}
	} else {
		simple := &SimpleJob{}
		if err := yaml.UnmarshalStrict(data, simple); err != nil {
			return nil, err
		}
		job = simple.Expand()
	}

	return job, nil
}
