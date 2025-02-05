// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"github.com/zhenorzz/goploy/model"
)

type Repo interface {
	Ping(url string) error
	// Create one repository
	Create(projectID int64) error
	// Follow the repository code and update to latest
	Follow(project model.Project, target string) error
	// RemoteBranchList list remote branches in the url
	RemoteBranchList(url string) ([]string, error)
	// BranchList list the local repository's branches
	BranchList(projectID int64) ([]string, error)
	// CommitLog list the local commit log
	CommitLog(projectID int64, rows int) ([]CommitInfo, error)
	// BranchLog list the local commit log from specific branch
	BranchLog(projectID int64, branch string, rows int) ([]CommitInfo, error)
	// TagLog list the local commit log from all tag
	TagLog(projectID int64, rows int) ([]CommitInfo, error)
}

type CommitInfo struct {
	Branch    string `json:"branch"`
	Commit    string `json:"commit"`
	Author    string `json:"author"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Tag       string `json:"tag"`
	Diff      string `json:"diff"`
}

func GetRepo(repoType string) (Repo, error) {
	if repoType == model.RepoGit {
		return GitRepo{}, nil
	} else if repoType == model.RepoSVN {
		return SvnRepo{}, nil
	} else if repoType == model.RepoFTP {
		return FtpRepo{}, nil
	} else if repoType == model.RepoSFTP {
		return SftpRepo{}, nil
	}
	return nil, fmt.Errorf("wrong repo type passed")
}
