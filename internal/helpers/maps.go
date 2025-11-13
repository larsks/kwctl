package helpers

func ReverseMap[K1 comparable, V1 comparable](inputMap map[K1]V1) map[V1]K1 {
	reversedMap := make(map[V1]K1, len(inputMap))
	for k, v := range inputMap {
		reversedMap[v] = k
	}
	return reversedMap
}
