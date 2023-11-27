package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type Crons struct{}

func (q *Crons) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Stack {
	case models.Laravel:
		answers.Crons = map[string]map[string]string{
			"scheduler": {
				"spec":      "*/5 * * * *",
				"cmd":      "php artisan schedule:run",
			},
		}
	}
	
	return nil
}
