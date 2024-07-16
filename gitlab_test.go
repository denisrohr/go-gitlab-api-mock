package gitlabapimock_test

import (
	"fmt"
	"net/http"

	"github.com/xanzy/go-gitlab"

	gitlabapimock "github.com/arkadiusjonczek/go-gitlab-api-mock"
)

var (
	GitlabHost = "127.0.0.1:31337"
)

func initGitlabApiMockServer(gitlabMock *gitlabapimock.GitlabMock) *http.Server {
	gitlabApiMock := gitlabapimock.NewGitlabApiMock(gitlabMock)
	return gitlabApiMock.CreateServer(GitlabHost)
}

func initGitlabClient() (*gitlab.Client, error) {
	return gitlab.NewClient("foobar", gitlab.WithBaseURL(fmt.Sprintf("http://%s", GitlabHost)))
}
