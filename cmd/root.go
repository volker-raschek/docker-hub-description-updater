package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/go-flucky/flucky/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/hub"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/types"
)

var (
	dockerHubNamespace  string
	dockerHubUser       string
	dockerHubPassword   string
	dockerHubRepository string
	file                string

	loglevel string
	flogger  logger.Logger = logger.NewDefaultLogger(logger.LogLevelDebug)
)

var rootCmd = &cobra.Command{
	Use:   "dhdu",
	Short: "docker hub description updater (dhdu)",
	Run: func(cmd *cobra.Command, args []string) {

		switch loglevel {
		case "debug":
			flogger = logger.NewDefaultLogger(logger.LogLevelDebug)
			break
		case "info":
			flogger = logger.NewDefaultLogger(logger.LogLevelDebug)
			break
		case "warn":
			flogger = logger.NewDefaultLogger(logger.LogLevelDebug)
			break
		case "error":
			flogger = logger.NewDefaultLogger(logger.LogLevelDebug)
			break
		case "fatal":
			flogger = logger.NewDefaultLogger(logger.LogLevelDebug)
			break
		}
		hub.SetLogger(flogger)

		if len(dockerHubUser) <= 0 {
			flogger.Fatal("No user defined over flags")
		}

		if len(dockerHubPassword) <= 0 {
			flogger.Fatal("No password defined over flags")
		}

		if len(dockerHubNamespace) <= 0 {
			flogger.Debug("No namespace defined over flags: Use docker username instead")
			dockerHubNamespace = dockerHubUser
		}

		if len(dockerHubRepository) <= 0 {
			flogger.Fatal("No repository defined over flags")
		}

		if _, err := os.Stat(file); os.IsNotExist(err) && len(file) <= 0 {
			flogger.Fatal("Can not find file: %v", file)
		}

		flogger.Debug("Read file %v", file)
		f, err := ioutil.ReadFile(file)
		if err != nil {
			flogger.Fatal("Can not read file %v: %v", file, err)
		}
		fullDescription := string(f)

		loginCredentials := &types.LoginCredentials{
			User:     dockerHubUser,
			Password: dockerHubPassword,
		}

		flogger.Debug("Get Token")
		token, err := hub.GetToken(loginCredentials)
		if err != nil {
			log.Fatalf("%v", err)
		}

		repository := &types.Repository{
			Name:            dockerHubRepository,
			Namespcace:      dockerHubNamespace,
			FullDescription: fullDescription,
		}

		flogger.Debug("Send Repository Patch")
		_, err = hub.PatchRepository(repository, token)
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}

// Execute a
func Execute(version string) {
	rootCmd.Version = version

	rootCmd.Flags().StringVarP(&loglevel, "loglevel", "l", "info", "Loglevel, possible values: debug, info, warn, error, fatal")

	rootCmd.Flags().StringVarP(&dockerHubNamespace, "namespace", "n", "", "Docker Hub Namespace (default \"username\")")
	rootCmd.Flags().StringVarP(&dockerHubPassword, "password", "p", "", "Docker Hub Password")
	rootCmd.Flags().StringVarP(&dockerHubRepository, "repository", "r", "", "Docker Hub Repository")
	rootCmd.Flags().StringVarP(&dockerHubUser, "username", "u", "", "Docker Hub Username")

	rootCmd.Flags().StringVarP(&file, "file", "f", "./README.md", "File which should be uploaded as docker hub description")

	rootCmd.Execute()
}
