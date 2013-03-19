package commands

import (
	"code.google.com/p/go.crypto/openpgp"
	"fmt"
	"github.com/jgrocho/passphrase"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

func encrypt(ringPath string, input io.Reader, output io.Writer) error {
	ringFile, err := os.Open(ringPath)
	if err != nil {
		return err
	}
	defer ringFile.Close()

	ring, err := openpgp.ReadKeyRing(ringFile)
	if err != nil {
		return err
	}
	// TODO: Select proper key
	key := ring[0]

	plaintext, err := openpgp.Encrypt(output, []*openpgp.Entity{key}, nil, nil, nil)
	if err != nil {
		return err
	}
	defer plaintext.Close()

	if _, err := io.Copy(plaintext, input); err != nil {
		return err
	}

	return nil
}

func decrypt(ringPath string, input io.Reader) (io.Reader, error) {
	ringFile, err := os.Open(ringPath)
	if err != nil {
		return nil, err
	}
	defer ringFile.Close()

	ring, err := openpgp.ReadKeyRing(ringFile)
	if err != nil {
		return nil, err
	}

	var keyToTry, attempt int
	var triedCache bool
	promptFunc := openpgp.PromptFunction(func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if keyToTry >= len(keys) {
			return nil, fmt.Errorf("no more keys to try")
		}
		if attempt > 2 {
			attempt = 0
			keyToTry++
			return nil, nil
		}
		defer func() { attempt++ }()

		key := keys[keyToTry]
		fingerprint := fmt.Sprintf("%X", key.PublicKey.Fingerprint)

		if !triedCache {
			triedCache = true
			if cachedPass, _ := passphrase.GetPassphrase(fingerprint, "", "", "", false, false); cachedPass != "" {
				if err := key.PrivateKey.Decrypt([]byte(cachedPass)); err == nil {
					return nil, nil
				}
			}
		}

		passphrase.ClearCachedPassphrase(fingerprint)
		prompt := ""
		description := fmt.Sprintf("Key %s; attempt %d", key.PublicKey.KeyIdShortString(), attempt+1)
		passwd, err := passphrase.GetPassphrase(fingerprint, prompt, description, "", true, false)
		if err != nil {
			return nil, err
		}
		key.PrivateKey.Decrypt([]byte(passwd))

		return nil, nil
	})

	msgDetails, err := openpgp.ReadMessage(input, ring, promptFunc, nil)
	if err != nil {
		return nil, err
	}

	return msgDetails.UnverifiedBody, nil
}

func getNameAndFile(prefix string, args []string) (string, string, error) {
	var name string
	var file string
	if len(args) == 0 {
		name = "(default)"
		file = filepath.Join(prefix, ".gpg")
	} else if len(args) == 1 {
		name = args[0]
		file = filepath.Join(prefix, name+".gpg")
	} else {
		return "", "", &CmdError{5, "wrong number of arguments"}
	}
	return name, file, nil
}

func isTerminal(file *os.File) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

func copyToClipboard(data string) error {
	xclip := exec.Command("xclip", "-in", "-loop", "2")
	stdin, err := xclip.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	if err := xclip.Start(); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdin, "%s", data); err != nil {
		return err
	}
	return nil
}
