package hub_test

import (
	"os"
	"testing"

	flogger "github.com/go-flucky/flucky/pkg/logger"
	"github.com/stretchr/testify/require"
	"github.com/volker-raschek/dhd/pkg/hub"
	"github.com/volker-raschek/dhd/pkg/types"
)

func TestPatchRepository(t *testing.T) {

	hub.SetLogger(flogger.NewDefaultLogger(flogger.LogLevelDebug))

	dockerHubUser := os.Getenv("REGISTRY_USER")
	if len(dockerHubUser) <= 0 {
		t.Fatalf("Environment variable REGISTRY_USER is empty")
	}

	dockerHubPassword := os.Getenv("REGISTRY_PASSWORD")
	if len(dockerHubPassword) <= 0 {
		t.Fatalf("Environment variable REGISTRY_PASSWORD is empty")
	}

	dockerHubNamespace := os.Getenv("REGISTRY_NAMESPACE")
	if len(dockerHubNamespace) <= 0 {
		t.Fatalf("Environment variable REGISTRY_NAMESPACE is empty")
	}

	dockerHubRepository := os.Getenv("CONTAINER_IMAGE_NAME")
	if len(dockerHubRepository) <= 0 {
		t.Fatalf("Environment variable CONTAINER_IMAGE_NAME is empty")
	}

	loginCredentials := &types.LoginCredentials{
		User:     dockerHubUser,
		Password: dockerHubPassword,
	}

	require := require.New(t)
	token, err := hub.GetToken(loginCredentials)
	require.NoError(err)

	readme, err := Asset("README.md")
	require.NoError(err)

	currentRepository, err := hub.GetRepository(dockerHubNamespace, dockerHubRepository, token)
	require.NoError(err)

	expectedRepository := *currentRepository
	expectedRepository.FullDescription = string(readme)

	actualRepository, err := hub.PatchRepository(&expectedRepository, token)
	require.NoError(err)

	require.NotEqual(currentRepository, actualRepository, "The repository properties have remained the same even though an update was performed")
	require.Equal(&expectedRepository, actualRepository, "The update was successfully")
}
