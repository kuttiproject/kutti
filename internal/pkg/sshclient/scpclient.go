package sshclient

import (
	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
)

func copyFileToRemote(config *ssh.ClientConfig, address string, localFilePath string, remoteFilePath string) error {
	scpClient, err := scp.NewClient(address, config, &scp.ClientOption{})
	if err != nil {
		return err
	}

	return scpClient.CopyFileToRemote(localFilePath, remoteFilePath, &scp.FileTransferOption{})
}

func copyFileFromRemote(config *ssh.ClientConfig, address string, remoteFilePath string, localFilePath string) error {
	scpClient, err := scp.NewClient(address, config, &scp.ClientOption{})
	if err != nil {
		return err
	}

	return scpClient.CopyFileFromRemote(remoteFilePath, localFilePath, &scp.FileTransferOption{})
}

func copyDirectoryToRemote(config *ssh.ClientConfig, address string, localPath string, remotePath string) error {
	scpClient, err := scp.NewClient(address, config, &scp.ClientOption{})
	if err != nil {
		return err
	}

	return scpClient.CopyDirToRemote(localPath, remotePath, &scp.DirTransferOption{})
}

func copyDirectoryFromRemote(config *ssh.ClientConfig, address string, remotePath string, localPath string) error {
	scpClient, err := scp.NewClient(address, config, &scp.ClientOption{})
	if err != nil {
		return err
	}

	return scpClient.CopyDirFromRemote(remotePath, localPath, &scp.DirTransferOption{})
}
