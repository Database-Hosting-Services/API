package projects

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"net/http"
)

func CreateProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := Project{}
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			response.BadRequest(w, "Invalid Input", err)
			return
		}

		has, err := CreateUserProject(r.Context(), config.DB, project.Name)
		if err != nil {
			app.ErrorLog.Println("Project creation failed:", err)
			response.BadRequest(w, "Project creation failed", err)
			return
		}

		if has {
			response.BadRequest(w, "Project with this name already exists", nil)
			return
		}

		response.Created(w, "Project Created", project)
	}
}
