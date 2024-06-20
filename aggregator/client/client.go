package client

import (
	"context"

	"github.com/KurobaneShin/tolling/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
