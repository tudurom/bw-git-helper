package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gobwas/glob"
	"github.com/tudurom/bw-git-helper/pinentry"
	"gopkg.in/ini.v1"
)

const MAX_PASSWORD_TRIALS = 3

type Mapping struct {
	Pattern glob.Glob
	Target string
}

type Config struct {
	Pinentry string
	Mappings []Mapping
}

type ErrBW struct {
	Message string
}

func (e *ErrBW) Error() string {
	return "bitwarden error: " + e.Message
}

var ErrUserCancel = errors.New("user canceled")

type StdinData map[string]string

func (sd StdinData) String() string {
	ret := ""
	for k, v := range(sd) {
		ret += fmt.Sprintf("%s=%s\n", k, v)
	}

	return ret
}

func readConfig(fp string) (Config, error) {
	config := Config{
		Pinentry: "pinentry",
		Mappings: []Mapping{},
	}

	cfg, err := ini.Load(fp)
	if err != nil {
		return Config{}, err
	}

	configSection, err := cfg.GetSection("config")
	if err == nil {
		pinentryKey, err := configSection.GetKey("pinentry")
		if err == nil {
			config.Pinentry = pinentryKey.String()
		}
	}
	for _, section := range(cfg.Sections()) {
		key, err := section.GetKey("target")
		if err == nil {
			patternGlob, err := glob.Compile(section.Name())
			if err != nil {
				return Config{}, fmt.Errorf("couldn't compile pattern '%s': %w", section.Name(), err)
			}
			config.Mappings = append(config.Mappings, Mapping{
				Pattern: patternGlob,
				Target: key.String(),
			})
		}
	}


	return config, nil
}

func readInput() StdinData {
	stdinArgs := StdinData{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		split := strings.SplitN(scanner.Text(), "=", 2)
		if (split[1] != "") {
			stdinArgs[split[0]] = split[1]
		}
	}

	return stdinArgs
}

func readPassword(host string, username string, protocol string, trial int, pinentryProgram string) (string, error) {
	if protocol != "" {
		protocol += "://"
	}
	request := pinentry.Request{ Prompt: "Passphrase:" }
	if username != "" {
		request.Desc = fmt.Sprintf("Requesting credentials for %s%s@%s", protocol, username, host)
	} else {
		request.Desc = fmt.Sprintf("Requesting credentials for %s%s", protocol, host)
	}
	if trial > 1 {
		request.Error = fmt.Sprintf("Wrong passphrase (%d of %d)", trial, MAX_PASSWORD_TRIALS)
	}
	pass, err := request.GetPIN(pinentryProgram)
	if err != nil {
		if errors.Is(err, pinentry.ErrCancel) {
			return "", ErrUserCancel
		}
		return "", err
	}

	return pass, nil
}

func invokeBW(args ...string) (string, error) {
	out, err := exec.Command("bw", args...).CombinedOutput()
	if err != nil {
		var e *exec.ExitError
		if errors.As(err, &e) {
			return "", &ErrBW{Message: string(out)}
		} else {
			return "", err
		}
	}
	return string(out), nil
}

func getSessionToken(password string) (string, error) {
	return invokeBW("unlock", "--raw", password)
}

func getPassword(target string, token string) (string, error) {
	return invokeBW("--session", token, "get", "password", target)
}

func tryGetSessionToken(host string, data StdinData, config *Config) (string, error) {
	for trial := 1; trial <= MAX_PASSWORD_TRIALS; trial++ {
		password, err := readPassword(host, data["username"], data["protocol"], trial, config.Pinentry)
		if err != nil {
			return "", err
		}

		token, err := getSessionToken(password)
		if err != nil {
			var e *ErrBW
			if (errors.As(err, &e)) {
				continue
			} else {
				return "", err
			}
		}
		return token, nil
	}

	return "", errors.New("too many requests")
}


func findFirstTarget(host string, mappings []Mapping) *Mapping {
	for _, mapping := range(mappings) {
		if mapping.Pattern.Match(host) {
			return &mapping
		}
	}

	return nil
}

func ifError(err interface{}) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	cfgFlag := flag.String("c", "", "path to config file")
	flag.Parse()

	if len(flag.Args()) != 1 {
		ifError("Expected an operation")
	}
	if flag.Args()[0] != "get" {
		// fail silently
		os.Exit(0)
	}

	var configFilePath string
	if *cfgFlag != "" {
		configFilePath = *cfgFlag
	} else {
		var err error
		configFilePath, err = xdg.ConfigFile("bw-git-helper/config.ini")
		ifError(err)
	}
	config, err := readConfig(configFilePath)
	ifError(err)

	data := readInput()
	host, ok := data["host"]
	if !ok {
		ifError("host parameter required")
	}
	if path, ok := data["path"]; ok {
		host += "/" + path
	}

	token, err := tryGetSessionToken(host, data, &config)
	ifError(err)

	firstMapping := findFirstTarget(host, config.Mappings)
	if firstMapping == nil {
		ifError("no matching mapping")
	}

	pass, err := getPassword(firstMapping.Target, token)
	ifError(err)

	fmt.Printf("password=%s\n", pass)
}
