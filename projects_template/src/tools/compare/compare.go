package compare

func Contains[T []K, K comparable](s T, contains K) bool {
	for _, item := range s {
		if item == contains {
			return true
		}
	}

	return false
}

func RemoveDuplicate[T comparable](list []T) (newList []T) {
	k := make(map[T]bool)

	for _, item := range list {
		if _, value := k[item]; !value {
			k[item] = true
			newList = append(newList, item)
		}
	}

	return
}

type Constructor[T, K any] func(t T) K

func Rebuild[T, K any](from []T, newK Constructor[T, K]) (to []K) {
	to = make([]K, 0, len(from))
	for _, item := range from {
		to = append(to, newK(item))
	}
	return
}
