package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const acceptHeader = "application/vnd.docker.distribution.manifest.v2+json"
const credentialsFile = ".credentials"

// Registry credentials structure
type Registry struct {
	Host       string `toml:"nexus_host"`
	Username   string `toml:"nexus_username"`
	Password   string `toml:"nexus_password"`
	Repository string `toml:"nexus_repository"`
}

type repositories struct {
	Images []string `json:"repositories"`
}

type imageTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// ImageManifest : docker registry manifest v2
type ImageManifest struct {
	SchemaVersion int64       `json:"schemaVersion"`
	MediaType     string      `json:"mediaType"`
	Config        layerInfo   `json:"config"`
	Layers        []layerInfo `json:"layers"`
}
type layerInfo struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

// NewRegistry : creates new Registry structure
func NewRegistry() (Registry, error) {
	r := Registry{}
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		return r, fmt.Errorf("%s file not found: %v", credentialsFile, err)
	} else if err != nil {
		return r, err
	}

	if _, err := toml.DecodeFile(credentialsFile, &r); err != nil {
		return r, err
	}
	return r, nil
}

// ListImages : List images in Nexus Docker registry
func (r Registry) ListImages() ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/_catalog", r.Host, r.Repository)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	var repositories repositories
	json.NewDecoder(resp.Body).Decode(&repositories)

	return repositories.Images, nil
}

// ListTagsByImage : list image tags in Nexus Docker registry
func (r Registry) ListTagsByImage(image string) ([]string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/tags/list", r.Host, r.Repository, image)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	var imageTags imageTags
	json.NewDecoder(resp.Body).Decode(&imageTags)

	return imageTags.Tags, nil
}

// ImageManifest : get docker image manifest from registry
func (r Registry) ImageManifest(image string, tag string) (ImageManifest, error) {
	var imageManifest ImageManifest
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return imageManifest, err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return imageManifest, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return imageManifest, fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	json.NewDecoder(resp.Body).Decode(&imageManifest)

	return imageManifest, nil
}

// DeleteImageByTag : delete specific image tag from registry
func (r Registry) DeleteImageByTag(image string, tag string) error {
	sha, err := r.getImageSHA(image, tag)
	if err != nil {
		return err
	}
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, sha)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	fmt.Printf("%s:%s has been successful deleted\n", image, tag)

	return nil
}

func (r Registry) getImageSHA(image string, tag string) (string, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	return resp.Header.Get("docker-content-digest"), nil
}

// GetImageTagDate : get last modified date for the image tag
func (r Registry) GetImageTagDate(image string, tag string) (time.Time, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/repository/%s/v2/%s/manifests/%s", r.Host, r.Repository, image, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return time.Now(), err
	}
	req.SetBasicAuth(r.Username, r.Password)
	req.Header.Add("Accept", acceptHeader)

	resp, err := client.Do(req)
	if err != nil {
		return time.Now(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return time.Now(), fmt.Errorf("HTTP Code: %d", resp.StatusCode)
	}

	t, err := time.Parse(time.RFC1123, resp.Header.Get("last-modified"))

	if err != nil {
		return time.Now(), err
	}

	return t, nil
}
