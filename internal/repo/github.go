package repo

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/v55/github"
	"github.com/zu1k/nali/internal/constant"
	"io"
	"net/http"
	"regexp"
)

var (
	client *github.Client
)

func getLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	client = github.NewClient(http.DefaultClient)
	rel, resp, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
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

func getTargetAsset(rel *github.RepositoryRelease, assetFilters []*regexp.Regexp) *github.ReleaseAsset {
	var tAsset *github.ReleaseAsset
	found := false

	for _, asset := range rel.Assets {
		name := asset.GetName()

		for _, filter := range assetFilters {
			if filter.MatchString(name) {
				tAsset = asset
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	return tAsset
}

func download(ctx context.Context, assetId int64) (data []byte, err error) {
	var rc io.ReadCloser

	rc, _, err = client.Repositories.DownloadReleaseAsset(ctx, constant.Owner, constant.Repo, assetId, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("failed to call GitHub Releases API for getting the asset ID %v on repository '%v/%v': %v", assetId, constant.Owner, constant.Repo, err)
	}
	defer func() { _ = rc.Close() }()
	data, err = io.ReadAll(rc)

	return
}

func validate(data, vData []byte) error {
	if len(vData) < sha256.BlockSize {
		return fmt.Errorf("incorrect checksum file format")
	}

	hash := fmt.Sprintf("%s", vData[:sha256.BlockSize])
	calculatedHash := fmt.Sprintf("%x", sha256.Sum256(data))

	if equal, err := hexStringEquals(sha256.Size, calculatedHash, hash); !equal {
		if err == nil {
			return fmt.Errorf("expected %q, found %q: sha256 validation failed", hash, calculatedHash)
		}
		return fmt.Errorf("%s: sha256 validation failed", err.Error())
	}
	return nil
}

func hexStringEquals(size int, a, b string) (equal bool, err error) {
	size *= 2
	if len(a) == size && len(b) == size {
		var bytesA, bytesB []byte
		if bytesA, err = hex.DecodeString(a); err == nil {
			if bytesB, err = hex.DecodeString(b); err == nil {
				equal = bytes.Equal(bytesA, bytesB)
			}
		}
	}
	return
}
