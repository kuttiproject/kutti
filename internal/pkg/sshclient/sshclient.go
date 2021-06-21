package sshclient

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/containerd/console"
	"github.com/kuttiproject/kuttilog"
	"golang.org/x/crypto/ssh"
)

// SSHClient represents a simple SSH client, which uses password
// authentication.
type SSHClient struct {
	config *ssh.ClientConfig
}

// RunWithResults connects to the specified address, runs the specified command, and
// fetches the results.
func (sc *SSHClient) RunWithResults(address string, command string) (string, error) {
	client, err := ssh.Dial("tcp", address, sc.config)
	if err != nil {
		return "", fmt.Errorf("could not connect to address %s:%v ", address, err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("could not create session at address %s:%v ", address, err)
	}
	defer session.Close()

	resultdata, err := session.Output(command)
	if err != nil {
		return string(resultdata), fmt.Errorf("command '%s' at address %s produced an error:%v ", command, address, err)
	}

	return string(resultdata), nil
}

// RunInterativeShell connects to the specified address and runs an interactive
// shell.
func (sc *SSHClient) RunInterativeShell(address string) {
	// Handle OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	// Start an interactive shell
	go func() {
		if err := sc.runclient(ctx, address); err != nil {
			kuttilog.Print(kuttilog.Quiet, err)
		}
		cancel()
	}()

	// Wait for either OS interrupt or shell finish
	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}

// Copied almost verbatim from https://gist.github.com/atotto/ba19155295d95c8d75881e145c751372
// Thanks, Ato Araki (atotto@github)
func (sc *SSHClient) runclient(ctx context.Context, address string) error {
	conn, err := ssh.Dial("tcp", address, sc.config)
	if err != nil {
		return fmt.Errorf("cannot connect to %v: %v", address, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %v", err)
	}
	defer session.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	// Set console to raw mode for proper vtty behaviour.
	// The original used Terminal.MakeRaw. After seeing
	// serious problems with that, containerd/console was
	// used instead.
	current := console.Current()
	defer current.Reset()

	err = current.SetRaw()
	if err != nil {
		return fmt.Errorf("while terminal make raw: %s", err)
	}

	ws, err := current.Size()
	if err != nil {
		return fmt.Errorf("while terminal get size: %s", err)
	}

	h := int(ws.Height)
	w := int(ws.Width)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	if err := session.RequestPty(term, h, w, modes); err != nil {
		return fmt.Errorf("while requesting session xterm: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("while requesting session shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("in ssh: %s", err)
	}
	return nil
}

// NewWithPassword creates a new SSH client with password authentication, and no host key check
func NewWithPassword(username string, password string) *SSHClient {
	return &SSHClient{
		config: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			Timeout:         5 * time.Second,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
}
