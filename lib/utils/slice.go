package utils

func TryMap[T, R any](a []T, mapper func(*T) (R, error)) ([]R, error) {
	var (
		err error
		out = make([]R, len(a))
	)

	for i, v := range a {
		out[i], err = mapper(&v)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func Map[T, R any](a []T, mapper func(*T) R) []R {
	out := make([]R, len(a))
	for i := range a {
		out[i] = mapper(&a[i])
	}

	return out
}

func TryForEach[T any](a []T, fn func(T) error) error {
	for _, v := range a {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

func ForEach[T any](a []T, fn func(T)) {
	for _, v := range a {
		fn(v)
	}
}
