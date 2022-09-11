package container

func Must[T any, E error](value T, err E) T {
	return Resuify(value, err).Value()
}
