package question

import (
	"context"
	// "fmt"
	"path"
	"path/filepath"

	// "golang.org/x/exp/slices"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
	// "github.com/platformsh/platformify/vendorization"
)

type BuildSteps struct{}

func (q *BuildSteps) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}

	for _, dm := range answers.DependencyManagers {
		switch dm {
		case models.Poetry:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Set PIP_USER to 0 so that Poetry does not complain",
				"export PIP_USER=0",
				"# Install poetry as a global tool",
				"python -m venv /app/.global",
				"pip install poetry==$POETRY_VERSION",
				"poetry install",
			)
		case models.Pipenv:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Set PIP_USER to 0 so that Pipenv does not complain",
				"export PIP_USER=0",
				"# Install Pipenv as a global tool",
				"python -m venv /app/.global",
				"pip install pipenv==$PIPENV_TOOL_VERSION",
				"pipenv install",
			)
		case models.Pip:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"# Upgrade to the latest version of pip.",
				"pip install --upgrade pip",
				"# Install dependencies.",
				"pip install -r requirements.txt",
			)
		case models.Yarn, models.Npm:
			if answers.Type.Runtime != models.NodeJS {
				answers.Dependencies["nodejs"]["n"] = "*"
				answers.Dependencies["nodejs"]["npx"] = "*"
				answers.Environment["N_PREFIX"] = "/app/.global"
				answers.BuildSteps = append(
					answers.BuildSteps,
					"n auto || n lts",
					"hash -r",
				)
			}

			if dm == models.Yarn {
				answers.BuildSteps = append(
					answers.BuildSteps,
					"yarn",
				)
			} else {
				answers.BuildSteps = append(
					answers.BuildSteps,
					"npm i",
				)
			}
			if _, ok := utils.GetJSONValue(
				[]string{"scripts", "build"},
				path.Join(answers.WorkingDirectory, "package.json"),
				true,
			); ok {
				if dm == models.Yarn {
					answers.BuildSteps = append(answers.BuildSteps, "yarn build")
				} else {
					answers.BuildSteps = append(answers.BuildSteps, "npm run build")
				}
			}
		case models.Composer:
			answers.BuildSteps = append(
				answers.BuildSteps,
				"composer --no-ansi --no-interaction install --no-progress --prefer-dist --optimize-autoloader --no-dev",
			)
		}
	}

	if answers.Stack == models.Django {
		if managePyPath := utils.FindFile(
			path.Join(answers.WorkingDirectory, answers.ApplicationRoot),
			managePyFile,
		); managePyPath != "" {
			// prefix := ""
			// if slices.Contains(answers.DependencyManagers, models.Pipenv) {
			// 	prefix = "pipenv run "
			// } else if slices.Contains(answers.DependencyManagers, models.Poetry) {
			// 	prefix = "poetry run "
			// }

			managePyPath, _ = filepath.Rel(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), managePyPath)
			// assets, _ := vendorization.FromContext(ctx)
			answers.BuildSteps = append(
				answers.BuildSteps,
				"pip install --upgrade pip",
				"pip install -r requirements.txt",
				// fmt.Sprintf(
				// 	"# Collect static files so that they can be served by %s.",
				// 	assets.ServiceName,
				// ),
				// fmt.Sprintf("%spython %s collectstatic --noinput\n", prefix, managePyPath),
				"npm install",
				"npm run build",
			)
		}
	}

	if answers.Stack == models.Laravel {
		answers.BuildSteps = append(
			answers.BuildSteps,
			"n auto || n lts",
			"hash -r",
			"npm install",
			"npm run production",
		)
	}

	return nil
}
