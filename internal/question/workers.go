package question

import (
	"context"

	"github.com/platformsh/platformify/internal/question/models"
)

type Workers struct{}

func (q *Workers) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Stack {
	case models.Laravel:
		answers.Workers = map[string]map[string]string{
			"queue": {
				"start":      "php artisan schedule:work",
			},
		}
	case models.Django:
		answers.Workers = map[string]map[string]string{
			"celery_worker": {
				"start":      "celery -A config.celery_app worker --loglevel=info",
			},
			"celery_beat": {
				"start":      "celery -A config.celery_app beat --loglevel=info",
			},
		}
	}

	return nil
}
