package sftp

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type Client interface {
	Connect() error
	UploadFile(localFileName string, remoteFileName string) error
	Upload(data []byte, remoteFileName string) error
	DownloadFile(fileName string) ([]byte, error)
}

type client struct {
	user       string
	certPath   string
	host       string
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func NewClient(host string, certPath string, user string) Client {
	return &client{
		user:     user,
		host:     host,
		certPath: certPath,
	}
}

func (c *client) Connect() error {
	// get auth method
	key, err := ioutil.ReadFile(c.certPath)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	sshConfig := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer)},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshClient, err := ssh.Dial("tcp", c.host, sshConfig)
	if err != nil {
		return err
	}
	c.sshClient = sshClient

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	c.sftpClient = sftpClient

	return nil
}
func (c *client) Upload(data []byte, remoteFileName string) error {
	err := c.upload(data, remoteFileName)
	return err
}

func (c *client) UploadFile(localFileName string, remoteFileName string) error {
	localFile, err := os.Open(localFileName)
	if err != nil {
		return err
	}
	defer localFile.Close()

	remoteFile, err := c.sftpClient.Create(remoteFileName)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	data := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		readCount, err := localFile.Read(buf)
		if readCount > 0 {
			data = append(data, buf...)
		}
		if err == io.EOF {
			return nil
		}
		return err
	}
	err = c.upload(data, remoteFileName)
	return err
}

func (c *client) DownloadFile(fileName string) ([]byte, error) {
	file, err := c.sftpClient.OpenFile(fileName, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *client) upload(data []byte, remoteFileName string) error {
	remoteFile, err := c.sftpClient.Create(remoteFileName)
	if err != nil {
		return err
	}
	defer remoteFile.Close()
	chunkSize := 1024
	var batchs [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize

		if end > len(data) {
			end = len(data)
		}

		batchs = append(batchs, data[i:end])
	}

	for _, batch := range batchs {
		_, err := remoteFile.Write(batch)
		if err != nil {
			return err
		}
	}
	return nil
}
