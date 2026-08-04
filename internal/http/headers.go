package http

import "strings"

type Headers struct {
	values map[string][]string
}

func NewHeaders() *Headers {
	h := Headers{}
	h.values = make(map[string][]string)
	return &h
}

func (h *Headers) Get(name string) string {
	if val, ok := h.values[strings.ToLower(name)]; ok {
		return val[0]
	}
	return ""
}

func (h *Headers) Values(name string) []string {
	if val, ok := h.values[strings.ToLower(name)]; ok {
		newVal := make([]string, len(val))
		copy(newVal, val)
		return newVal
	}
	return nil
}

func (h *Headers) Has(name string) bool {
	_, ok := h.values[strings.ToLower(name)]
	return ok
}

func (h *Headers) add(name, value string) {
	name = strings.ToLower(name)
	h.values[name] = append(h.values[name], value)
}
