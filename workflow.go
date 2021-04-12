package workflow

import (
	"context"
	"errors"
	"fmt"
)

type Context interface{}

type Workflow struct {
	OnSuccess       SuccessFunc
	FailureCallback Failure
	OnFailure       FailureFunc
	Context         context.Context
	Start           *Step
	channel         chan *Step
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
	w.loadQueue()
	for {
		select {
		case <-w.Context.Done():
			return errors.New("context failed")
		case step, ok := <-w.channel:
			if !ok {
				fmt.Println("COMPLETE")
				w.OnSuccess(step, w.Context, successCallback)
				return nil
			}
			resp, err := step.Run(w.Context)

			if err != nil {
				if err = w.OnFailure(err, step, w.Context); err != nil {
					fmt.Println("FAILED")
					if w.FailureCallback != nil {
						w.FailureCallback(err, step, w.Context, failCallback)
					}
					if step.Hook != nil {
						return step.Hook(resp, err)
					}
					return err
				}
			}

		}
	}
}

func (w *Workflow) AddSteps(steps ...*Step) {
	if w.queue == nil {
		w.queue = make([]*Step, 0)
	}
	for _, s := range steps {
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

func (w *Workflow) loadQueue() {
	w.loadStep(w.Start)
	w.channel = make(chan *Step, len(w.queue))
	for _, s := range w.queue {
		w.channel <- s
	}
	close(w.channel)
}
