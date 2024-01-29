package setup

import (
	"fmt"
	"strings"

	commonNs "github.com/longhorn/go-common-libs/ns"

	"github.com/longhorn/go-spdk-helper/pkg/types"
)

const (
	spdkSetupPath = "/usr/src/spdk/scripts/setup.sh"
)

func Bind(deviceAddr, deviceDriver string, executor *commonNs.Executor) (string, error) {
	if deviceAddr == "" {
		return "", fmt.Errorf("device address is empty")
	}

	envs := []string{
		fmt.Sprintf("%s=%s", "PCI_ALLOWED", deviceAddr),
		fmt.Sprintf("%s=%s", "DRIVER_OVERRIDE", deviceDriver),
	}

	cmdArgs := []string{
		spdkSetupPath,
		"bind",
	}

	outputStr, err := executor.Execute(envs, "bash", cmdArgs, types.ExecuteTimeout)
	if err != nil {
		return "", err
	}

	return outputStr, nil
}

func Unbind(deviceAddr string, executor *commonNs.Executor) (string, error) {
	if deviceAddr == "" {
		return "", fmt.Errorf("device address is empty")
	}

	cmdArgs := []string{
		spdkSetupPath,
		"unbind",
		deviceAddr,
	}

	outputStr, err := executor.Execute(nil, "bash", cmdArgs, types.ExecuteTimeout)
	if err != nil {
		return "", err
	}

	return outputStr, nil
}

func GetDiskDriver(deviceAddr string, executor *commonNs.Executor) (string, error) {
	if deviceAddr == "" {
		return "", fmt.Errorf("device address is empty")
	}

	cmdArgs := []string{
		spdkSetupPath,
		"disk-driver",
		deviceAddr,
	}

	outputStr, err := executor.Execute(nil, "bash", cmdArgs, types.ExecuteTimeout)
	if err != nil {
		return "", err
	}

	return extractJSON(outputStr)
}

func GetDiskStatus(deviceAddr string, executor *commonNs.Executor) (string, error) {
	if deviceAddr == "" {
		return "", fmt.Errorf("device address is empty")
	}

	cmdArgs := []string{
		spdkSetupPath,
		"disk-status",
		deviceAddr,
	}

	outputStr, err := executor.Execute(nil, "bash", cmdArgs, types.ExecuteTimeout)
	if err != nil {
		return "", err
	}

	return extractJSON(outputStr)
}

func extractJSON(outputStr string) (string, error) {
	// Find the first '{' and last '}' characters, assuming valid JSON format
	start := strings.Index(outputStr, "{")
	end := strings.LastIndex(outputStr, "}")
	if start != -1 && end != -1 {
		return outputStr[start : end+1], nil
	}
	return "", fmt.Errorf("failed to extract JSON from output: %s", outputStr)
}
