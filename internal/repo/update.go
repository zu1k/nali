package repo

import (
	"bytes"
	"context"
	"fmt"
	semver "github.com/Masterminds/semver/v3"
	"github.com/google/go-github/v55/github"
	"github.com/zu1k/nali/internal/constant"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	ctx          = context.Background()
	reVersion    = regexp.MustCompile(`\d+\.\d+\.\d+`)
	assetSuffixs = []string{".zip", ".gz"}
	versionRe    = "(v)?\\d+\\.\\d+\\.\\d+"
	tAsset       *github.ReleaseAsset
	shaAsset     *github.ReleaseAsset
)

func UpdateRepo(version string) error {
	rel, err := getLatestRelease(constant.Owner, constant.Repo)
	if err != nil {
		return fmt.Errorf("failed to get latest release: %v", err)
	}

	// get the latest version and compare version numbers, if not the latest, update it
	if constant.Version != "unknown version" {
		latest, _ := semver.NewVersion(rel.GetTagName())
		cur, _ := semver.NewVersion(constant.Version)

		if cur.GreaterThan(latest) {
			return fmt.Errorf("current version %v is greater or equal to the latest version %v, no update", constant.Version, rel.GetTagName())
		}
	}

	assetFilters := makeFilters()

	//Filtering assets by GOOS and GOARCH
	if tAsset = getTargetAsset(rel, assetFilters); tAsset == nil {
		return fmt.Errorf("could not find target asset from release: %v", rel.GetTagName())
	}

	shaFilter := tAsset.GetName() + ".sha256"
	shaRe, err := regexp.Compile(shaFilter)
	if err != nil {
		return fmt.Errorf("could not compile regular expression %v for filtering releases sha256: %v", shaFilter, err)
	}
	if shaAsset = getTargetAsset(rel, []*regexp.Regexp{shaRe}); shaAsset == nil {
		return fmt.Errorf("could not find target sha256 asset from release: %v", rel.GetTagName())
	}

	//Download the new version nali and its sha256
	data, err := download(ctx, tAsset.GetID())
	if err != nil {
		return fmt.Errorf("failed to download asset %v: %v", tAsset.GetID(), err)
	}

	vData, err := download(ctx, shaAsset.GetID())
	if err != nil {
		return fmt.Errorf("failed to download asset %v: %v", tAsset.GetID(), err)
	}

	// Verifying files with sha256
	if err = validate(data, vData); err != nil {
		return fmt.Errorf("failed to validate asset %v: %v", tAsset.GetID(), err)
	}

	// Unzip and replace nali itself
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not locate executable path: %v", err)
	}

	asset, err := decompress(bytes.NewReader(data), tAsset.GetName())
	if err != nil {
		return fmt.Errorf("error occurred while decompress: %v", err)
	}

	log.Printf("Will update %v to %v downloaded from %v", exe, version, tAsset.GetBrowserDownloadURL())
	if err = update(asset, exe); err != nil {
		return fmt.Errorf("update executable failed: %v", err)
	}

	log.Printf("Successfully updated to version %v", rel.GetTagName())
	return nil
}

func update(asset io.Reader, cmdPath string) error {
	newBytes, err := io.ReadAll(asset)
	if err != nil {
		return err
	}

	// get the directory the executable exists in
	updateDir := filepath.Dir(cmdPath)
	filename := filepath.Base(cmdPath)

	// Copy the contents of new binary to a new executable file
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("create the new executable file failed: %v", err)
	}

	if _, err = io.Copy(fp, bytes.NewReader(newBytes)); err != nil {
		return fmt.Errorf("copy the new executable file failed: %v", err)
	}
	fp.Close()

	oldPath := filepath.Join(updateDir, fmt.Sprintf(".%s.old", filename))

	// delete any existing old exec file - this is necessary on Windows for two reasons:
	// 1. after a successful asset, Windows can't remove the .old file because the process is still running
	// 2. windows rename operations fail if the destination file already exists
	_ = os.Remove(oldPath)

	if err = os.Rename(cmdPath, oldPath); err != nil {
		return fmt.Errorf("rename the old executable file failed: %v", err)
	}

	if err = os.Rename(newPath, cmdPath); err != nil {
		// move unsuccessful
		// The filesystem is now in a bad state. We have successfully
		// moved the existing binary to a new location, but we couldn't move the new
		// binary to take its place. That means there is no file where the current executable binary
		// used to be!
		// Try to rollback by restoring the old binary to its original path.
		if rerr := os.Rename(oldPath, cmdPath); rerr != nil {
			return fmt.Errorf("unable to rollback binary: %v", rerr)
		}

		return fmt.Errorf("unable to move new binary to executable path: %v", err)
	}

	if err = os.Remove(oldPath); err != nil {
		// windows has trouble with removing old binaries, so do nothing only print log
		log.Printf("remove old binary failed, please remove the old binary manually: %v", err)
	}

	return nil
}

func makeFilters() []*regexp.Regexp {
	assetFilters := make([]*regexp.Regexp, 0)

	for _, ext := range assetSuffixs {
		filter := fmt.Sprintf("%s%c%s%c%s%c%s%s",
			constant.Repo, constant.Sep,
			constant.OS, constant.Sep,
			constant.Arch, constant.Sep,
			versionRe, ext)

		re, err := regexp.Compile(filter)
		if err != nil {
			log.Printf("could not compile regular expression %v for filtering releases: %v", filter, err)
			continue
		}
		assetFilters = append(assetFilters, re)
	}
	return assetFilters
}
