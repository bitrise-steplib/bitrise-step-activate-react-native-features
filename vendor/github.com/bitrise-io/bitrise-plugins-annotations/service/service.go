package service

import (
	"fmt"
	"net/http"
)

type Annotation struct {
	Context  string `json:"context"`
	Markdown string `json:"markdown"`
	Style    string `json:"style"`
}

func Annotate(annotation Annotation) error {
	_, err := do(http.MethodPut, "/", annotation)

	return err
}

func List() (string, error) {
	resp, err := do(http.MethodGet, "/", nil)

	return string(resp), err
}

func Get(context string) (string, error) {
	resp, err := do(http.MethodGet, fmt.Sprintf("/%s", context), nil)

	return string(resp), err
}

func Delete(context string) error {
	_, err := do(http.MethodDelete, fmt.Sprintf("/%s", context), nil)

	return err
}

func DeleteAll() error {
	_, err := do(http.MethodDelete, "/", nil)

	return err
}
