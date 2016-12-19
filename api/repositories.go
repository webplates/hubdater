package api

import (
	"fmt"
)

type RepositoriesService service

type Autobuild struct {
	Repository int    `json:"repository"`
	BuildName  string `json:"build_name"`
	Provider   string `json:"provider"`
	SourceUrl  string `json:"source_url"`
	DockerUrl  string `json:"docker_url"`
	RepoWebUrl string `json:"repo_web_url"`
	RepoType   string `json:"repo_type"`
	Active     bool   `json:"active"`
	RepoId     string `json:"repo_id"`
	BuildTags  []Tag  `json:"build_tags"`
}

type Tag struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	DockerfileLocation string `json:"dockerfile_location"`
	SourceName         string `json:"source_name"`
	SourceType         string `json:"source_type"`
}

func (s *RepositoriesService) GetAutobuild(name string) (*Autobuild, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("repositories/%s/autobuild/", name), nil)
	if err != nil {
		return nil, err
	}

	autobuild := new(Autobuild)

	resp, err := s.client.Do(req, autobuild)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, s.client.ParseError(fmt.Sprintf("get details for %s", name), resp)
	}

	return autobuild, nil
}

func (s *RepositoriesService) CreateTag(repo string, tag *Tag) error {
	req, err := s.client.NewRequest("POST", fmt.Sprintf("repositories/%s/autobuild/tags/", repo), tag)
	if err != nil {
		return err
	}

	if err := s.client.Authenticate(req); err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return s.client.ParseError(fmt.Sprintf("create tag %s for repo %s", tag.Name, repo), resp)
	}

	return nil
}

func (s *RepositoriesService) UpdateTag(repo string, tag *Tag) error {
	req, err := s.client.NewRequest("PUT", fmt.Sprintf("repositories/%s/autobuild/tags/%d/", repo, tag.Id), tag)
	if err != nil {
		return err
	}

	if err := s.client.Authenticate(req); err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return s.client.ParseError(fmt.Sprintf("update tag %s in repo %s", tag.Name, repo), resp)
	}

	return nil
}

func (s *RepositoriesService) DeleteTag(repo string, tag *Tag) error {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("repositories/%s/autobuild/tags/%d/", repo, tag.Id), nil)
	if err != nil {
		return err
	}

	if err := s.client.Authenticate(req); err != nil {
		return err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return s.client.ParseError(fmt.Sprintf("delete tag %s from repo %s", tag.Name, repo), resp)
	}

	return nil
}
