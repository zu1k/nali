package repo

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/zu1k/nali/internal/constant"
	"github.com/zu1k/nali/pkg/common"

	"github.com/google/go-github/v55/github"
)

var (
	client     *github.Client
	httpClient = common.GetHttpClient().Client
)

func getLatestRelease() (*github.RepositoryRelease, error) {
	client = github.NewClient(httpClient)
	rel, resp, err := client.Repositories.GetLatestRelease(ctx, constant.Owner, constant.Repo)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// 404 means repository not found or release not found. It's not an error here.
			return nil, fmt.Errorf("repository or release not found")
		}
		return nil, fmt.Errorf("API returned an error response: %s", err)
	}
	if rel == nil {
		return nil, fmt.Errorf("repository release is nil")
	}

	return rel, nil
}

func getTargetAsset(rel *github.RepositoryRelease, sha bool) *github.ReleaseAsset {
	for _, asset := range rel.Assets {
		name := asset.GetName()

		if strings.Contains(name, constant.OS) && strings.Contains(name, constant.Arch) {
			if sha && strings.Contains(name, ".sha256") {
				return asset
			}
			if !sha && !strings.Contains(name, ".sha256") {
				return asset
			}
		}
	}
	return nil
}

func download(ctx context.Context, assetId int64) (data []byte, err error) {
	var rc io.ReadCloser

	rc, _, err = client.Repositories.DownloadReleaseAsset(ctx, constant.Owner, constant.Repo, assetId, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to call GitHub Releases API for getting the asset ID %v on repository '%v/%v': %v", assetId, constant.Owner, constant.Repo, err)
	}
	defer rc.Close()
	data, err = io.ReadAll(rc)

	return
}

func validate(data []byte, vHash []byte) bool {
	cHash := sha256.Sum256(data)

	return bytes.Equal(vHash, cHash[:])
}
