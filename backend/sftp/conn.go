// Package sftp: connections sftp
package sftp

import (
	"fmt"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func InitSFTP(host, port, username, password string) (*sftp.Client, error) {
	addr := host + ":" + port

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		sshClient.Close()
		return nil, fmt.Errorf("[!] Impossible de créer la connection ssh avec %s : %v", addr, err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		sftpClient.Close()
		return nil, fmt.Errorf("[!] Échec de la création du client SFTP : %s", err)
	}

	return sftpClient, nil
}
