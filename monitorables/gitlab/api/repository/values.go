package repository

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func valueToString(values url.Values, key string) *string {
	if value := values.Get(key); value != "" {
		return &value
	}
	return nil
}

func valueToInt(values url.Values, key string) (*int, error) {
	if value := values.Get(key); value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %q into %s (int)", value, key)
		}
		return &i, nil
	}
	return nil, nil
}

func valueToLabels(values url.Values, key string) gitlab.Labels {
	if value := values.Get(key); value != "" {
		return strings.Split(value, ",")
	}
	return nil
}
