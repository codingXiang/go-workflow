package workflow

import (
	"time"
)

type StepFunc func(context Context) error

type Step struct {
	Timeout   time.Duration
	Label     string
	Run       StepFunc
	DependsOn []*Step
}
