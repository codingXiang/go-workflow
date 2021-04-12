package workflow_test

import (
	"errors"
	"github.com/codingXiang/go-workflow"
	"strings"
	"testing"
)

func TestFailureFunc(t *testing.T) {
	var testVar bool

	w := workflow.New()
	w.OnFailure = func(err error, step *workflow.Step, context workflow.Context) error {
		testVar = true
		return nil
	}
	w.Start = &workflow.Step{
		Label: "fail workflow",
		Run: func(c workflow.Context) (map[string]interface{}, error) {
			return nil, errors.New("generic error")
		},
	}

	err := w.Run(nil, nil)
	if err != nil {
		t.Error(err)
	}
	if testVar != true {
		t.Fail()
	}
}

func TestInteractiveFailure(t *testing.T) {
	var testVar bool

	workflow.InputFile = strings.NewReader("s\n")

	w := workflow.New()
	w.OnFailure = workflow.InteractiveFailure
	w.AddSteps(
		nil,
		&workflow.Step{
			Label: "fail workflow",
			Run: func(c workflow.Context) (map[string]interface{}, error) {
				return nil, errors.New("generic error")
			},
		},
		&workflow.Step{
			Label: "succeed workflow",
			Run: func(c workflow.Context) (map[string]interface{}, error) {
				testVar = true
				return nil, nil
			},
		},
	)

	err := w.Run(nil, nil)
	if err != nil {
		t.Error(err)
	}
	if testVar != true {
		t.Fail()
	}
}
