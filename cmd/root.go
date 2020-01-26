package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/hub"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/types"
)

var (
	dockerHubNamespace  string
	dockerHubUser       string
	dockerHubPassword   string
	dockerHubRepository string
)

// Execute a
func Execute(version string) {
	rootCmd := &cobra.Command{
		Use:     "dhdu",
		Short:   "docker hub description updater (dhdu)",
		RunE:    runE,
		Args:    cobra.ExactArgs(1),
		Version: version,
	}
	rootCmd.Flags().StringVarP(&dockerHubNamespace, "namespace", "n", "", "Docker Hub Namespace (default \"username\")")
	rootCmd.Flags().StringVarP(&dockerHubPassword, "password", "p", "", "Docker Hub Password")
	rootCmd.Flags().StringVarP(&dockerHubRepository, "repository", "r", "", "Docker Hub Repository")
	rootCmd.Flags().StringVarP(&dockerHubUser, "username", "u", "", "Docker Hub Username")
	rootCmd.Execute()
}

func runE(cmd *cobra.Command, args []string) error {

	file := args[0]

	if len(dockerHubUser) <= 0 {
		return fmt.Errorf("No user defined over flags")
	}

	if len(dockerHubPassword) <= 0 {
		return fmt.Errorf("No password defined over flags")
	}

	if len(dockerHubNamespace) <= 0 {
		log.Printf("No namespace defined over flags: Use docker username %v instead", dockerHubUser)
		dockerHubNamespace = dockerHubUser
	}

	if len(dockerHubRepository) <= 0 {
		return fmt.Errorf("No repository defined over flags")
	}

	if _, err := os.Stat(file); os.IsNotExist(err) && len(file) <= 0 {
		return fmt.Errorf("Can not find file: %v", file)
	}

	f, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Can not read file %v: %v", file, err)
	}
	fullDescription := string(f)

	loginCredentials := &types.LoginCredentials{
		User:     dockerHubUser,
		Password: dockerHubPassword,
	}

	h := hub.New(loginCredentials)

	repository := &types.Repository{
		Name:            dockerHubRepository,
		Namespcace:      dockerHubNamespace,
		FullDescription: fullDescription,
	}

	_, err = h.PatchRepository(repository)
	if err != nil {
		return err
	}

	return nil
}
