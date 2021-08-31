package ssh

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

	"golang.org/x/crypto/ssh"
)

func ParseKeyFile(keyFilePath string) ([]ssh.AuthMethod, error) {
	buf, err := readKeyFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	} // ~が読めないので絶対パスに変換してください
	k, _ := ssh.ParsePrivateKey(buf)
	return []ssh.AuthMethod{ssh.PublicKeys(k)}, nil
}

func readKeyFile(keyFilePath string) ([]byte, error) {
	if strings.Contains(keyFilePath, "~") {
		usr, _ := user.Current()
		keyFilePath = strings.Replace(keyFilePath, "~", usr.HomeDir, 1)
	}
	// use assets
	bb, err := asset(keyFilePath[1:])
	if err == nil {
		return bb, nil
	}
	// fallback to read file system
	return ioutil.ReadFile(keyFilePath)
}

func asset(name string) ([]byte, error) {
	return nil, fmt.Errorf("asset %s not found", name)
}
