package todo

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrEmptyTitle = errors.New("title is required")