package gitlabapimock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	gitlabapimock "github.com/arkadiusjonczek/go-gitlab-api-mock.git"
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
