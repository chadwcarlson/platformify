package question

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"golang.org/x/exp/slices"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type DeployCommand struct{}

func (q *DeployCommand) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	switch answers.Stack {
	case models.Django:
		managePyPath := utils.FindFile(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), managePyFile)
		if managePyPath != "" {
			managePyPath, _ = filepath.Rel(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), managePyPath)
			prefix := ""
			if slices.Contains(answers.DependencyManagers, models.Pipenv) {
				prefix = "pipenv run "
			} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
				prefix = "poetry run "
			}
			answers.DeployCommand = append(answers.DeployCommand,
				fmt.Sprintf("%spython %s collectstatic --noinput", prefix, managePyPath),
				fmt.Sprintf("%spython %s migrate", prefix, managePyPath),
			)
		}
	case models.Laravel:
		answers.DeployCommand = append(answers.DeployCommand,
			"php artisan optimize:clear",
			"php artisan optimize",
			"php artisan view:clear",
			"php artisan view:cache",
			"php artisan event:clear",
			"php artisan event:cache",
			"php artisan migrate --force",
			
		)
	case models.Flask:
		prefix := ""
		if slices.Contains(answers.DependencyManagers, models.Pipenv) {
			prefix = "pipenv run "
		} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
			prefix = "poetry run "
		}
		answers.DeployCommand = append(answers.DeployCommand,
			"# npm run build",
			fmt.Sprintf("# %sflask run db upgrade", prefix),
		)
	}

	return nil
}
