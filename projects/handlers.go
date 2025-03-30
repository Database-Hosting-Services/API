package projects

import (
	"DBHS/config"
	"DBHS/response"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func CreateProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := Project{}
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			response.BadRequest(w, "Invalid Input", err)
			return
		}

		has, err := CreateUserProject(r.Context(), config.DB, project.Name, project.Description)
		if err != nil {
			app.ErrorLog.Println("Project creation failed:", err)
			response.InternalServerError(w, "Internal Server Error", errors.New("Project creation failed"))
			return
		}

		if has {
			response.BadRequest(w, "Project with this name already exists", nil)
			return
		}

		response.Created(w, "Project Created", project)
	}
}

// this function returns all projects which the use is the owner of these project
// NOTE : in future plans this function will return also the projects which the user is a member in these projects
func GetProjects(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int)
		data, err := getUserProjects(r.Context(), config.DB, userId)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		response.OK(w, "Projects Retrieved Successfully", data)
	}
}

func getSpecificProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		data, err := getUserSpecificProject(r.Context(), config.DB, userId, projectOid)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		response.OK(w, "Project Retrieved Successfully", data)
	}
}
