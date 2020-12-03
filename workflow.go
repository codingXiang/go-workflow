package workflow

import (
	"fmt"
)

type Context interface{}

type Workflow struct {
	Start           *Step
	OnSuccess       SuccessFunc
	FailureCallback Failure
	OnFailure       FailureFunc
	Context         Context

	queue   []*Step
	inQueue map[*Step]bool
}

func New() *Workflow {
	w := &Workflow{}
	w.inQueue = make(map[*Step]bool)
	w.OnSuccess = SuccessCallback()
	w.FailureCallback = FailureCallback()
	return w
}

func (w *Workflow) Run(successCallback func(objs ...interface{}) error, failCallback func(objs ...interface{}) error) error {
	for _, step := range w.queue {
		fmt.Printf("Running step: %s ", step.Label)
		if err := step.Run(w.Context); err != nil {
			if err := w.OnFailure(err, step, w.Context); err != nil {
				fmt.Println("FAILED")
				w.FailureCallback(err, step, w.Context, failCallback)
				return err
			}
		}
		fmt.Println("COMPLETE")
		w.OnSuccess(step, w.Context, successCallback)
	}
	return nil
}

func (w *Workflow) AddStep(s *Step) {
	if w.queue == nil {
		w.queue = make([]*Step, 0)
	}
	w.queue = append(w.queue, s)
}

func (w *Workflow) loadQueue(s *Step) {
	if s == nil {
		return
	}

	for _, step := range s.DependsOn {
		w.loadQueue(step)
	}

	if !w.inQueue[s] {
		w.inQueue[s] = true
		w.queue = append(w.queue, s)
	}
	return
}
