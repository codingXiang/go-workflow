package workflow

type SuccessFunc func(step *Step, context Context, callback func() error) error

func SuccessCallback() SuccessFunc {
	return func(step *Step, context Context, callback func() error) error {
		return callback()
	}
}
