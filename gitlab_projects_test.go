package gitlabapimock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	gitlabapimock "github.com/arkadiusjonczek/go-gitlab-api-mock.git"
)

func Test_Projects_ListProjects_ReturnsEmptyList(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listProjectsOptions := &gitlab.ListProjectsOptions{}
	projects, response, err := gitlabClient.Projects.ListProjects(listProjectsOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, projects, 0)
}

func Test_Projects_ListProjects_ReturnsProjects(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	gitlabMock.AddProject("project1", group1)
	gitlabMock.AddProject("project2", group1)
	gitlabMock.AddProject("project3", group1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listProjectsOptions := &gitlab.ListProjectsOptions{}
	projects, response, err := gitlabClient.Projects.ListProjects(listProjectsOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, projects, 3)
}

func Test_Projects_ListProjectMembers_ReturnsEmptyList(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listProjectMemberOptions := &gitlab.ListProjectMembersOptions{}
	_, _, err = gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.Error(t, err)
}

func Test_Projects_ListProjectMembers_ReturnsProjectMembers(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	project1 := gitlabMock.AddProject("project1", group1)

	projectMember1 := &gitlab.ProjectMember{
		ID:          1,
		Username:    "member1",
		Email:       "member1@gitlab.com",
		Name:        "member1",
		AccessLevel: gitlab.ReporterPermissions,
	}

	gitlabMock.AddProjectMember(projectMember1, project1)

	projectMember2 := &gitlab.ProjectMember{
		ID:          2,
		Username:    "member2",
		Email:       "member21@gitlab.com",
		Name:        "member2",
		AccessLevel: gitlab.DeveloperPermissions,
	}

	gitlabMock.AddProjectMember(projectMember2, project1)

	projectMember3 := &gitlab.ProjectMember{
		ID:          3,
		Username:    "member3",
		Email:       "member3@gitlab.com",
		Name:        "member3",
		AccessLevel: gitlab.MaintainerPermissions,
	}

	gitlabMock.AddProjectMember(projectMember3, project1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listProjectMemberOptions := &gitlab.ListProjectMembersOptions{}
	projectMembers, response, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, projectMembers, 3)
}

func Test_Projects_AddProjectMember_ReturnsProjectMember(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	gitlabMock.AddProject("project1", group1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	addProjectMemberOptions := &gitlab.AddProjectMemberOptions{
		UserID:      1,
		AccessLevel: gitlab.Ptr(gitlab.DeveloperPermissions),
	}
	projectMember, response, err := gitlabClient.ProjectMembers.AddProjectMember(1, addProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, projectMember)
	require.Equal(t, gitlab.DeveloperPermissions, projectMember.AccessLevel)

	listProjectMemberOptions := &gitlab.ListProjectMembersOptions{}
	projectMembers, response, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, projectMembers, 1)

	require.Equal(t, gitlab.DeveloperPermissions, projectMembers[0].AccessLevel)
}

func Test_Projects_EditProjectMember_ReturnsProjectMember(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	gitlabMock.AddProject("project1", group1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	addProjectMemberOptions := &gitlab.AddProjectMemberOptions{
		UserID:      1,
		AccessLevel: gitlab.Ptr(gitlab.MaintainerPermissions),
	}
	addedProjectMember, addProjectMemberResponse, err := gitlabClient.ProjectMembers.AddProjectMember(1, addProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, addProjectMemberResponse.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, addedProjectMember)
	require.Equal(t, gitlab.MaintainerPermissions, addedProjectMember.AccessLevel)

	listProjectMemberOptions := &gitlab.ListProjectMembersOptions{}
	projectMembersAfterAdd, responseAfterAdd, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterAdd.StatusCode)
	require.Len(t, projectMembersAfterAdd, 1)

	require.Equal(t, gitlab.MaintainerPermissions, projectMembersAfterAdd[0].AccessLevel)

	editProjectMemberOptions := &gitlab.EditProjectMemberOptions{
		AccessLevel: gitlab.Ptr(gitlab.DeveloperPermissions),
	}
	updatedProjectMember, updatedProjectMemberResponse, err := gitlabClient.ProjectMembers.EditProjectMember(1, 1, editProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, updatedProjectMemberResponse.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, updatedProjectMember)
	require.Equal(t, gitlab.DeveloperPermissions, updatedProjectMember.AccessLevel)

	listProjectMemberOptions2 := &gitlab.ListProjectMembersOptions{}
	projectMembersAfterEdit, responseAfterEdit, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions2)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterEdit.StatusCode)
	require.Len(t, projectMembersAfterEdit, 1)

	require.Equal(t, gitlab.DeveloperPermissions, projectMembersAfterEdit[0].AccessLevel)
}

