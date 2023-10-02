package repo

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/zu1k/nali/internal/constant"

	"github.com/google/go-github/v55/github"
)

var (
	ctx      = context.Background()
	tAsset   *github.ReleaseAsset
	shaAsset *github.ReleaseAsset
)

func UpdateRepo() error {
	rel, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %v", err)
	}

	if !canUpdate(rel) {
		log.Printf("current version is already the latest version, no update \n")
		return nil
	}

	//Filtering assets by GOOS and GOARCH
	if tAsset = getTargetAsset(rel, false); tAsset == nil {
		return fmt.Errorf("no target asset found for %s %s", constant.OS, constant.Arch)
	}
	if shaAsset = getTargetAsset(rel, true); shaAsset == nil {
		return fmt.Errorf("no sha256 asset found for %s %s", constant.OS, constant.Arch)
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
	vHash := make([]byte, sha256.Size)
	if _, err := hex.Decode(vHash, vData[:sha256.BlockSize]); err != nil {
		return fmt.Errorf("failed to decode sha256 hash: %v", err)
	}
	if !validate(data, vHash) {
		return fmt.Errorf("failed to validate asset %v, sha256 check failed", tAsset.GetID())
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

	log.Printf("Updating %v to %v \n", exe, rel.GetTagName())
	if err = update(asset, exe); err != nil {
		return fmt.Errorf("update executable failed: %v", err)
	}

	log.Printf("Successfully updated to version %v \n", rel.GetTagName())
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

	// some of our users may install nali through package management, we need to check the permissions before updating
	if !canWriteDir(updateDir) {
		return fmt.Errorf("no write permissions on the directory, consider updating nali manually")
	}

	// Copy the contents of new binary to a new executable file
	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)

	if err != nil {
		_ = fp.Close()
		return fmt.Errorf("create the new executable file failed: %v", err)
	}

	if _, err = io.Copy(fp, bytes.NewReader(newBytes)); err != nil {
		_ = fp.Close()
		return fmt.Errorf("copy the new executable file failed: %v", err)
	}
	// if we don't call fp.Close(), windows won't let us move the new executable
	// because the file will still be "in use"
	if err = fp.Close(); err != nil {
		return fmt.Errorf("failed to close file, may cause file corruption, nali updation was cancelled: %v", err)
	}

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
		log.Printf("remove old binary failed, please remove the old binary manually: %v \n", err)
	}

	return nil
}

func canUpdate(rel *github.RepositoryRelease) bool {
	// unknown version means that the user compiled it manually instead of downloading it from the release,
	// in which case we don't take the liberty of updating it to a potentially older version.
	if constant.Version == "unknown version" {
		return false
	}

	latest, err := parseVersion(rel.GetTagName())
	if err != nil {
		log.Printf("failed to parse latest version: %v, err: %v \n", rel.GetTagName(), err)
		return false
	}

	cur, err := parseVersion(constant.Version)
	if err != nil {
		log.Printf("failed to parse current version: %v, err: %v \n", constant.Version, err)
		return false
	}

	return latest.GreaterThan(cur)
}

func canWriteDir(path string) bool {
	fp := filepath.Join(path, ".tempWriteCheck")
	defer os.Remove(fp)

	file, err := os.Create(fp)
	if err == nil {
		file.Close()
	}

	return err == nil
}
