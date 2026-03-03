package updater

import (
	"bytes"
	"io"
	"net/http"
	"runtime"
	"testing"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	return m.response, m.err
}

func TestGetAssetName(t *testing.T) {
	t.Parallel()

	got := GetAssetName()
	if got == "" {
		t.Error("GetAssetName() should not be empty")
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH

	if goos == "windows" {
		if got != "nordlys-installer-windows-"+goarch+".exe" {
			t.Errorf("GetAssetName() = %q, want nordlys-installer-windows-%s.exe", got, goarch)
		}
	} else {
		expectedName := "nordlys-installer-" + goos + "-" + goarch
		if got != expectedName {
			t.Errorf("GetAssetName() = %q, want %q", got, expectedName)
		}
	}
}

func TestGetUpdateCheckCachePath(t *testing.T) {
	t.Parallel()

	got := GetUpdateCheckCachePath()
	if got == "" {
		t.Error("GetUpdateCheckCachePath() should not be empty")
	}

	if got[0] != '/' && got[1] != ':' {
		t.Errorf("GetUpdateCheckCachePath() should return absolute path, got %q", got)
	}
}

func TestCheckForUpdate(t *testing.T) {
	t.Parallel()

	_, _, err := CheckForUpdate()
	if err != nil {
		t.Logf("CheckForUpdate() error = %v (expected if no network or repo doesn't exist yet)", err)
	}
}

func TestUpdater_CheckForUpdate_NewerVersion(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "v2.0.0", "assets": []}`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	version, needsUpdate, err := u.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate() error = %v", err)
	}

	if !needsUpdate {
		t.Error("CheckForUpdate() needsUpdate = false, want true")
	}

	if version != "2.0.0" {
		t.Errorf("CheckForUpdate() version = %q, want %q", version, "2.0.0")
	}
}

func TestUpdater_CheckForUpdate_SameVersion(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "v1.0.0", "assets": []}`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	version, needsUpdate, err := u.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate() error = %v", err)
	}

	if needsUpdate {
		t.Error("CheckForUpdate() needsUpdate = true, want false")
	}

	if version != "1.0.0" {
		t.Errorf("CheckForUpdate() version = %q, want %q", version, "1.0.0")
	}
}

func TestUpdater_CheckForUpdate_NoPrefix(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "2.0.0", "assets": []}`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	version, needsUpdate, err := u.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate() error = %v", err)
	}

	if !needsUpdate {
		t.Error("CheckForUpdate() needsUpdate = false, want true")
	}

	if version != "2.0.0" {
		t.Errorf("CheckForUpdate() version = %q, want %q", version, "2.0.0")
	}
}

func TestUpdater_CheckForUpdate_EmptyTag(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "", "assets": []}`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	version, needsUpdate, err := u.CheckForUpdate()
	if err != nil {
		t.Fatalf("CheckForUpdate() error = %v", err)
	}

	if needsUpdate {
		t.Error("CheckForUpdate() needsUpdate = true, want false")
	}

	if version != "" {
		t.Errorf("CheckForUpdate() version = %q, want empty", version)
	}
}

func TestUpdater_CheckForUpdate_HTTPError(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	_, _, err := u.CheckForUpdate()
	if err == nil {
		t.Error("CheckForUpdate() should return error for non-200 status")
	}
}

func TestUpdater_CheckForUpdate_InvalidJSON(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	_, _, err := u.CheckForUpdate()
	if err == nil {
		t.Error("CheckForUpdate() should return error for invalid JSON")
	}
}

func TestNewUpdater(t *testing.T) {
	t.Parallel()

	u := NewUpdater()
	if u.client == nil {
		t.Error("NewUpdater() client should not be nil")
	}
	if u.currentVer == "" {
		t.Error("NewUpdater() currentVer should not be empty")
	}
	if u.releasesURL == "" {
		t.Error("NewUpdater() releasesURL should not be empty")
	}
}

func TestUpdater_SelfUpdate_AlreadyLatest(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "v1.0.0", "assets": []}`)),
		},
	}

	u := NewUpdaterWithClient(mockClient, "1.0.0", "https://api.test.com/releases/latest")
	err := u.SelfUpdate()
	if err == nil {
		t.Error("SelfUpdate() should return error when already at latest version")
	}
}

func TestUpdater_SelfUpdate_NoAsset(t *testing.T) {
	t.Parallel()

	callCount := 0
	mockClient := &mockHTTPClient{}

	u := &Updater{
		client: &mockSequentialHTTPClient{
			responses: []*http.Response{
				{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "v2.0.0", "assets": []}`)),
				},
				{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`{"tag_name": "v2.0.0", "assets": []}`)),
				},
			},
			callCount: &callCount,
		},
		currentVer:  "1.0.0",
		releasesURL: "https://api.test.com/releases/latest",
	}
	_ = mockClient

	err := u.SelfUpdate()
	if err == nil {
		t.Error("SelfUpdate() should return error when no asset found")
	}
}

type mockSequentialHTTPClient struct {
	responses []*http.Response
	callCount *int
}

func (m *mockSequentialHTTPClient) Get(url string) (*http.Response, error) {
	idx := *m.callCount
	*m.callCount++
	if idx < len(m.responses) {
		return m.responses[idx], nil
	}
	return nil, io.EOF
}
