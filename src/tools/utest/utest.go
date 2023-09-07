package utest

import "testing"

type TestFunc[T any] func(*testing.T, T) bool
