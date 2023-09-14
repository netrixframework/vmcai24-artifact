package client

import (
	"os/exec"
	"path"
	"strings"
)

type BFTSmartClientConfig struct {
	CodePath string
}

type BFTSmartClient struct {
	config *BFTSmartClientConfig
}

func NewBFTSmartClient(config *BFTSmartClientConfig) *BFTSmartClient {
	return &BFTSmartClient{
		config: config,
	}
}

func (b *BFTSmartClient) execCmd(args ...string) (string, error) {
	cmdArgs := []string{
		// "-Djava.security.properties=\"" + path.Join(b.config.CodePath, "config/java.security") + "\"",
		// "-Dlogback.configurationFile=\"" + path.Join(b.config.CodePath, "config/logback.xml") + "\"",
		// "-cp", path.Join(b.config.CodePath, "build/install/library/lib/*"),
		path.Join(b.config.CodePath, "runscripts/smartrun.sh"),
		"bftsmart.demo.map.MapArgsClient", "1001",
	}

	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("bash", cmdArgs...)
	cmd.Dir = b.config.CodePath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(output), "\n"), nil
}

func (b *BFTSmartClient) Get(key string) (string, error) {
	return b.execCmd("get", key)
}

func (b *BFTSmartClient) Set(key, value string) (string, error) {
	return b.execCmd("set", key, value)
}

func (b *BFTSmartClient) Remove(key string) (string, error) {
	return b.execCmd("delete", key)
}
