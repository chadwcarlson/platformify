package question

import (
	"context"
	"path/filepath"

	"github.com/platformsh/platformify/internal/question/models"
	"github.com/platformsh/platformify/internal/utils"
)

type Locations struct{}

func (q *Locations) Ask(ctx context.Context) error {
	answers, ok := models.FromContext(ctx)
	if !ok {
		return nil
	}
	answers.Locations = make(map[string]map[string]interface{})
	switch answers.Stack {
	case models.Django:
		answers.Locations["/static"] = map[string]interface{}{
			"root":    "static",
			"expires": "1h",
			"allow":   true,
		}
	case models.Laravel:
		answers.Locations["/"] = map[string]interface{}{
			"root":    "public",
			"allow":   true,
			"passthru": "/index.php",
			"index": []string{
				"index.php",
			},
		}
		answers.Locations["/storage"] = map[string]interface{}{
			"root":    "storage/app/public",
			"scripts":   false,
		}
	default:
		if answers.Type.Runtime == models.PHP {
			locations := map[string]interface{}{
				"passthru": "/index.php",
				"root":     "",
			}
			if indexPath := utils.FindFile(answers.WorkingDirectory, "index.php"); indexPath != "" {
				indexRelPath, _ := filepath.Rel(answers.WorkingDirectory, indexPath)
				if filepath.Dir(indexRelPath) != "." {
					locations["root"] = filepath.Dir(indexRelPath)
				}
			}
			answers.Locations["/"] = locations
		}
	}

	return nil
}
