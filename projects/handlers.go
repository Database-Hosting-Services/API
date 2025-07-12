package projects

import (
	"DBHS/config"
	"DBHS/response"
	"DBHS/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project with the provided details
// @Tags projects
// @Accept json
// @Produce json
// @Param project body projects.CreateProjectRequest true "Project information"
// @Security BearerAuth
// @Success 201 {object} response.SuccessResponse{data=SafeProjectData}
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 500 {object} response.ErrorResponse500
// @Router /projects [post]
func CreateProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project := Project{}
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			response.BadRequest(w, r, "Invalid Input", err)
			return
		}

		has, err, ProjectData := CreateUserProject(r.Context(), config.DB, project.Name, project.Description)
		if err != nil {
			app.ErrorLog.Println("Project creation failed:", err)
			if err.Error() == "Project already exists" {
				response.BadRequest(w, r, "Project already exists", errors.New("Project creation failed"))
			} else if err.Error() == "database name must start with a letter or underscore and contain only letters, numbers, underscores, or $" {
				response.BadRequest(w, r, "database name must start with a letter or underscore and contain only letters, numbers, underscores, or $", errors.New("Project creation failed"))
			} else {
				response.InternalServerError(w, r, "Internal Server Error", err)
			}
			return
		}

		if has {
			response.BadRequest(w, r, "Project with this name already exists", nil)
			return
		}

		response.Created(w, r, "Project Created Successfully", ProjectData)
	}
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse "Project deleted successfully"
// @Failure 400 {object} response.ErrorResponse400 "Project ID is required"
// @Failure 401 {object} response.ErrorResponse401 "Unauthorized"
// @Failure 404 {object} response.ErrorResponse404 "Project not found"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects/{project_id} [delete]
func DeleteProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlVariables := mux.Vars(r)
		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		// Call the function to delete the project
		err := DeleteUserProject(r.Context(), config.DB, projectOid)
		if err != nil {
			app.ErrorLog.Println("Project deletion failed:", err)

			switch err.Error() {
			case "Project not found":
				response.NotFound(w, r, "Project not found", err)
			case "Unauthorized":
				response.UnAuthorized(w, r, "Unauthorized", err)
			default:
				response.InternalServerError(w, r, "Internal Server Error", errors.New("Project deletion failed"))
			}
			return
		}
		response.OK(w, r, "Project Deleted Successfully", nil)
	}
}

// GetProjects godoc
// @Summary Get all user projects
// @Description Get all projects owned by the authenticated user
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]Project} "Projects retrieved successfully"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects [get]
// this function returns all projects which the use is the owner of these project
// NOTE : in future plans this function will return also the projects which the user is a member in these projects
func GetProjects(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int64)
		data, err := getUserProjects(r.Context(), config.DB, userId)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		response.OK(w, r, "Projects Retrieved Successfully", data)
	}
}

// getSpecificProject godoc
// @Summary Get a specific project
// @Description Get details of a specific project by its ID
// @Tags projects
// @Produce json
// @Param project_id path string true "Project ID"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=Project} "Project retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Project ID is required"
// @Failure 404 {object} response.ErrorResponse "Project not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /projects/{project_id} [get]
func getSpecificProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int64)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		data, err := GetUserSpecificProject(r.Context(), config.DB, userId, projectOid)
		if err != nil {
			if errors.Is(err, ErrorProjectNotFound) {
				response.NotFound(w, r, "Project not found", nil)
				return
			}
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		response.OK(w, r, "Project Retrieved Successfully", data)
	}
}

// updateProject godoc
// @Summary Update a project
// @Description Update a project's details by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "Project ID"
// @Param project body updateProjectDataModel true "Project update information"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=Project} "Project updated successfully"
// @Failure 400 {object} response.ErrorResponse400 "Invalid input or Project ID is required"
// @Failure 500 {object} response.ErrorResponse500 "Internal server error"
// @Router /projects/{project_id} [patch]
func updateProject(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user-id").(int64)
		urlVariables := mux.Vars(r)

		projectOid := urlVariables["project_id"]
		if projectOid == "" {
			response.BadRequest(w, r, "Project Id is required", nil)
			return
		}

		var data updateProjectDataModel
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			response.BadRequest(w, r, "Invalid Input", errors.New("The request body is empty"))
			return
		}
		defer r.Body.Close()

		fieldsToUpdate, Values, err := utils.GetNonZeroFieldsFromStruct(&data)
		if err != nil {
			response.BadRequest(w, r, "Invalid Input", err)
			return
		}

		// check if fields to update has a 'name'
		// if it has then validate it
		for idx, field := range fieldsToUpdate {
			if field == "name" {
				err := validateProjectData(r.Context(), config.DB, Values[idx].(string), userId)
				if err != nil {
					response.BadRequest(w, r, "Invalid Input Data", err)
					return
				}
			}
		}

		query, err := BuildProjectUpdateQuery(projectOid, fieldsToUpdate)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", errors.New("error in generating the updating query"))
			return
		}

		transaction, err := config.DB.Begin(r.Context())
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", errors.New("Cannot begin transaction"))
			return
		}

		err = updateProjectData(r.Context(), transaction, query, Values)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", errors.New("Cannot update project data"))
			return
		}

		projectData, err := GetUserSpecificProject(r.Context(), transaction, userId, projectOid)
		if err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", nil)
			return
		}

		if err := transaction.Commit(r.Context()); err != nil {
			app.ErrorLog.Println(err)
			response.InternalServerError(w, r, "Internal Server Error", errors.New("Cannot commit transaction"))
			return
		}

		response.OK(w, r, "Project Retrieved Successfully", projectData)
	}
}
