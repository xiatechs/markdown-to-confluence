package confluence

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/confluence/test/confluencemocks"
	"github.com/xiatechs/markdown-to-confluence/markdown"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestAPIClient_UpdatePage(t *testing.T) {
	newPages := []struct {
		name          string
		pageVersion   int64
		pageID        int
		pageContent   *markdown.FileContents
		setup         func(*confluencemocks.MockHTTPClient)
		expectedError error
	}{
		{
			name:        "happy path, updates page successfully",
			pageVersion: int64(1),
			pageID:      321,
			pageContent: &markdown.FileContents{
				MetaData: map[string]interface{}{"title": "pageTitle"},
				Body:     []byte("some text"),
			},
			setup: func(m *confluencemocks.MockHTTPClient) {
				m.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Status:     "OKI page updated",
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader("some text")),
				}, nil)
			},
			expectedError: nil,
		},
	}

	for _, test := range newPages {
		test := test
		t.Run(test.name, func(t *testing.T) {
			asserts := assert.New(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mock := confluencemocks.NewMockHTTPClient(mockCtrl)

			test.setup(mock)

			setEnvs([]string{"INPUT_CONFLUENCE_USERNAME", "INPUT_CONFLUENCE_API_KEY", "INPUT_CONFLUENCE_SPACE"})

			apiClient := APIClientWithAuths(mock)

			err := apiClient.UpdatePage(test.pageID, test.pageVersion, test.pageContent)

			asserts.Equal(err, test.expectedError)
		})
	}
}

func setEnvs(envs []string) {
	for _, env := range envs {
		os.Setenv(env, "username_space_password")
	}
}
