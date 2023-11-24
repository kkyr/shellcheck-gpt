package main

import (
	"fmt"
	"strings"
)

type Model int

const (
	GPT35Turbo Model = iota
	GPT4Turbo
)

var modelNames = map[Model]string{
	GPT35Turbo: "gpt-3.5-turbo",
	GPT4Turbo:  "gpt-4-turbo",
}

// String returns a string representation of the model.
func (m Model) String() string {
	return modelNames[m]
}

// Set satisfies the flag.Value interface and sets a Model based on the input string,
// returning an error for illegal values.
func (m *Model) Set(s string) error {
	for k, v := range modelNames {
		if v == s {
			*m = k
			return nil
		}
	}

	return fmt.Errorf("invalid model '%s', valid models are %s", s, strings.Join(validModels(), " or "))
}

func validModels() []string {
	models := make([]string, 0, len(modelNames))
	for _, v := range modelNames {
		models = append(models, v)
	}

	return models
}
