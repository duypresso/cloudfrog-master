package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"f/cloudfrog/backend/config"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPConnect establishes connection to SFTP server
func SFTPConnect() (*sftp.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.SFTPUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SFTPPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Connect directly using format: user@hostname with port
	addr := fmt.Sprintf("%s@%s", config.SFTPUser, config.SFTPHost)
	log.Printf("Connecting to SFTP server at %s:%s", addr, config.SFTPPort)

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.SFTPHost, config.SFTPPort), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP server: %v", err)
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %v", err)
	}

	// Verify connection
	log.Printf("Verifying SFTP connection to %s", config.SFTPPath)
	_, err = sftpClient.ReadDir(config.SFTPPath)
	if err != nil {
		sftpClient.Close()
		sshClient.Close()
		return nil, fmt.Errorf("failed to access SFTP directory %s: %v", config.SFTPPath, err)
	}

	return sftpClient, nil
}

// UploadToSFTP uploads a file to SFTP server
func UploadToSFTP(localFile *os.File, remotePath string) error {
	sftpClient, err := SFTPConnect()
	if err != nil {
		return fmt.Errorf("SFTP connection failed: %v", err)
	}
	defer sftpClient.Close()

	// Convert Windows path separators to Unix style
	remotePath = filepath.ToSlash(remotePath)
	log.Printf("Uploading file to SFTP path: %s", remotePath)

	// Check if directory exists first
	dirPath := filepath.ToSlash(filepath.Dir(remotePath))
	info, err := sftpClient.Stat(dirPath)
	if err != nil {
		log.Printf("Directory %s not found, attempting to create", dirPath)
		err = sftpClient.MkdirAll(dirPath)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", dirPath)
	}

	// Create remote file
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file (%s): %v", remotePath, err)
	}
	defer remoteFile.Close()

	// Reset file pointer and copy contents
	_, err = localFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %v", err)
	}

	written, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	log.Printf("Successfully uploaded %d bytes to %s", written, remotePath)
	return nil
}

// DeleteFromSFTP deletes a file from SFTP server
func DeleteFromSFTP(remotePath string) error {
	sftpClient, err := SFTPConnect()
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	// Delete the file
	err = sftpClient.Remove(remotePath)
	if err != nil {
		return fmt.Errorf("failed to delete remote file: %v", err)
	}

	return nil
}

// DownloadFromSFTP downloads a file from SFTP server and writes it to the provided writer
func DownloadFromSFTP(remotePath string) (io.ReadCloser, error) {
	sftpClient, err := SFTPConnect()
	if err != nil {
		log.Printf("SFTP connection failed: %v", err)
		return nil, fmt.Errorf("SFTP connection failed: %v", err)
	}

	// Convert path to Unix style and ensure it starts with /
	remotePath = filepath.ToSlash(remotePath)
	if !strings.HasPrefix(remotePath, "/") {
		remotePath = "/" + remotePath
	}

	log.Printf("Attempting to download file from: %s", remotePath)

	// Open remote file
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		sftpClient.Close()
		log.Printf("Failed to open remote file %s: %v", remotePath, err)
		return nil, fmt.Errorf("failed to open remote file: %v", err)
	}

	log.Printf("Successfully opened file: %s", remotePath)

	// Create a wrapper that will close both the file and client
	return &sftpReadCloser{
		Reader: remoteFile,
		closeFunc: func() error {
			remoteFile.Close()
			return sftpClient.Close()
		},
	}, nil
}

// sftpReadCloser wraps an SFTP file to ensure both file and client are closed
type sftpReadCloser struct {
	io.Reader
	closeFunc func() error
}

func (s *sftpReadCloser) Close() error {
	return s.closeFunc()
}
