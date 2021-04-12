package workflow

import (
	"time"
)

type StepFunc func(context Context) (map[string]interface{}, error)
type StepHook func(step *Step, data map[string]interface{}, err error) error

type Step struct {
	Timeout   time.Duration
	Label     string
	Run       StepFunc
	Hook      StepHook
	DependsOn []*Step
}