func Test_Projects_DeleteProjectMember_ReturnsOK(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	gitlabMock.AddProject("project1", group1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	addProjectMemberOptions := &gitlab.AddProjectMemberOptions{
		UserID:      1,
		AccessLevel: gitlab.Ptr(gitlab.MaintainerPermissions),
	}
	addedProjectMember, addProjectMemberResponse, err := gitlabClient.ProjectMembers.AddProjectMember(1, addProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, addProjectMemberResponse.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, addedProjectMember)
	require.Equal(t, gitlab.MaintainerPermissions, addedProjectMember.AccessLevel)

	listProjectMemberOptions := &gitlab.ListProjectMembersOptions{}
	projectMembersAfterAdd, responseAfterAdd, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterAdd.StatusCode)
	require.Len(t, projectMembersAfterAdd, 1)

	require.Equal(t, gitlab.MaintainerPermissions, projectMembersAfterAdd[0].AccessLevel)

	deletedProjectMemberResponse, err := gitlabClient.ProjectMembers.DeleteProjectMember(1, 1)

	require.NoError(t, err)
	require.Equal(t, 200, deletedProjectMemberResponse.StatusCode)

	listProjectMemberOptions2 := &gitlab.ListProjectMembersOptions{}
	projectMembersAfterDelete, responseAfterDelete, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions2)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterDelete.StatusCode)
	require.Len(t, projectMembersAfterDelete, 0)
}

func Test_Projects_DeleteMultipleProjectMembers_ReturnsOK(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	gitlabMock.AddProject("project1", group1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	addProjectMemberOptions1 := &gitlab.AddProjectMemberOptions{
		UserID:      1,
		AccessLevel: gitlab.Ptr(gitlab.ReporterPermissions),
	}
	addedProjectMember1, addProjectMemberResponse1, err := gitlabClient.ProjectMembers.AddProjectMember(1, addProjectMemberOptions1)

	require.NoError(t, err)
	require.Equal(t, 200, addProjectMemberResponse1.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, addedProjectMember1)
	require.Equal(t, gitlab.ReporterPermissions, addedProjectMember1.AccessLevel)

	addProjectMemberOptions2 := &gitlab.AddProjectMemberOptions{
		UserID:      2,
		AccessLevel: gitlab.Ptr(gitlab.DeveloperPermissions),
	}
	addedProjectMember2, addProjectMemberResponse2, err := gitlabClient.ProjectMembers.AddProjectMember(1, addProjectMemberOptions2)

	require.NoError(t, err)
	require.Equal(t, 200, addProjectMemberResponse2.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, addedProjectMember2)
	require.Equal(t, gitlab.DeveloperPermissions, addedProjectMember2.AccessLevel)

	addProjectMemberOptions3 := &gitlab.AddProjectMemberOptions{
		UserID:      3,
		AccessLevel: gitlab.Ptr(gitlab.MaintainerPermissions),
	}
	addedProjectMember3, addProjectMemberResponse3, err := gitlabClient.ProjectMembers.AddProjectMember(1, addProjectMemberOptions3)

	require.NoError(t, err)
	require.Equal(t, 200, addProjectMemberResponse3.StatusCode)
	require.IsType(t, &gitlab.ProjectMember{}, addedProjectMember3)
	require.Equal(t, gitlab.MaintainerPermissions, addedProjectMember3.AccessLevel)

	listProjectMemberOptions := &gitlab.ListProjectMembersOptions{}
	projectMembersAfterAdd, responseAfterAdd, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterAdd.StatusCode)
	require.Len(t, projectMembersAfterAdd, 3)

	require.Equal(t, gitlab.ReporterPermissions, projectMembersAfterAdd[0].AccessLevel)
	require.Equal(t, gitlab.DeveloperPermissions, projectMembersAfterAdd[1].AccessLevel)
	require.Equal(t, gitlab.MaintainerPermissions, projectMembersAfterAdd[2].AccessLevel)

	deletedProjectMemberResponse, err := gitlabClient.ProjectMembers.DeleteProjectMember(1, 1)

	require.NoError(t, err)
	require.Equal(t, 200, deletedProjectMemberResponse.StatusCode)

	projectMembersAfterDelete, responseAfterDelete, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterDelete.StatusCode)
	require.Len(t, projectMembersAfterDelete, 2)

	require.Equal(t, gitlab.DeveloperPermissions, projectMembersAfterDelete[0].AccessLevel)
	require.Equal(t, gitlab.MaintainerPermissions, projectMembersAfterDelete[1].AccessLevel)

	deletedProjectMemberResponse2, err := gitlabClient.ProjectMembers.DeleteProjectMember(1, 3)

	require.NoError(t, err)
	require.Equal(t, 200, deletedProjectMemberResponse2.StatusCode)

	projectMembersAfterDelete2, responseAfterDelete2, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterDelete2.StatusCode)
	require.Len(t, projectMembersAfterDelete2, 1)

	require.Equal(t, gitlab.DeveloperPermissions, projectMembersAfterDelete2[0].AccessLevel)

	deletedProjectMemberResponse3, err := gitlabClient.ProjectMembers.DeleteProjectMember(1, 2)

	require.NoError(t, err)
	require.Equal(t, 200, deletedProjectMemberResponse3.StatusCode)

	projectMembersAfterDelete3, responseAfterDelete3, err := gitlabClient.ProjectMembers.ListProjectMembers(1, listProjectMemberOptions)

	require.NoError(t, err)
	require.Equal(t, 200, responseAfterDelete3.StatusCode)
	require.Len(t, projectMembersAfterDelete3, 0)
}
