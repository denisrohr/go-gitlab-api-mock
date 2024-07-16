package gitlabapimock

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xanzy/go-gitlab"
)

const (
	GitlabApiPrefix = "/api/v4/"
)

// GitlabApiMock
// TODO: separate handler from business logic
type GitlabApiMock struct {
	gitlabMock *GitlabMock
}

func NewGitlabApiMock(gitlabMock *GitlabMock) *GitlabApiMock {
	return &GitlabApiMock{
		gitlabMock: gitlabMock,
	}
}

func (mock *GitlabApiMock) CreateServer(addr string) *http.Server {
	r := mux.NewRouter().PathPrefix(GitlabApiPrefix).Subrouter()

	r.HandleFunc("/groups", mock.ListGroupsHandler).Methods(http.MethodGet)
	r.HandleFunc("/projects", mock.ListProjectsHandler).Methods(http.MethodGet)
	r.HandleFunc("/projects/{id}/members", mock.ListAllMembersOfAProjectsHandler).Methods(http.MethodGet)
	r.HandleFunc("/projects/{id}/members", mock.AddMemberToAProjectsHandler).Methods(http.MethodPost)
	r.HandleFunc("/projects/{id}/members/{user_id}", mock.EdifMemberOfAProjectHandler).Methods(http.MethodPut)
	r.HandleFunc("/projects/{id}/members/{user_id}", mock.DeleteMemberFromAProjectHandler).Methods(http.MethodDelete)

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return server
}

// ListGroupsHandler implements https://docs.gitlab.com/ee/api/groups.html#list-groups
func (mock *GitlabApiMock) ListGroupsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	err := json.NewEncoder(responseWriter).Encode(mock.gitlabMock.groups)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	return
}

// ListProjectsHandler implements https://docs.gitlab.com/ee/api/projects.html#list-all-projects
func (mock *GitlabApiMock) ListProjectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var projects []*gitlab.Project

	for _, group := range mock.gitlabMock.groups {
		projects = append(projects, group.Projects...)
	}

	err := json.NewEncoder(responseWriter).Encode(projects)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	return
}

// ListAllMembersOfAProjectsHandler implements https://docs.gitlab.com/ee/api/members.html#list-all-members-of-a-group-or-project
func (mock *GitlabApiMock) ListAllMembersOfAProjectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id := vars["id"]
	idInteger, _ := strconv.Atoi(id)

	_, projectExists := mock.gitlabMock.projects[idInteger]
	if !projectExists {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	projectMembers := mock.gitlabMock.projectMembers[idInteger]

	err := json.NewEncoder(responseWriter).Encode(projectMembers)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	return
}

// AddMemberToAProjectsHandler implements https://docs.gitlab.com/ee/api/members.html#add-a-member-to-a-group-or-project
func (mock *GitlabApiMock) AddMemberToAProjectsHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id := vars["id"]
	idInteger, _ := strconv.Atoi(id)

	_, projectExists := mock.gitlabMock.projects[idInteger]
	if !projectExists {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	var addProjectMemberOptions gitlab.AddProjectMemberOptions
	err := json.NewDecoder(request.Body).Decode(&addProjectMemberOptions)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	projectMembers := mock.gitlabMock.projectMembers[idInteger]

	found := false
	for _, member := range projectMembers {
		if member.ID == addProjectMemberOptions.UserID {
			found = true
			break
		}
	}

	userId := addProjectMemberOptions.UserID.(float64)
	userIdInteger := int(userId)

	projectMember := &gitlab.ProjectMember{
		ID:          userIdInteger,
		AccessLevel: *addProjectMemberOptions.AccessLevel,
	}

	if !found {
		mock.gitlabMock.projectMembers[idInteger] = append(projectMembers, projectMember)
	}

	err = json.NewEncoder(responseWriter).Encode(projectMember)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	return
}

// EdifMemberOfAProjectHandler implements https://docs.gitlab.com/ee/api/members.html#edit-a-member-of-a-group-or-project
func (mock *GitlabApiMock) EdifMemberOfAProjectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id := vars["id"]
	idInteger, _ := strconv.Atoi(id)

	_, projectExists := mock.gitlabMock.projects[idInteger]
	if !projectExists {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	userId := vars["user_id"]
	userIdInteger, _ := strconv.Atoi(userId)

	var editProjectMemberOptions gitlab.EditProjectMemberOptions
	err := json.NewDecoder(request.Body).Decode(&editProjectMemberOptions)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	var projectMember *gitlab.ProjectMember

	projectMembers := mock.gitlabMock.projectMembers[idInteger]

	found := false
	for idx, member := range projectMembers {
		if member.ID == userIdInteger {
			found = true
			projectMember = projectMembers[idx]
			break
		}
	}

	if !found {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	if editProjectMemberOptions.AccessLevel != nil {
		log.Printf("UserID = %d, Old Access Level = %d, New Access Level = %d", userIdInteger, projectMember.AccessLevel, *editProjectMemberOptions.AccessLevel)
		projectMember.AccessLevel = *editProjectMemberOptions.AccessLevel
	}

	err = json.NewEncoder(responseWriter).Encode(projectMember)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}

	return
}

// DeleteMemberFromAProjectHandler implements https://docs.gitlab.com/ee/api/members.html#remove-a-member-from-a-group-or-project
func (mock *GitlabApiMock) DeleteMemberFromAProjectHandler(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id := vars["id"]
	idInteger, _ := strconv.Atoi(id)

	_, projectExists := mock.gitlabMock.projects[idInteger]
	if !projectExists {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	userId := vars["user_id"]
	userIdInteger, _ := strconv.Atoi(userId)

	projectMembers := mock.gitlabMock.projectMembers[idInteger]

	found := false
	for idx, member := range projectMembers {
		if member.ID == userIdInteger {
			found = true

			copy(projectMembers[idx:], projectMembers[idx+1:])
			projectMembers[len(projectMembers)-1] = nil
			projectMembers = projectMembers[:len(projectMembers)-1]

			mock.gitlabMock.projectMembers[idInteger] = projectMembers
			break
		}
	}

	if !found {
		responseWriter.WriteHeader(http.StatusNotFound)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	return
}
