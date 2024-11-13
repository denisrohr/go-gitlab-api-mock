package gitlabapimock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	gitlabapimock "github.com/arkadiusjonczek/go-gitlab-api-mock"
)

func Test_Groups_ListGroups_ReturnsEmptyList(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listGroupOptions := &gitlab.ListGroupsOptions{}
	groups, response, err := gitlabClient.Groups.ListGroups(listGroupOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, groups, 0)
}

func Test_Groups_ListGroups_ReturnsGroups(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	gitlabMock.AddGroup("group1")
	gitlabMock.AddGroup("group2")
	gitlabMock.AddGroup("group3")

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listGroupOptions := &gitlab.ListGroupsOptions{}
	groups, response, err := gitlabClient.Groups.ListGroups(listGroupOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, groups, 3)
}

func Test_Groups_GetInheritedGroupsAccessLevelWithHigherProjectPermissions(t *testing.T) {
	// group1->group2->group3->project1
	// projectMember1 is reporter in group2, but owner in project1, so we expect owner accesslevel
	// projectMember2 is maintainer in group1 and not in project1, so we expect maintainer accesslevel
	// projectMember3 is not in any group, but reporter in project1, so we expect reporter accesslevel
	gitlabMock := gitlabapimock.NewGitlabMock()
	group1 := gitlabMock.AddGroup("group1")
	group2 := gitlabMock.AddGroupWithParent("group2", group1.ID)
	group3 := gitlabMock.AddGroupWithParent("group3", group2.ID)

	project1 := gitlabMock.AddProject("project1", group3)
	projectMember1 := &gitlab.ProjectMember{
		ID:          1,
		Username:    "member1",
		Email:       "member1@gitlab.com",
		Name:        "member1",
		AccessLevel: gitlab.OwnerPermissions,
	}
	projectMember3 := &gitlab.ProjectMember{
		ID:          3,
		Username:    "member3",
		Email:       "member3@gitlab.com",
		Name:        "member3",
		AccessLevel: gitlab.ReporterPermissions,
	}
	gitlabMock.AddProjectMember(projectMember1, project1)
	gitlabMock.AddProjectMember(projectMember3, project1)

	groupMember1 := &gitlab.GroupMember{
		ID:          1,
		Username:    "member1",
		Email:       "member1@gitlab.com",
		Name:        "member1",
		AccessLevel: gitlab.OwnerPermissions,
	}
	groupMember2 := &gitlab.GroupMember{
		ID:          2,
		Username:    "member2",
		Email:       "member2@gitlab.com",
		Name:        "member2",
		AccessLevel: gitlab.MaintainerPermissions,
	}

	gitlabMock.AddGroupMember(groupMember1, group2)
	gitlabMock.AddGroupMember(groupMember2, group1)

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listProjectMembersOptions := &gitlab.ListProjectMembersOptions{}
	projectMembers, response, err := gitlabClient.ProjectMembers.ListAllProjectMembers(project1.ID, listProjectMembersOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, projectMembers, 3)

	membershipsFound := 0
	for _, projectMember := range projectMembers {
		switch projectMember.ID {
		case 1:
			require.Equal(t, gitlab.OwnerPermissions, projectMember.AccessLevel)
			membershipsFound++
		case 2:
			require.Equal(t, gitlab.MaintainerPermissions, projectMember.AccessLevel)
			membershipsFound++
		case 3:
			require.Equal(t, gitlab.ReporterPermissions, projectMember.AccessLevel)
			membershipsFound++
		default:
		}
	}
	require.Equal(t, 3, membershipsFound)
}
