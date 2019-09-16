package types

import "time"

type Repository struct {
	User            string       `json:"user"`
	Name            string       `json:"name"`
	Namespcace      string       `json:"namespace"`
	Type            string       `json:"repository_type"`
	Status          int          `json:"status"`
	Description     string       `json:"description"`
	Private         bool         `json:"is_private"`
	Automated       bool         `json:"is_automated"`
	Edit            bool         `json:"can_edit"`
	StarCount       int          `json:"start_count"`
	PullCount       int          `json:"pull_count"`
	LastUpdated     time.Time    `json:"last_updated"`
	IsMigrated      bool         `json:"is_migrated"`
	HasStarred      bool         `json:"has_starred"`
	FullDescription string       `json:"full_description"`
	Affiliation     string       `json:"affilition"`
	Permissions     *Permissions `json:"permissions"`
}

type Permissions struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
	Admin bool `json:"admin"`
}
