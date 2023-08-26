package rate

import (
	"context"
	"identity-v2/internal/model"
)

type RateControl interface {
	Ban(context.Context, model.ID, int32) error
}
