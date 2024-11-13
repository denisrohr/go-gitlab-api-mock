package gitlabapimock

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/xanzy/go-gitlab"
)

// GitlabMock
// TODO: Add error handling to Add methods
type GitlabMock struct {
	userIds    atomic.Int32
	groupIds   atomic.Int32
	projectIds atomic.Int32

	//ID->entity
	users    map[int]*gitlab.User
	groups   map[int]*gitlab.Group
	projects map[int]*gitlab.Project

	//groupID->GroupMembers
	groupMembers map[int][]*gitlab.GroupMember
	//projectID->ProjectMembers
	projectMembers map[int][]*gitlab.ProjectMember

	//projectID->groupID
	projectsToGroup map[int]int
	//groupID->projectIDs
	groupToProjects map[int][]int
}

func NewGitlabMock() *GitlabMock {
	mock := &GitlabMock{
		users:          make(map[int]*gitlab.User),
		groups:         make(map[int]*gitlab.Group),
		projects:       make(map[int]*gitlab.Project),
		groupMembers:   make(map[int][]*gitlab.GroupMember),
		projectMembers: make(map[int][]*gitlab.ProjectMember),

		projectsToGroup: make(map[int]int),
		groupToProjects: make(map[int][]int),
	}
	return mock
}

func (mock *GitlabMock) AddUser(name string, username string, email string) (*gitlab.User, error) {
	for _, user := range mock.users {
		if user.Username == username {
			return nil, errors.New("user with that username already exists")
		} else if user.Email == email {
			return nil, errors.New("user with that email already exists")
		}
	}

	id := int(mock.userIds.Add(1))

	user := &gitlab.User{
		ID:       id,
		Name:     name,
		Username: username,
		Email:    email,
	}

	mock.users[id] = user

	return user, nil
}

func (mock *GitlabMock) GetUsers() []*gitlab.User {
	var users []*gitlab.User

	for _, user := range mock.users {
		users = append(users, user)
	}

	return users
}

func (mock *GitlabMock) AddGroup(name string) *gitlab.Group {
	id := int(mock.groupIds.Add(1))

	group := &gitlab.Group{
		ID:       id,
		Name:     name,
		FullName: name,
		Path:     name,
		ParentID: 0,
	}

	mock.groups[id] = group

	return group
}

func (mock *GitlabMock) AddGroupWithParent(name string, parentID int) *gitlab.Group {
	group := mock.AddGroup(name)
	group.ParentID = parentID

	return group
}

func (mock *GitlabMock) GetGroups() []*gitlab.Group {
	var groups []*gitlab.Group

	for _, group := range mock.groups {
		groups = append(groups, group)
	}

	return groups
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
	mock.projectsToGroup[id] = group.ID

	return project
}

func (mock *GitlabMock) GetProjects() []*gitlab.Project {
	var projects []*gitlab.Project

	for _, project := range mock.projects {
		projects = append(projects, project)
	}

	return projects
}

func (mock *GitlabMock) AddProjectMember(projectMember *gitlab.ProjectMember, project *gitlab.Project) *gitlab.Project {
	mock.projectMembers[project.ID] = append(mock.projectMembers[project.ID], projectMember)

	return project
}

func (mock *GitlabMock) AddGroupMember(groupMember *gitlab.GroupMember, group *gitlab.Group) *gitlab.Group {
	mock.groupMembers[group.ID] = append(mock.groupMembers[group.ID], groupMember)

	return group
}

func (mock *GitlabMock) GetProjectMembers(projectID int) ([]*gitlab.ProjectMember, error) {
	projectMembers, projectExists := mock.projectMembers[projectID]
	if !projectExists {
		return nil, fmt.Errorf("project %d not found", projectID)
	}

	return projectMembers, nil
}

func (mock *GitlabMock) GetGroupMembers(groupID int) ([]*gitlab.GroupMember, error) {
	groupMembers, groupExists := mock.groupMembers[groupID]
	if !groupExists {
		return nil, fmt.Errorf("group %d not found", groupID)
	}

	return groupMembers, nil
}

func (mock *GitlabMock) MapGroupMemberToProjectMember(groupMember *gitlab.GroupMember) *gitlab.ProjectMember {
	projectMember := &gitlab.ProjectMember{
		ID:          groupMember.ID,
		Username:    groupMember.Username,
		Email:       groupMember.Email,
		Name:        groupMember.Name,
		State:       groupMember.State,
		CreatedAt:   groupMember.CreatedAt,
		ExpiresAt:   groupMember.ExpiresAt,
		AccessLevel: groupMember.AccessLevel,
		WebURL:      groupMember.WebURL,
		AvatarURL:   groupMember.AvatarURL,
	}
	return projectMember
}

func (mock *GitlabMock) GetProjectMembersWithInheritedGroupMembers(projectID int) ([]*gitlab.ProjectMember, error) {
	_, projectExists := mock.projects[projectID]
	if !projectExists {
		return nil, fmt.Errorf("project %d not found", projectID)
	}

	directProjectMembers := mock.projectMembers[projectID]
	groupIDOfProject, groupExists := mock.projectsToGroup[projectID]
	if !groupExists {
		return directProjectMembers, nil
	}

	//group inheritance is involved
	//map to keep track of users
	mapOfMembers := make(map[int]*gitlab.ProjectMember)
	for _, directProjectMember := range directProjectMembers {
		mapOfMembers[directProjectMember.ID] = directProjectMember
	}

	currentGroupID := groupIDOfProject
	for currentGroupID != 0 {
		currentGroup, groupExists := mock.groups[currentGroupID]
		if !groupExists {
			break
		}
		groupMembers, err := mock.GetGroupMembers(currentGroupID)
		if err == nil {
			for _, groupMember := range groupMembers {
				member, memberAlreadyParsed := mapOfMembers[groupMember.ID]
				if memberAlreadyParsed {
					// upgrade accesslevel with group access level if higher
					// real gitlab tries to prevent the user from creating a membership with lower access than the parentGroup membership
					if groupMember.AccessLevel > member.AccessLevel {
						mapOfMembers[groupMember.ID].AccessLevel = groupMember.AccessLevel
					}
				} else {
					// this is the usual case where a user has implicit accessLevels via groupMembership
					mapOfMembers[groupMember.ID] = mock.MapGroupMemberToProjectMember(groupMember)
				}
			}
		}
		currentGroupID = currentGroup.ParentID
	}

	result := make([]*gitlab.ProjectMember, 0)
	for _, value := range mapOfMembers {
		result = append(result, value)
	}
	return result, nil
}
