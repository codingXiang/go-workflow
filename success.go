package workflow

type SuccessFunc func(step *Step, context Context, callback func(objs ...interface{}) error) error

func SuccessCallback() SuccessFunc {
	return func(step *Step, context Context, callback func(objs ...interface{}) error) error {
		return callback()
	}
}
