package workflow

import (
	"context"
	"fmt"
)

type Context interface{}

type Workflow struct {
	OnSuccess       SuccessFunc
	FailureCallback Failure
	OnFailure       FailureFunc
	Context         context.Context
	Start           *Step
	queue           []*Step
	inQueue         map[*Step]bool
}

func New() *Workflow {
	w := &Workflow{
		Context: context.Background(),
	}
	w.inQueue = make(map[*Step]bool)
	w.OnSuccess = SuccessCallback()
	w.FailureCallback = FailureCallback()
	return w
}

func (w *Workflow) Run(successCallback func(objs ...interface{}) error, failCallback func(objs ...interface{}) error) error {
	for _, step := range w.queue {
		if resp, err := step.Run(w.Context); err == nil {
			// step work fine
			if step.Hook != nil {
				if err := step.Hook(step, resp, err); err != nil {
					return err
				}
			}
		} else {
			// step run failed
			if err = w.OnFailure(err, step, w.Context); err != nil {
				fmt.Println("FAILED")
				if w.FailureCallback != nil {
					return w.FailureCallback(err, step, w.Context, failCallback)
				}
				if step.Hook != nil {
					if err = step.Hook(step, resp, err); err != nil {
						return err
					}
				}
			}
		}
	}
	return w.OnSuccess(nil, w.Context, successCallback)
}

func (w *Workflow) AddSteps(hook StepHook, steps ...*Step) {
	if w.queue == nil {
		w.queue = make([]*Step, 0)
	}
	for _, s := range steps {
		if hook != nil {
			s.Hook = hook
		}
		w.queue = append(w.queue, s)
	}
}

func (w *Workflow) loadStep(s *Step) {
	if s == nil {
		return
	}

	for _, step := range s.DependsOn {
		w.loadStep(step)
	}

	if !w.inQueue[s] {
		w.inQueue[s] = true
		w.queue = append(w.queue, s)
	}
	return
}
