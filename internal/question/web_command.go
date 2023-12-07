package question

import (
	"context"
	"fmt"
	// "os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
	"github.com/platformsh/platformify/vendorization"
)

type WebCommand struct{}

func (q *WebCommand) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	if answers.WebCommand != "" {
		// Skip the step
		return nil
	}

	// Do not ask the command for PHP applications
	if answers.Type.Runtime == models.PHP {
		return nil
	}

	assets, _ := vendorization.FromContext(ctx)


	//nolint:lll
	answers.WebCommand = fmt.Sprintf(
		"echo 'Put your web server command in here! You need to listen to \"$UNIX\" unix socket. Read more about it here: %s#web-commands'; sleep 60",
		assets.Docs().AppReference,
	)

	if answers.SocketFamily == models.TCP {
		//nolint:lll
		answers.WebCommand = fmt.Sprintf(
			"echo 'Put your web server command in here! You need to listen to \"$PORT\" port. Read more about it here: %s#web-commands'; sleep 60",

			assets.Docs().AppReference,
		)
	}

	switch answers.Stack {
	case models.Django:
		prefix := ""
		pythonPath := ""
		wsgi := "app.wsgi"
		// try to find the wsgi.py file to change the default command
		wsgiPath := utils.FindFile(path.Join(answers.WorkingDirectory, answers.ApplicationRoot), "wsgi.py")
		if wsgiPath != "" {
			wsgiParentDir := path.Base(path.Dir(wsgiPath))
			wsgi = fmt.Sprintf("%s.wsgi", wsgiParentDir)

			// add the pythonpath if the wsgi.py file is not in the root of the app
			wsgiRel, _ := filepath.Rel(
				path.Join(answers.WorkingDirectory, answers.ApplicationRoot),
				path.Dir(path.Dir(wsgiPath)),
			)
			if wsgiRel != "." {
				pythonPath = " --pythonpath=" + path.Base(path.Dir(path.Dir(wsgiPath)))
			}
		}

		if slices.Contains(answers.DependencyManagers, models.Pipenv) {
			prefix = "pipenv run "
		} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
			prefix = "poetry run "
		}
		if answers.SocketFamily == models.TCP {
			answers.WebCommand = fmt.Sprintf("%sgunicorn%s -b 0.0.0.0:$PORT %s --log-file -", prefix, pythonPath, wsgi)
			return nil
		}

		answers.WebCommand = fmt.Sprintf("%sgunicorn%s -b unix:$SOCKET %s --log-file -", prefix, pythonPath, wsgi)
		return nil
	case models.NextJS:
		answers.WebCommand = "npx next start -p $PORT"
		return nil
	case models.Strapi:
		if _, ok := utils.GetJSONValue(
			[]string{"scripts", "start"},
			path.Join(answers.WorkingDirectory, "package.json"),
			true,
		); ok {
			if slices.Contains(answers.DependencyManagers, models.Yarn) {
				answers.WebCommand = "NODE_ENV=production yarn start"
			} else {
				answers.WebCommand = "NODE_ENV=production npm start"
			}
		}
	case models.Express:
		if _, ok := utils.GetJSONValue(
			[]string{"scripts", "start"},
			path.Join(answers.WorkingDirectory, "package.json"),
			true,
		); ok {
			if slices.Contains(answers.DependencyManagers, models.Yarn) {
				answers.WebCommand = "NODE_ENV=production yarn start"
			} else {
				answers.WebCommand = "NODE_ENV=production npm start"
			}
			return nil
		}

		if mainPath, ok := utils.GetJSONValue(
			[]string{"main"},
			path.Join(answers.WorkingDirectory, "package.json"),
			true,
		); ok {
			answers.WebCommand = fmt.Sprintf("node %s", mainPath.(string))
			return nil
		}

		if indexFile := utils.FindFile(answers.WorkingDirectory, "index.js"); indexFile != "" {
			indexFile, _ = filepath.Rel(answers.WorkingDirectory, indexFile)
			answers.WebCommand = fmt.Sprintf("node %s", indexFile)
			return nil
		}
	case models.Flask:

		prefix := ""
		if slices.Contains(answers.DependencyManagers, models.Pipenv) {
			prefix = "pipenv run "
		} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
			prefix = "poetry run "
		}

		// Using FLASK_APP + either gunicorn/flask run to build start command.
		if ok, flask_app, _ := utils.FindFlaskApp(answers.WorkingDirectory); ok {
			if slices.Contains(answers.DependencyManagers, models.Pip) {
				requirementsPath := utils.FindFile(answers.WorkingDirectory, "requirements.txt")
				if requirementsPath != "" {
					if _, ok := utils.DepInNestedRequirements("gunicorn", requirementsPath, true); ok {
						answers.WebCommand = fmt.Sprintf("%sgunicorn -b unix:$SOCKET %s:app --log-file -", prefix, strings.TrimSuffix(flask_app, ".py"))
						return nil
					} else {
						answers.WebCommand = fmt.Sprintf("%sflask run -p $PORT", prefix)
						answers.SocketFamily = ""
						return nil
					}
				}
			} else if slices.Contains(answers.DependencyManagers, models.Poetry) {
				pyProjectPath := utils.FindFile(answers.WorkingDirectory, "pyproject.toml")
				if pyProjectPath != "" {
					if _, ok := utils.GetTOMLValue([]string{"tool", "poetry", "dependencies", "gunicorn"}, pyProjectPath, true); ok {
						answers.WebCommand = fmt.Sprintf("%sgunicorn -b unix:$SOCKET %s:app --log-file -", prefix, strings.TrimSuffix(flask_app, ".py"))
						return nil
					} else {
						answers.WebCommand = fmt.Sprintf("%sflask run -p $PORT", prefix)
						answers.SocketFamily = ""
						return nil
					}
				}
			} else if slices.Contains(answers.DependencyManagers, models.Pipenv) {
				pipfilePath := utils.FindFile(answers.WorkingDirectory, pipenvFile)
				if pipfilePath != "" {
					if _, ok := utils.GetTOMLValue([]string{"packages", "gunicorn"}, pipfilePath, true); ok {
						answers.WebCommand = fmt.Sprintf("%sgunicorn -b unix:$SOCKET %s:app --log-file -", prefix, strings.TrimSuffix(flask_app, ".py"))
						return nil
					} else {
						answers.WebCommand = fmt.Sprintf("%sflask run -p $PORT", prefix)
						answers.SocketFamily = ""
						return nil
					}
				}
			}
		}


		// appPath := ""
		// // try to find the app.py, api.py or server.py files
		// for _, name := range []string{"app.py", "server.py", "api.py"} {
		// 	if _, err := os.Stat(path.Join(answers.WorkingDirectory, name)); err == nil {
		// 		appPath = fmt.Sprintf("'%s:app'", strings.TrimSuffix(name, ".py"))
		// 		break
		// 	}
		// }
		// if appPath == "" {
		// 	return nil
		// }

		// prefix := ""
		// if slices.Contains(answers.DependencyManagers, models.Pipenv) {
		// 	prefix = "pipenv run "
		// } else if slices.Contains(answers.DependencyManagers, models.Poetry) {
		// 	prefix = "poetry run "
		// }

		// if answers.SocketFamily == models.TCP {
		// 	answers.WebCommand = fmt.Sprintf("%sgunicorn -b 0.0.0.0:$PORT %s --log-file -", prefix, appPath)
		// 	return nil
		// }

		// answers.WebCommand = fmt.Sprintf("%sgunicorn -b unix:$SOCKET %s --log-file -", prefix, appPath)
		return nil


	}

	return nil
}
