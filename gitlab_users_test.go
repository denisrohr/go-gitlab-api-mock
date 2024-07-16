package gitlabapimock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	gitlabapimock "github.com/arkadiusjonczek/go-gitlab-api-mock"
)

func Test_Users_ListUsers_ReturnsEmptyList(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listUserOptions := &gitlab.ListUsersOptions{}
	users, response, err := gitlabClient.Users.ListUsers(listUserOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, users, 0)
}

func Test_Users_ListUsers_ReturnsUsers(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	gitlabMock.AddUser("Peter Pan", "peter.pan", "peter.pan@telekom.de")
	gitlabMock.AddUser("Petra Pan", "petra.pan", "petra.pan@telekom.de")
	gitlabMock.AddUser("Fred Feuerstein", "fred.feuerstein", "fred.feuerstein@telekom.de")

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listUserOptions := &gitlab.ListUsersOptions{}
	users, response, err := gitlabClient.Users.ListUsers(listUserOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, users, 3)
}

func Test_Users_ListUsers_ReturnsSearchedUsername(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	gitlabMock.AddUser("Peter Pan", "peter.pan", "peter.pan@telekom.de")
	gitlabMock.AddUser("Petra Pan", "petra.pan", "petra.pan@telekom.de")
	gitlabMock.AddUser("Fred Feuerstein", "fred.feuerstein", "fred.feuerstein@telekom.de")

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listUserOptions := &gitlab.ListUsersOptions{
		Username: gitlab.Ptr("peter.pan"),
	}
	users, response, err := gitlabClient.Users.ListUsers(listUserOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, users, 1)

	require.Equal(t, "Peter Pan", users[0].Name)
	require.Equal(t, "peter.pan", users[0].Username)
	require.Equal(t, "peter.pan@telekom.de", users[0].Email)
}

func Test_Users_ListUsers_ReturnsSearchedUser(t *testing.T) {
	gitlabMock := gitlabapimock.NewGitlabMock()
	gitlabMock.AddUser("Peter Pan", "peter.pan", "peter.pan@telekom.de")
	gitlabMock.AddUser("Petra Pan", "petra.pan", "petra.pan@telekom.de")
	gitlabMock.AddUser("Fred Feuerstein", "fred.feuerstein", "fred.feuerstein@telekom.de")

	server := initGitlabApiMockServer(gitlabMock)

	go server.ListenAndServe()
	defer server.Close()

	time.Sleep(1 * time.Second)

	gitlabClient, err := initGitlabClient()

	require.NoError(t, err)

	listUserOptions := &gitlab.ListUsersOptions{
		Search: gitlab.Ptr("peter.pan@telekom.de"),
	}
	users, response, err := gitlabClient.Users.ListUsers(listUserOptions)

	require.NoError(t, err)
	require.Equal(t, 200, response.StatusCode)
	require.Len(t, users, 1)

	require.Equal(t, "Peter Pan", users[0].Name)
	require.Equal(t, "peter.pan", users[0].Username)
	require.Equal(t, "peter.pan@telekom.de", users[0].Email)
}
