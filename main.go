package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/Masterminds/semver"
	"github.com/go-flucky/flucky/pkg/logger"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/hub"
	"github.com/volker-raschek/docker-hub-description-updater/pkg/types"
)

var (
	dockerHubAPI        string = "https://hub.docker.com/v2"
	dockerHubUser       string
	dockerHubPassword   string
	dockerHubNamespace  string
	dockerHubRepository string

	shortDescription     string
	shortDescriptionFile string
	fullDescription      string
	fullDescriptionFile  string

	semVersion *semver.Version
	version    string

	flogger logger.Logger
)

func init() {
	// sVersion, err := semver.NewVersion(version)
	// if err != nil {
	// 	log.Fatalf("Can not create new semantic version from %v: %v", version, err)
	// }
	// semVersion = sVersion

	flogger = logger.NewDefaultLogger(logger.LogLevelDebug)
}

func main() {

	flogger.Debug("Parse flags")
	flag.StringVar(&dockerHubUser, "user", "", "Docker Hub Username")
	flag.StringVar(&dockerHubPassword, "password", "", "Docker Hub Password")
	flag.StringVar(&dockerHubNamespace, "namespace", "", "Docker Hub Namespace")
	flag.StringVar(&dockerHubRepository, "repository", "", "Docker Hub Repository")
	flag.StringVar(&shortDescription, "short-description", "", "Short description of the repository ")
	flag.StringVar(&shortDescriptionFile, "short-description-file", "", "Short description of the repository. Override short-description if defined.")
	flag.StringVar(&fullDescription, "full-description", "", "Full description of the repository")
	flag.StringVar(&fullDescriptionFile, "full-description-file", "./README.md", "Full description of the repository. Override full-description if defined.")
	flag.Parse()

	if len(dockerHubUser) <= 0 {
		flogger.Fatal("No user defined over flags")
	}

	if len(dockerHubPassword) <= 0 {
		flogger.Fatal("No password defined over flags")
	}

	if len(dockerHubNamespace) <= 0 {
		flogger.Fatal("No namespace defined over flags")
	}

	if len(dockerHubRepository) <= 0 {
		flogger.Fatal("No repository defined over flags")
	}

	hub.SetLogger(flogger)

	loginCredentials := &types.LoginCredentials{
		User:     dockerHubUser,
		Password: dockerHubPassword,
	}

	actualShortDescription := ""
	if len(shortDescription) > 0 {
		actualShortDescription = shortDescription
		flogger.Debug("Select short description from flag")
	} else if len(shortDescriptionFile) > 0 {
		f, err := ioutil.ReadFile(shortDescriptionFile)
		if err != nil {
			log.Fatalf("Can not read file %v", shortDescriptionFile)
		}
		actualShortDescription = string(f)
		flogger.Debug("Select short description from file")
	}

	actualFullDescription := ""
	if len(fullDescription) > 0 {
		actualFullDescription = fullDescription
		flogger.Debug("Select full description from flag")
	} else if len(fullDescriptionFile) > 0 {
		f, err := ioutil.ReadFile(fullDescriptionFile)
		if err != nil {
			log.Fatalf("Can not read file %v", fullDescriptionFile)
		}
		actualFullDescription = string(f)
		flogger.Debug("Select full description from file")
	}

	flogger.Debug("Get Token")
	token, err := hub.GetToken(loginCredentials)
	if err != nil {
		log.Fatalf("%v", err)
	}

	repository := &types.Repository{
		Name:            dockerHubRepository,
		Namespcace:      dockerHubNamespace,
		Description:     actualShortDescription,
		FullDescription: actualFullDescription,
	}

	flogger.Debug("Send Repository Patch")
	_, err = hub.PatchRepository(repository, token)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
