package validator

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsUUIDValid[T any](f func(T) string, message string) func(T) error {
	return func(t T) error {
		str := f(t)
		if _, err := uuid.Parse(str); err != nil {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func IsDigit[T any](f func(T) string, message string) func(T) error {
	return func(t T) error {
		str := f(t)
		if _, err := strconv.Atoi(str); err != nil {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func EqualTo[T any, K comparable](f func(T) K, equalTo K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		if field != equalTo {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func NotEqualTo[T any, K comparable](f func(T) K, notEqualTo K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		if field == notEqualTo {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func Greater[T any, K constraints.Ordered](f func(T) K, compareTo K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		fmt.Println(field)
		if field <= compareTo {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func GreaterOrEqual[T any, K constraints.Ordered](f func(T) K, compareTo K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		if field >= compareTo {
			return nil
		}
		return status.New(codes.InvalidArgument, message).Err()
	}
}

func Less[T any, K constraints.Ordered](f func(T) K, compareTo K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		if field < compareTo {
			return nil
		}
		return status.New(codes.InvalidArgument, message).Err()
	}
}

func LessOrEqual[T any, K constraints.Ordered](f func(T) K, compareTo K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		if field > compareTo {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func Between[T any, K constraints.Ordered](f func(T) K, from K, to K, message string) func(T) error {
	return func(t T) error {
		field := f(t)
		if field < from || field > to {
			return status.New(codes.InvalidArgument, message).Err()
		}
		return nil
	}
}

func CheckAll[T any, K any](f func(K) []T, s Stage[T]) Stage[K] {
	return func(k K) error {
		list := f(k)
		for i := range list {
			if err := s(list[i]); err != nil {
				return err
			}
		}
		return nil
	}
}
