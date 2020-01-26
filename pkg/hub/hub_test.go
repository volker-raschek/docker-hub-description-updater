package hub

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/types"
)

func TestPatchRepository(t *testing.T) {

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

	h := New(loginCredentials)

	require := require.New(t)

	readme, err := Asset("README.md")
	require.NoError(err)

	currentRepository, err := h.GetRepository(dockerHubNamespace, dockerHubRepository)
	require.NoError(err)

	expectedRepository := *currentRepository
	expectedRepository.FullDescription = string(readme)

	actualRepository, err := h.PatchRepository(&expectedRepository)
	require.NoError(err)
	require.EqualValues(&expectedRepository.FullDescription, actualRepository.FullDescription, "Full description not equal")
}
