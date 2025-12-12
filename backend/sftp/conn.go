// Package sftp: connections sftp
package sftp

import (
	"fmt"
	"log"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func InitSFTP(host, port, username string, userPrivateKey *ssh.Signer, hostKey *ssh.PublicKey) (*sftp.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("[i] Tentative de connexion à SFTP: host=%s port=%s user=%s\n", host, port, username)

	var auths []ssh.AuthMethod
	auths = append(auths, ssh.PublicKeys(*userPrivateKey))

	sshConfig := &ssh.ClientConfig{
		User:            username,
		Auth:            auths,
		HostKeyCallback: ssh.FixedHostKey(*hostKey),
	}

	sshClient, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("[!] impossible de créer la connection ssh avec %s : %v", addr, err)
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		sftpClient.Close()
		return nil, fmt.Errorf("[!] échec de la création du client SFTP : %s", err)
	}

	return sftpClient, nil
}
