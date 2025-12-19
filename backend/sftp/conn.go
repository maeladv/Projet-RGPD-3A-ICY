// Package sftp: connections sftp
package sftp

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func InitSFTP(host, username string, userPrivateKey *ssh.Signer) (*sftp.Client, error) {
	log.Printf("[i] Tentative de connexion à SFTP: host=%s user=%s\n", host, username)

	var auths []ssh.AuthMethod
	auths = append(auths, ssh.PublicKeys(*userPrivateKey))

	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var sshClient *ssh.Client
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		sshClient, err = ssh.Dial("tcp", host, sshConfig)
		if err == nil {
			break
		}

		log.Printf("[!] Attempt %d/%d failed: %v", attempt, 10, err)
		if attempt < 10 {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("[!] impossible de créer la connection ssh avec %s : %v", host, err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		sftpClient.Close()
		return nil, fmt.Errorf("[!] échec de la création du client SFTP : %s", err)
	}

	return sftpClient, nil
}
