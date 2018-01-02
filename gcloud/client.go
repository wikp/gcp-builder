package gcloud

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const installerLocation = "https://dl.google.com/dl/cloudsdk/release/install_google_cloud_sdk.bash"
const installerScriptLocation = "./install_google_cloud_sdk.bash"

type Client struct {
	log *log.Logger
}

func NewClient() *Client {
	return &Client{log.New(
		os.Stdout, "[gcloud client] ", log.Lmicroseconds,
	)}
}

func (i *Client) InstallDir() string {
	if dir, err := homedir.Dir(); err == nil {
		return dir + "/.tmp-travis-builder"
	} else {
		return "~/.tmp-travis-builder"
	}
}

func (i *Client) IsInstalled(command string) bool {
	if _, err := os.Stat(i.sdkBinaryLocation(command)); os.IsNotExist(err) {
		return false
	}

	return true
}

func (i *Client) Install() error {

	if i.IsInstalled("gcloud") {
		i.log.Printf("Updating Google Cloud Platform SDK components from %s", installerLocation)

		if err := i.RunCommand("gcloud", []string{"components", "update"}); err != nil {
			return err
		}
	} else {
		i.log.Printf("Downloading Google Cloud Platform SDK installer from %s", installerLocation)

		response, err := http.Get(installerLocation)
		if err != nil {
			return err
		}

		defer response.Body.Close()

		file, err := os.Create(installerScriptLocation)
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, response.Body); err != nil {
			return err
		}

		if err := file.Chmod(0760); err != nil {
			return err
		}

		if err := file.Close(); err != nil {
			return err
		}

		cmd := exec.Command(installerScriptLocation)
		cmd.Env = []string{
			"CLOUDSDK_CORE_DISABLE_PROMPTS=1",
			fmt.Sprintf("CLOUDSDK_INSTALL_DIR=%s", i.InstallDir()),
			fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		}

		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			i.log.Printf("There was an error while installing SDK: %s", err)
			return err
		} else {
			i.log.Printf("Google Cloud Platform SDK installed")
		}
	}

	if !i.IsInstalled("kubectl") {
		i.log.Printf("Installing kubectl")

		if output, err := i.CaptureCommand("gcloud", []string{"components", "install", "kubectl"}); err != nil {
			i.log.Printf("Could not install kubectl due to: %s", err)
			return err
		} else {
			i.log.Printf("Kubectl was installed: %s", output)
		}
	}

	if !i.IsInstalled("kubectl") || !i.IsInstalled("gcloud") {
		return errors.New("InstallationFailed")
	}

	return nil
}

func (i *Client) CaptureCommand(command string, args []string) (string, error) {
	cmd := exec.Command(i.sdkBinaryLocation(command), args...)

	i.log.Printf("Running command %s %+v", command, args)

	cmd.Env = []string{
		fmt.Sprintf("PATH=%s:%s", i.sdkBinaryLocation(""), os.Getenv("PATH")),
		fmt.Sprintf("KUBECONFIG=%s/.kube", i.InstallDir()),
	}

	out := bytes.Buffer{}

	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return string(out.Bytes()), nil
}

func (i *Client) RunCommand(command string, args []string) error {
	cmd := exec.Command(i.sdkBinaryLocation(command), args...)

	i.log.Printf("Running command %s %+v", command, args)

	cmd.Env = []string{
		fmt.Sprintf("PATH=%s:%s", i.sdkBinaryLocation(""), os.Getenv("PATH")),
		fmt.Sprintf("KUBECONFIG=%s/.kube", i.InstallDir()),
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (i *Client) sdkBinaryLocation(command string) string {
	return fmt.Sprintf("%s/google-cloud-sdk/bin/%s", i.InstallDir(), command)
}
