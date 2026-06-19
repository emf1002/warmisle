// Package plugin defines the standard interfaces for cloud drive adapters.
package plugin

import "time"

// CloudFileInfo holds metadata about a file stored in a cloud drive.
type CloudFileInfo struct {
	FileID     string    `json:"file_id"`
	FileName   string    `json:"file_name"`
	Size       int64     `json:"size"`
	Path       string    `json:"path"`
	CreateTime time.Time `json:"create_time"`
}

// CloudDrive defines the interface for cloud drive file operations.
// Implementations must handle authentication internally via TokenProvider.
type CloudDrive interface {
	// Upload uploads a local file to the specified cloud directory.
	// localPath is the absolute path to the file on disk.
	// cloudDir is the target directory path in the cloud drive.
	// Returns metadata about the uploaded file on success.
	Upload(localPath, cloudDir string) (*CloudFileInfo, error)

	// Download downloads a file from the cloud drive to a local path.
	// fileID is the cloud drive file identifier.
	// localPath is the destination path on disk.
	Download(fileID, localPath string) error

	// List returns a list of files in the specified cloud directory.
	List(cloudDir string) ([]CloudFileInfo, error)

	// Delete removes a file from the cloud drive by its identifier.
	Delete(fileID string) error
}

// TokenProvider defines the interface for OAuth-based token management
// used by cloud drive implementations.
type TokenProvider interface {
	// GetAccessToken returns a valid access token, refreshing it if necessary.
	GetAccessToken() (string, error)

	// GetAuthURL returns the authorization URL the user should visit to grant access.
	// state is an anti-CSRF token that will be returned in the callback.
	GetAuthURL(state string) string

	// ExchangeCode exchanges an authorization code for an access token.
	// This is called upon receiving the OAuth callback.
	ExchangeCode(code string) error
}
