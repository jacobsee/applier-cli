package apiresponses

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetLatestVersionInfo() GitHubRelease {
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get("https://api.github.com/repos/redhat-cop/openshift-applier/releases/latest")
	if err != nil {
		// return err
	}
	defer r.Body.Close()
	response := &GitHubRelease{}
	err = json.NewDecoder(r.Body).Decode(response)
	return *response
}

// GitHubRelease is what you get back when you ask GitHub's API for info on the latest release.
type GitHubRelease struct {
	URL             string        `json:"url"`
	AssetsURL       string        `json:"assets_url"`
	UploadURL       string        `json:"upload_url"`
	HTMLURL         string        `json:"html_url"`
	ID              int           `json:"id"`
	NodeID          string        `json:"node_id"`
	TagName         string        `json:"tag_name"`
	TargetCommitish string        `json:"target_commitish"`
	Name            string        `json:"name"`
	Draft           bool          `json:"draft"`
	Author          GitHubAuthor  `json:"author"`
	Prerelease      bool          `json:"prerelease"`
	CreatedAt       time.Time     `json:"created_at"`
	PublishedAt     time.Time     `json:"published_at"`
	Assets          []interface{} `json:"assets"`
	TarballURL      string        `json:"tarball_url"`
	ZipballURL      string        `json:"zipball_url"`
	Body            string        `json:"body"`
}

// GitHubAuthor is a sub-component of GitHubLatestRelease
type GitHubAuthor struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
