package confluence

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xiatechs/markdown-to-confluence/confluence/test/confluencemocks"
	"github.com/xiatechs/markdown-to-confluence/markdown"
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
					Body:       io.NopCloser(strings.NewReader("some text")),
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

			envs := []string{"INPUT_CONFLUENCE_USERNAME", "INPUT_CONFLUENCE_API_KEY", "INPUT_CONFLUENCE_SPACE"}
			setEnvs(envs, true)
			defer setEnvs(envs, false)

			apiClient := APIClientWithAuths(mock)

			oldResults := PageResults{
				Results: []Page{
					{
						ID: "test",
						Body: BodyObj{
							Storage: StorageObj{
								Value: "testing",
							},
						},
					},
				},
			}

			_, err := apiClient.UpdatePage(test.pageID, test.pageVersion, test.pageContent, oldResults)

			asserts.Equal(err, test.expectedError)
		})
	}
}

func setEnvs(envs []string, setEnvs bool) {
	if setEnvs {
		for _, env := range envs {
			os.Setenv(env, "username_space_password")
		}
	} else {
		for _, env := range envs {
			os.Unsetenv(env)
		}
	}
}

func TestAPIClient_FindPage(t *testing.T) {
	returnedPage := PageResults{Results: []Page{{
		ID:      "321",
		Type:    "page",
		Title:   "PageTitle",
		Version: VersionObj{2},
		Body: BodyObj{Storage: StorageObj{
			Value: "some text",
		}},
	}}}

	returnedJSON, err := json.Marshal(returnedPage)
	if err != nil {
		fmt.Println("error marshaling test data: ", err)
	}

	fmt.Println("test data: ", string(returnedJSON))

	pageInputs := []struct {
		Name            string
		PageTitle       string
		Setup           func(m *confluencemocks.MockHTTPClient)
		ExpectedID      int
		ExpectedVersion int64
		ExpectedBool    bool
		ExpectedErr     error
	}{
		{
			Name:      "happy path found page",
			PageTitle: "TestPageHappy",
			Setup: func(m *confluencemocks.MockHTTPClient) {
				m.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Status:     "OK, Page Found",
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(string(returnedJSON))),
				}, nil)
			},
			ExpectedBool:    true,
			ExpectedErr:     nil,
			ExpectedID:      321,
			ExpectedVersion: int64(2),
		},
	}

	for _, test := range pageInputs {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			asserts := assert.New(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mock := confluencemocks.NewMockHTTPClient(mockCtrl)

			test.Setup(mock)

			envs := []string{"INPUT_CONFLUENCE_USERNAME", "INPUT_CONFLUENCE_API_KEY", "INPUT_CONFLUENCE_SPACE"}
			setEnvs(envs, true)
			defer setEnvs(envs, false)

			client := APIClientWithAuths(mock)
			_, err := client.FindPage(test.PageTitle, false)

			asserts.Equal(test.ExpectedErr, err)
		})
	}
}

func TestAPIClient_CreatePage(t *testing.T) {
	returnedPage := Page{
		ID:      "321",
		Type:    "page",
		Title:   "PageTitle",
		Version: VersionObj{2},
		Body: BodyObj{Storage: StorageObj{
			Value: "some text",
		},
		},
	}

	returnedJSON, err := json.Marshal(returnedPage)
	if err != nil {
		fmt.Println("error marshaling test data: ", err)
	}

	inputs := []struct {
		name          string
		pageContent   *markdown.FileContents
		setUp         func(*confluencemocks.MockHTTPClient)
		expectedError error
	}{
		{
			name: "Happy path, creates page",
			pageContent: &markdown.FileContents{
				MetaData: map[string]interface{}{"title": "pageTitle"},
				Body:     []byte("some text"),
			},
			setUp: func(m *confluencemocks.MockHTTPClient) {
				m.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Status:     "OK, Page Found",
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(string(returnedJSON))),
				}, nil)
			},
			expectedError: nil,
		},
		{
			name: "un-happy path",
			pageContent: &markdown.FileContents{
				MetaData: map[string]interface{}{"title": "pageTitle"},
				Body:     []byte("some text"),
			},
			setUp: func(m *confluencemocks.MockHTTPClient) {
				m.EXPECT().Do(gomock.Any()).Return(&http.Response{
					Status:     "Not Found",
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil)
			},
			expectedError: fmt.Errorf("failed to create confluence page: %s", "Not Found"),
		},
	}

	for _, test := range inputs {
		test := test
		t.Run(test.name, func(t *testing.T) {
			asserts := assert.New(t)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mock := confluencemocks.NewMockHTTPClient(mockCtrl)

			test.setUp(mock)

			envs := []string{"INPUT_CONFLUENCE_USERNAME", "INPUT_CONFLUENCE_API_KEY", "INPUT_CONFLUENCE_SPACE"}
			setEnvs(envs, true)
			defer setEnvs(envs, false)

			client := APIClientWithAuths(mock)
			_, err := client.CreatePage(0, test.pageContent, true)
			asserts.Equal(test.expectedError, err)
		})
	}
}

func TestAPIClient_UploadAttachment(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	asserts := assert.New(t)
	mock := confluencemocks.NewMockHTTPClient(mockCtrl)

	defer mockCtrl.Finish()

	envs := []string{"INPUT_CONFLUENCE_USERNAME", "INPUT_CONFLUENCE_API_KEY", "INPUT_CONFLUENCE_SPACE"}
	setEnvs(envs, true)

	defer setEnvs(envs, false)

	client := APIClientWithAuths(mock)
	err := client.UploadAttachment("thisfiledoesnotexist", 0, false, 0)

	asserts.Equal(err.Error(), "file upload error: open thisfiledoesnotexist: no such file or directory")
}
