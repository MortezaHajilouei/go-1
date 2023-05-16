package validator

type Stage[T any] func(T) error

func NewPipeline[T any](stages ...Stage[T]) Stage[T] {
	return func(t T) error {
		for _, f := range stages {
			if err := f(t); err != nil {
				return err
			}
		}
		return nil
	}
}

func (s Stage[T]) Validate(t T) error {
	return s(t)
}
