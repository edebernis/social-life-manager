package utils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var (
	// MockContextMatcher matches every context object
	MockContextMatcher = mock.MatchedBy(func(ctx context.Context) bool { return true })
)
