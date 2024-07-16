package gitlabapimock

import (
	"fmt"
	"sync/atomic"

	"github.com/xanzy/go-gitlab"
)

type GitlabMock struct {
	groupIds         atomic.Int32
	projectIds       atomic.Int32
	projectMemberIds atomic.Int32

	groups         []*gitlab.Group
	projects       map[int]*gitlab.Project
	projectMembers map[int][]*gitlab.ProjectMember
}

func NewGitlabMock() *GitlabMock {
	return &GitlabMock{
		groups:         make([]*gitlab.Group, 0),
		projects:       make(map[int]*gitlab.Project),
		projectMembers: make(map[int][]*gitlab.ProjectMember),
	}
}

func (mock *GitlabMock) AddGroup(name string) *gitlab.Group {
	id := int(mock.groupIds.Add(1))

	group := &gitlab.Group{
		ID:       id,
		Name:     name,
		FullName: name,
		Path:     name,
	}

	mock.groups = append(mock.groups, group)

	return group
}

func (mock *GitlabMock) GetGroups() []*gitlab.Group {
	return mock.groups
}

func (mock *GitlabMock) AddProject(name string, group *gitlab.Group) *gitlab.Project {
	id := int(mock.projectIds.Add(1))

	project := &gitlab.Project{
		ID:   id,
		Name: name,
		Path: name,
	}

	group.Projects = append(group.Projects, project)
	mock.projects[id] = project

	return project
}

func (mock *GitlabMock) GetProjects() []*gitlab.Project {
	var projects []*gitlab.Project

	for _, group := range mock.groups {
		projects = append(projects, mock.projects[group.ID])
	}

	return projects
}

func (mock *GitlabMock) AddProjectMember(projectMember *gitlab.ProjectMember, project *gitlab.Project) *gitlab.Project {
	mock.projectMembers[project.ID] = append(mock.projectMembers[project.ID], projectMember)

	return project
}

func (mock *GitlabMock) GetProjectMembers(projectID int) ([]*gitlab.ProjectMember, error) {
	projectMembers, projectExists := mock.projectMembers[projectID]
	if !projectExists {
		fmt.Errorf("project %d not found", projectID)
	}

	return projectMembers, nil
}
