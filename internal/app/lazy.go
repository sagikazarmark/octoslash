package app

import "errors"

type LazyValue[T any] func() T

func (fn LazyValue[T]) Resolve() T {
	if fn == nil {
		panic("lazy value: no resolver")
	}

	return fn()
}

type LazyResult[T any] func() (T, error)

func (fn LazyResult[T]) Resolve() (T, error) {
	if fn == nil {
		var t T

		return t, errors.New("lazy result: no resolver")
	}

	return fn()
}

type LazyOption[T any] func() (T, bool)

func (fn LazyOption[T]) Resolve() (T, bool) {
	if fn == nil {
		var t T

		return t, false
	}

	return fn()
}
