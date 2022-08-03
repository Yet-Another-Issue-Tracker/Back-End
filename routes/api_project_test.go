package routes

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateProject(testCase *testing.T) {
	expectedResponse := CreateProjectResponse{
		Id: "123",
	}
	expectedJsonReponse, _ := json.Marshal(expectedResponse)

	testCase.Run("createProject return the new id", func(t *testing.T) {
		response, err := createProject()
		require.Equal(t, nil, err)
		require.Equal(t, string(expectedJsonReponse), string(response))
	})
}
