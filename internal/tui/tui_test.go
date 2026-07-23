package tui

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRenderBreadcrumb_WindowsPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-only test")
	}
	result := renderBreadcrumb(`C:\Users\test\Documents`)
	if !strings.Contains(result, "Users") {
		t.Errorf("should contain 'Users', got: %s", result)
	}
	if !strings.Contains(result, "Documents") {
		t.Errorf("should contain 'Documents', got: %s", result)
	}
	if !strings.Contains(result, ">") {
		t.Errorf("should contain separator '>', got: %s", result)
	}
}

func TestRenderBreadcrumb_RootPath(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-only test")
	}
	result := renderBreadcrumb(`C:\`)
	if result == "" {
		t.Error("should not be empty for root path")
	}
}

func TestRenderBreadcrumb_ShortPath(t *testing.T) {
	dir := filepath.Join("some", "path")
	result := renderBreadcrumb(dir)
	if !strings.Contains(result, "some") {
		t.Errorf("should contain 'some', got: %s", result)
	}
	if !strings.Contains(result, "path") {
		t.Errorf("should contain 'path', got: %s", result)
	}
}

func TestParentDir(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{filepath.Join("a", "b", "c"), filepath.Join("a", "b")},
		{filepath.Join("a", "b"), "a"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parentDir(tt.input)
			if got != tt.want {
				t.Errorf("parentDir(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParentDir_RootDoesNotLoop(t *testing.T) {
	root := filepath.VolumeName("C:") + string(filepath.Separator)
	if root == "" {
		root = string(filepath.Separator)
	}
	got := parentDir(root)
	if got != root {
		t.Errorf("parentDir(%q) should return same path for root, got %q", root, got)
	}
}
