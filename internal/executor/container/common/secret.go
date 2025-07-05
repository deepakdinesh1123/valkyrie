package common

import "fmt"

func ConvertSecretsMapToSlice(secrets map[string]string) []string {
	var resultSlice []string

	for key, value := range secrets {
		resultSlice = append(resultSlice, fmt.Sprintf("%s=%s", key, value))
	}

	return resultSlice
}
