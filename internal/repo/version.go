package repo

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func parseVersion(vStr string) (*Version, error) {
	vStr = strings.TrimPrefix(vStr, "v")

	// split by hyphen to remove pre-release info, then split by dot to get version info
	parts := strings.Split(strings.Split(vStr, "-")[0], ".")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid version format: %s", vStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	return &Version{Major: major, Minor: minor, Patch: patch}, nil
}

func (v *Version) Equal(other *Version) bool {
	return v.compare(other) == 0
}

func (v *Version) GreaterThan(other *Version) bool {
	return v.compare(other) > 0
}

func (v *Version) LessThan(other *Version) bool {
	return v.compare(other) < 0
}

func (v *Version) compare(other *Version) int {
	if v.Major != other.Major {
		return v.Major - other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}
	return v.Patch - other.Patch
}
