package projects

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
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

		has, err, ProjectData := CreateUserProject(r.Context(), config.DB, project.Name, project.Description)
		if err != nil {
			app.ErrorLog.Println("Project creation failed:", err)
			if err.Error() == "Project already exists" {
				response.BadRequest(w, "Project already exists", errors.New("Project creation failed"))
			} else if err.Error() == "database name must start with a letter or underscore and contain only letters, numbers, underscores, or $" {
				response.BadRequest(w, "database name must start with a letter or underscore and contain only letters, numbers, underscores, or $", errors.New("Project creation failed"))
			} else {
				response.InternalServerError(w, "Internal Server Error", errors.New("Project creation failed"))
			}
			return
		}

		if has {
			response.BadRequest(w, "Project with this name already exists", nil)
			return
		}

		response.Created(w, "Project Created Successfully", ProjectData)
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

func updateProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, "Project Id is required", nil)
			return
		}

		var data updateProjectDataModel
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			response.BadRequest(w, "Invalid Input", errors.New("The request body is empty"))
			return
		}
		defer r.Body.Close()

		fieldsToUpdate, Values, err := utils.GetNonZeroFieldsFromStruct(&data)
		if err != nil {
			response.BadRequest(w, "Invalid Input", err)
			return
		}

		query, err := BuildProjectUpdateQuery(projectOid, fieldsToUpdate)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", errors.New("error in generating the updating query"))
			return
		}

		transaction, err := config.DB.Begin(r.Context())
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", errors.New("Cannot begin transaction"))
			return
		}

		err = updateProjectData(r.Context(), transaction, query, Values)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", errors.New("Cannot update project data"))
			return
		}

		projectData, err := getUserSpecificProject(r.Context(), transaction, userId, projectOid)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", nil)
			return
		}

		if err := transaction.Commit(r.Context()); err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, "Internal Server Error", errors.New("Cannot commit transaction"))
			return
		}

		response.OK(w, "Project Retrieved Successfully", projectData)
	}
}
