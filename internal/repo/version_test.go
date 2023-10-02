package repo

import (
	"reflect"
	"testing"
)

func Test_parseVersion(t *testing.T) {
	tests := []struct {
		name    string
		vStr    string
		want    *Version
		wantErr bool
	}{
		{
			name:    "happy path for normal release version name",
			vStr:    "v0.7.1",
			want:    &Version{Major: 0, Minor: 7, Patch: 1},
			wantErr: false,
		},
		{
			name:    "happy path for old aur-git release version name",
			vStr:    "0.7.3.r15.g43a3080",
			want:    &Version{Major: 0, Minor: 7, Patch: 3},
			wantErr: false,
		},
		{
			name:    "happy path for new aur-git release version name",
			vStr:    "0.7.3-r15-g43a3080",
			want:    &Version{Major: 0, Minor: 7, Patch: 3},
			wantErr: false,
		},
		{
			name:    "empty version name",
			vStr:    "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "only major and minor version",
			vStr:    "0.7",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format version name, major/minor/patch version is not a number",
			vStr:    "xx.xx.xx",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "user customized version name",
			vStr:    "test version",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "only dots, no version number",
			vStr:    "...",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.vStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_GreaterThan(t *testing.T) {
	tests := []struct {
		name  string
		v     *Version
		other *Version
		want  bool
	}{
		{
			name:  "happy path",
			v:     &Version{Major: 0, Minor: 7, Patch: 1},
			other: &Version{Major: 0, Minor: 7, Patch: 0},
			want:  true,
		},
		{
			name:  "same version number",
			v:     &Version{Major: 0, Minor: 7, Patch: 1},
			other: &Version{Major: 0, Minor: 7, Patch: 1},
			want:  false,
		},
		{
			name:  "less than other version number",
			v:     &Version{Major: 0, Minor: 7, Patch: 1},
			other: &Version{Major: 0, Minor: 7, Patch: 3},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.GreaterThan(tt.other); got != tt.want {
				t.Errorf("GreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}
