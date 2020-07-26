package auth

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	goGitSSH "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/hsnodgrass/pufctl/internal/logging"
)

const sshKeyPassphrasePrompt = "Please enter your SSH key passphrase"

// Auth holds a transport.AuthMethod used by go-git
type Auth struct {
	Method transport.AuthMethod
}

// GetAuth returns a pointer to Auth
func (a *Auth) GetAuth() *Auth {
	return a
}

type notHTTPError struct {
	url string
}

func (e *notHTTPError) Error() string {
	return fmt.Sprintf("URL %s is not an HTTP URL", e.url)
}

type notGitError struct {
	url string
}

func (e *notGitError) Error() string {
	return fmt.Sprintf("URL %s is not a Git URL", e.url)
}

type sshKeyReadError struct {
	keyPath string
	err     error
}

func (e *sshKeyReadError) Error() string {
	return fmt.Sprintf("Failed to read SSH key %s: %s", e.keyPath, e.err.Error())
}

type sshKeyPasswordError struct {
	keyPath string
	err     error
}

func (e *sshKeyPasswordError) Error() string {
	return fmt.Sprintf("Password failed for SSH key %s: %s", e.keyPath, e.err.Error())
}

type authError struct {
	url        string
	user       string
	sshKeyPath string
	httpErr    error
	sshErr     error
}

func (e *authError) SetHTTPErr(err error) {
	e.httpErr = err
}

func (e *authError) SetSSHErr(err error) {
	e.sshErr = err
}

func (e *authError) Error() string {
	if e.httpErr != nil && e.sshErr == nil {
		return fmt.Sprintf("Failed to create Git HTTP Auth - URL: %s, User: %s, Error: %#v", e.url, e.user, e.httpErr)
	} else if e.sshErr != nil && e.httpErr == nil {
		return fmt.Sprintf("Failed to create Git SSH Auth - URL: %s, SSH Key: %s, Error: %#v", e.url, e.user, e.sshErr)
	}
	return fmt.Sprintf("GENERIC: Failed to create Git Auth - URL: %s, User: %s, SSH Key: %s, HTTP Error: %#v, SSH Error: %#v", e.url, e.user, e.sshKeyPath, e.httpErr, e.sshErr)
}

// GitAuth parses a given URL and returns the appropriate transport.AuthMethod.
// The first return value of GitAuth is a boolean that's true if the URL is a
// properly formed git URL.
func GitAuth(url, user, pass, token, sshKeyPath string) (Auth, error) {
	genErr := authError{url, user, sshKeyPath, nil, nil}
	hAuth, hErr := httpAuth(url, user, pass, token)
	if hErr == nil {
		return hAuth, nil
	}
	genErr.SetHTTPErr(hErr)
	sAuth, sErr := sshAuth(url, sshKeyPath)
	if sErr == nil {
		return sAuth, nil
	}
	genErr.SetSSHErr(sErr)
	return Auth{}, &genErr
}

func httpAuth(url, user, pass, token string) (Auth, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		logging.Debugln("Detected HTTP(S) URL, proceeding with HTTP(S) authentication strategy")
		if token != "" {
			logging.Debugln("Using HTTP(S) token authentication")
			return Auth{Method: &http.BasicAuth{Username: user, Password: token}}, nil
		} else if pass != "" {
			logging.Debugln("Using HTTP(S) password authentication")
			return Auth{Method: &http.BasicAuth{Username: user, Password: pass}}, nil
		}
		logging.Debugln("Token and password not set, using no authentication")
		return Auth{}, nil
	}
	return Auth{}, &notHTTPError{url}
}

func sshAuth(url, sshKeyPath string) (Auth, error) {
	var auth transport.AuthMethod
	if strings.HasPrefix(url, "ssh://") || strings.HasPrefix(url, "git@") {
		// Check for a running SSH agent before any other auth methods.
		// If a SSH agent is running, return a blank auth struct to
		// offload authentication to the SSH agent itself.
		if sshAgent() {
			logging.Debugln("SSH Agent detected, deferring authentication to SSH agent")
			return Auth{}, nil
		}
		var signer ssh.Signer
		sshKey, err := ioutil.ReadFile(sshKeyPath)
		if err != nil {
			return Auth{}, &sshKeyReadError{sshKeyPath, err}
		}
		signer, err = ssh.ParsePrivateKey([]byte(sshKey))
		if err != nil {
			logging.Warnln("Unable to parse SSH key, trying with password prompt")
		}
		passphrase := promptForPassword(sshKeyPassphrasePrompt)
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(sshKey), passphrase)
		if err != nil {
			return Auth{}, &sshKeyPasswordError{sshKeyPath, err}
		}
		auth = &goGitSSH.PublicKeys{User: "git", Signer: signer}
		return Auth{Method: auth}, nil
	}
	return Auth{}, &notGitError{url}
}

func sshAgent() bool {
	if _, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return true
	}
	return false
}

func promptForPassword(prompt string) []byte {
	fmt.Printf("%s  ", prompt)
	password, err := terminal.ReadPassword(0)
	if err != nil {
		logging.Errorln(err)
	}
	return password
}

// The below interfaces aren't implemented yet
// They will be used eventually to the functions in this
// package to use dependency injection for proper test support.

// not used
type authMethodHandler interface {
	Name() string
	fmt.Stringer
}

// not used
type sshSignerHandler interface {
	PublicKey() ssh.PublicKey
	Sign(io.Reader, []byte) (*ssh.Signature, error)
}

// not used
type sshKeyHandler interface {
	ParsePrivateKey([]byte) (ssh.Signer, error)
	ParsePrivateKeyWithPassphrase([]byte, string) (ssh.Signer, error)
}

// not used
func sshPublicKeys(signer sshSignerHandler) {}
