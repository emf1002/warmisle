package alipan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"warmisle/internal/plugin"
)

const alipanDriveAPI = "https://openapi.alipan.com/adrive/v1.0"

// Ensure AlipanClient implements plugin.CloudDrive at compile time.
var _ plugin.CloudDrive = (*AlipanClient)(nil)

// AlipanClient is the Aliyun Drive API client implementing the CloudDrive interface.
type AlipanClient struct {
	oauth      *AlipanOAuth
	httpClient *http.Client
}

// NewAlipanClient creates a new AlipanClient with the given OAuth provider.
func NewAlipanClient(oauth *AlipanOAuth) *AlipanClient {
	return &AlipanClient{
		oauth:      oauth,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// setAuthHeader sets the Authorization Bearer header on the request.
func (c *AlipanClient) setAuthHeader(req *http.Request) error {
	token, err := c.oauth.GetAccessToken()
	if err != nil {
		return fmt.Errorf("获取访问令牌失败: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

// getDriveID retrieves the user's default drive ID.
func (c *AlipanClient) getDriveID() (string, error) {
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/user/getDriveInfo", alipanDriveAPI), nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	if err := c.setAuthHeader(req); err != nil {
		return "", err
	}

	log.Printf("[alipan] POST %s/user/getDriveInfo", alipanDriveAPI)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取网盘信息失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取网盘信息响应失败: %w", err)
	}

	log.Printf("[alipan] getDriveInfo response status=%d body=%s", resp.StatusCode, string(respBytes))

	var result struct {
		DriveID string `json:"default_drive_id"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("解析网盘信息失败: %w", err)
	}

	if result.DriveID == "" {
		return "", fmt.Errorf("无法获取默认网盘ID，响应: %s", string(respBytes))
	}

	return result.DriveID, nil
}

// resolveDirectory ensures the target cloud directory exists and returns its file_id.
// If the directory does not exist it will be created.
func (c *AlipanClient) resolveDirectory(driveID, cloudDir string) (string, error) {
	// First, try to find existing directory by searching
	fileID, err := c.searchDirectory(driveID, cloudDir)
	if err == nil && fileID != "" {
		return fileID, nil
	}

	// Create directory under root
	log.Printf("[alipan] directory '%s' not found, creating...", cloudDir)

	// Strip leading and trailing slashes for the folder name
	dirName := strings.Trim(cloudDir, "/")

	reqBody := map[string]interface{}{
		"drive_id":        driveID,
		"parent_file_id":  "root",
		"name":            dirName,
		"type":            "folder",
		"check_name_mode": "refuse",
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/create", alipanDriveAPI),
		bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建目录请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(req); err != nil {
		return "", err
	}

	log.Printf("[alipan] POST %s/openFile/create (folder)", alipanDriveAPI)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("创建目录请求失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[alipan] create folder response status=%d body=%s", resp.StatusCode, string(respBytes))

	var result struct {
		FileID string `json:"file_id"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("解析创建目录响应失败: %w", err)
	}

	return result.FileID, nil
}

// searchDirectory searches for a directory by name and returns its file_id.
func (c *AlipanClient) searchDirectory(driveID, cloudDir string) (string, error) {
	reqBody := map[string]interface{}{
		"drive_id": driveID,
		"query":    fmt.Sprintf("name = '%s' and type = 'folder'", cloudDir),
		"limit":    1,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/search", alipanDriveAPI),
		bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建搜索请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(req); err != nil {
		return "", err
	}

	log.Printf("[alipan] POST %s/openFile/search", alipanDriveAPI)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("搜索目录请求失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[alipan] search response status=%d body=%s", resp.StatusCode, string(respBytes))

	var result struct {
		Items []struct {
			FileID string `json:"file_id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("解析搜索响应失败: %w", err)
	}

	if len(result.Items) > 0 {
		return result.Items[0].FileID, nil
	}
	return "", fmt.Errorf("目录未找到")
}

// Upload uploads a local file to the specified cloud directory.
func (c *AlipanClient) Upload(localPath, cloudDir string) (*plugin.CloudFileInfo, error) {
	// Get drive ID
	driveID, err := c.getDriveID()
	if err != nil {
		return nil, fmt.Errorf("上传失败: %w", err)
	}

	// Resolve parent directory
	parentFileID, err := c.resolveDirectory(driveID, cloudDir)
	if err != nil {
		return nil, fmt.Errorf("解析云端目录失败: %w", err)
	}
	if parentFileID == "" {
		parentFileID = "root"
	}

	// Get local file info
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return nil, fmt.Errorf("读取本地文件失败: %w", err)
	}
	fileSize := fileInfo.Size()
	fileName := fileInfo.Name()

	log.Printf("[alipan] uploading %s (%d bytes) to dir %s (parent=%s)", localPath, fileSize, cloudDir, parentFileID)

	// Step 1: Pre-upload — create file entry
	createBody := map[string]interface{}{
		"drive_id":         driveID,
		"parent_file_id":   parentFileID,
		"name":             fileName,
		"type":             "file",
		"size":             fileSize,
		"check_name_mode":  "auto_rename",
		"pre_hash":         "",
	}
	createBytes, _ := json.Marshal(createBody)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/create", alipanDriveAPI),
		bytes.NewReader(createBytes))
	if err != nil {
		return nil, fmt.Errorf("创建上传请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(req); err != nil {
		return nil, fmt.Errorf("上传失败: %w", err)
	}

	log.Printf("[alipan] POST %s/openFile/create (file pre-upload)", alipanDriveAPI)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("预上传请求失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[alipan] pre-upload response status=%d body=%s", resp.StatusCode, string(respBytes))

	type partInfo struct {
		PartNumber int    `json:"part_number"`
		UploadURL  string `json:"upload_url"`
	}

	var createResult struct {
		RapidUpload  bool       `json:"rapid_upload"`
		FileID       string     `json:"file_id"`
		UploadID     string     `json:"upload_id"`
		PartInfoList []partInfo `json:"part_info_list"`
	}

	if err := json.Unmarshal(respBytes, &createResult); err != nil {
		return nil, fmt.Errorf("解析预上传响应失败: %w", err)
	}

	// If rapid_upload is true, the file already exists on the server
	if createResult.RapidUpload {
		log.Printf("[alipan] rapid upload succeeded, file_id=%s", createResult.FileID)
		return &plugin.CloudFileInfo{
			FileID:     createResult.FileID,
			FileName:   fileName,
			Size:       fileSize,
			Path:       cloudDir + "/" + fileName,
			CreateTime: time.Now(),
		}, nil
	}

	// Step 2: Upload file content if needed
	if len(createResult.PartInfoList) == 0 {
		// Some responses may not have part_info_list — the file entry is created directly
		log.Printf("[alipan] file entry created, file_id=%s", createResult.FileID)
		return &plugin.CloudFileInfo{
			FileID:     createResult.FileID,
			FileName:   fileName,
			Size:       fileSize,
			Path:       cloudDir + "/" + fileName,
			CreateTime: time.Now(),
		}, nil
	}

	// Upload each part
	file, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("打开本地文件失败: %w", err)
	}
	defer func() { _ = file.Close() }()

	for _, part := range createResult.PartInfoList {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("定位文件失败: %w", err)
		}

		partReq, err := http.NewRequest(http.MethodPut, part.UploadURL, file)
		if err != nil {
			return nil, fmt.Errorf("创建分片上传请求失败: %w", err)
		}

		log.Printf("[alipan] PUT part %d to %s", part.PartNumber, part.UploadURL)

		partResp, err := c.httpClient.Do(partReq)
		if err != nil {
			return nil, fmt.Errorf("分片上传失败 (part %d): %w", part.PartNumber, err)
		}
		_ = partResp.Body.Close()

		if partResp.StatusCode >= 400 {
			return nil, fmt.Errorf("分片上传返回错误状态 %d (part %d)", partResp.StatusCode, part.PartNumber)
		}

		log.Printf("[alipan] part %d uploaded successfully", part.PartNumber)
	}

	// Step 3: Complete the upload
	completeBody := map[string]interface{}{
		"drive_id":  driveID,
		"file_id":   createResult.FileID,
		"upload_id": createResult.UploadID,
	}
	completeBytes, _ := json.Marshal(completeBody)

	compReq, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/complete", alipanDriveAPI),
		bytes.NewReader(completeBytes))
	if err != nil {
		return nil, fmt.Errorf("创建完成上传请求失败: %w", err)
	}
	compReq.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(compReq); err != nil {
		return nil, fmt.Errorf("上传确认失败: %w", err)
	}

	log.Printf("[alipan] POST %s/openFile/complete", alipanDriveAPI)

	compResp, err := c.httpClient.Do(compReq)
	if err != nil {
		return nil, fmt.Errorf("完成上传请求失败: %w", err)
	}
	defer func() { _ = compResp.Body.Close() }()

	compRespBytes, _ := io.ReadAll(compResp.Body)
	log.Printf("[alipan] complete upload response status=%d body=%s", compResp.StatusCode, string(compRespBytes))

	if compResp.StatusCode >= 400 {
		return nil, fmt.Errorf("完成上传返回错误状态 %d: %s", compResp.StatusCode, string(compRespBytes))
	}

	log.Printf("[alipan] upload completed, file_id=%s", createResult.FileID)
	return &plugin.CloudFileInfo{
		FileID:     createResult.FileID,
		FileName:   fileName,
		Size:       fileSize,
		Path:       cloudDir + "/" + fileName,
		CreateTime: time.Now(),
	}, nil
}

// Download downloads a file from the cloud drive to a local path.
func (c *AlipanClient) Download(fileID, localPath string) error {
	driveID, err := c.getDriveID()
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	// Step 1: Get download URL
	dlBody := map[string]interface{}{
		"drive_id": driveID,
		"file_id":  fileID,
	}
	dlBytes, _ := json.Marshal(dlBody)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/getDownloadUrl", alipanDriveAPI),
		bytes.NewReader(dlBytes))
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(req); err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}

	log.Printf("[alipan] POST %s/openFile/getDownloadUrl", alipanDriveAPI)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("获取下载URL失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[alipan] getDownloadUrl response status=%d body=%s", resp.StatusCode, string(respBytes))

	var dlResult struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(respBytes, &dlResult); err != nil {
		return fmt.Errorf("解析下载URL失败: %w", err)
	}

	if dlResult.URL == "" {
		return fmt.Errorf("获取下载URL失败，响应: %s", string(respBytes))
	}

	// Step 2: Download the file
	dlReq, err := http.NewRequest(http.MethodGet, dlResult.URL, nil)
	if err != nil {
		return fmt.Errorf("创建下载请求失败: %w", err)
	}

	log.Printf("[alipan] GET %s", dlResult.URL)

	dlResp, err := c.httpClient.Do(dlReq)
	if err != nil {
		return fmt.Errorf("下载文件失败: %w", err)
	}
	defer func() { _ = dlResp.Body.Close() }()

	if dlResp.StatusCode >= 400 {
		return fmt.Errorf("下载文件返回错误状态 %d", dlResp.StatusCode)
	}

	// Step 3: Write to local file
	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer func() { _ = outFile.Close() }()

	written, err := io.Copy(outFile, dlResp.Body)
	if err != nil {
		return fmt.Errorf("写入本地文件失败: %w", err)
	}

	log.Printf("[alipan] downloaded %d bytes to %s", written, localPath)
	return nil
}

// List returns a list of files in the specified cloud directory.
// It navigates from root by name to find the target backup directory.
func (c *AlipanClient) List(cloudDir string) ([]plugin.CloudFileInfo, error) {
	driveID, err := c.getDriveID()
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %w", err)
	}

	dirName := strings.Trim(cloudDir, "/")
	log.Printf("[alipan] List: driveID=%s, cloudDir=%s, dirName=%s", driveID, cloudDir, dirName)

	// Step 1: list root to find backup folder
	rootFiles, listErr := c.listFiles(driveID, "root")
	if listErr != nil {
		log.Printf("[alipan] failed to list root: %v, falling back to root listing", listErr)
		return c.listFiles(driveID, "root")
	}

	log.Printf("[alipan] root has %d entries, looking for '%s'", len(rootFiles), dirName)
	for i, f := range rootFiles {
		cleanName := strings.TrimLeft(f.FileName, "/")
		log.Printf("[alipan] root[%d]: name=%q file_id=%q size=%d", i, f.FileName, f.FileID, f.Size)
		if f.FileName == dirName || cleanName == dirName {
			log.Printf("[alipan] FOUND backup folder: file_id=%s (raw name: %s)", f.FileID, f.FileName)
			return c.listFiles(driveID, f.FileID)
		}
	}

	log.Printf("[alipan] backup folder '%s' not found in root", dirName)
	return nil, fmt.Errorf("备份目录 %s 未找到", dirName)
}

// listFiles lists files in a directory by its file_id.
func (c *AlipanClient) listFiles(driveID, parentFileID string) ([]plugin.CloudFileInfo, error) {
	listBody := map[string]interface{}{
		"drive_id":        driveID,
		"parent_file_id":  parentFileID,
		"limit":           100,
		"order_by":        "updated_at",
		"order_direction": "DESC",
	}
	listBytes, _ := json.Marshal(listBody)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/list", alipanDriveAPI),
		bytes.NewReader(listBytes))
	if err != nil {
		return nil, fmt.Errorf("创建列表请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(req); err != nil {
		return nil, fmt.Errorf("列出文件失败: %w", err)
	}

	log.Printf("[alipan] POST %s/openFile/list (parent=%s)", alipanDriveAPI, parentFileID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("列出文件请求失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[alipan] list response status=%d body=%s", resp.StatusCode, string(respBytes))

	var listResult struct {
		Items []struct {
			FileID    string `json:"file_id"`
			Name      string `json:"name"`
			FileName  string `json:"file_name"`
			Size      int64  `json:"size"`
			Type      string `json:"type"`
			CreatedAt string `json:"created_at"`
		} `json:"items"`
	}

	if err := json.Unmarshal(respBytes, &listResult); err != nil {
		return nil, fmt.Errorf("解析文件列表失败: %w", err)
	}

	infos := make([]plugin.CloudFileInfo, 0, len(listResult.Items))
	for _, item := range listResult.Items {
		name := item.Name
		if name == "" {
			name = item.FileName
		}
		createTime, _ := time.Parse(time.RFC3339, item.CreatedAt)
		infos = append(infos, plugin.CloudFileInfo{
			FileID:     item.FileID,
			FileName:   name,
			Size:       item.Size,
			Path:       name,
			CreateTime: createTime,
		})
	}
	return infos, nil
}

// Delete removes a file from the cloud drive by its file identifier.
func (c *AlipanClient) Delete(fileID string) error {
	driveID, err := c.getDriveID()
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	delBody := map[string]interface{}{
		"drive_id": driveID,
		"file_id":  fileID,
	}
	delBytes, _ := json.Marshal(delBody)

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/openFile/delete", alipanDriveAPI),
		bytes.NewReader(delBytes))
	if err != nil {
		return fmt.Errorf("创建删除请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := c.setAuthHeader(req); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	log.Printf("[alipan] POST %s/openFile/delete", alipanDriveAPI)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("删除文件请求失败: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[alipan] delete response status=%d body=%s", resp.StatusCode, string(respBytes))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("删除文件返回错误状态 %d: %s", resp.StatusCode, string(respBytes))
	}

	return nil
}
