package fake

import (
	"github.com/stretchr/testify/mock"
	batchv1 "k8s.io/api/batch/v1"
)

type JobParser struct {
	mock.Mock
}

func (p *JobParser) Parse(b []byte) (*batchv1.Job, error) {
	args := p.Called(b)
	job, _ := args.Get(0).(*batchv1.Job)

	return job, args.Error(1)
}
