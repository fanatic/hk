package main

import (
	"bufio"
	"code.google.com/p/go-netrc/netrc"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	Version = "0.0.2"
)

var (
	apiURL = "https://api.heroku.com"
	hkHome = os.Getenv("HOME") + "/.hk"
)

var stdin = bufio.NewReader(os.Stdin)

var updater = Updater{
	url: "https://github.com/downloads/kr/hk/",
	dir: hkHome + "/update/",
}

type Command struct {
	// args does not include the command name
	Run func(cmd *Command, args []string)

	Usage string // first word is the command name
	Short string // `hk help` output
	Long  string // `hk help <cmd>` output
}

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Running `hk help` will list commands in this order.
var commands = []*Command{
	cmdCreate,
	cmdDestroy,
	cmdCreds,
	cmdEnv,
	cmdFetchUpdate,
	cmdGet,
	cmdSet,
	cmdInfo,
	cmdList,
	cmdOpen,
	cmdPs,
	cmdVersion,
	cmdHelp,
}

var flagApp = flag.String("a", "", "app")

func main() {
	defer updater.run() // doesn't run if os.Exit is called

	if s := os.Getenv("HEROKU_API_URL"); s != "" {
		apiURL = strings.TrimRight(s, "/")
	}

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	name := args[0]
	os.Args = args
	flag.Parse()
	args = flag.Args()

	for _, cmd := range commands {
		if cmd.Name() == name {
			cmd.Run(cmd, args)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", name)
	usage()
}

func getCreds(u *url.URL) (user, pass string) {
	if u.User != nil {
		pw, _ := u.User.Password()
		return u.User.Username(), pw
	}

	m, err := netrc.FindMachine(os.Getenv("HOME")+"/.netrc", u.Host)
	if err != nil {
		log.Fatal(err)
	}

	return m.Login, m.Password
}

func app() string {
	if *flagApp != "" {
		return *flagApp
	}
	out, err := exec.Command("git", "remote", "show", "-n", "heroku").Output()
	if err != nil {
		log.Fatal(err)
	}
	s := string(out)
	const sign = "Fetch URL: "
	i := strings.Index(s, sign)
	if i < 0 {
		log.Fatal("could not find git remote named 'heroku'")
	}
	s = s[i+len(sign):]
	i = strings.Index(s, "\n")
	if i >= 0 {
		s = s[:i]
	}
	if !strings.HasPrefix(s, gitURLPre) || !strings.HasSuffix(s, gitURLSuf) {
		log.Fatal("could not find app name in heroku git remote")
	}
	return s[len(gitURLPre) : len(s)-len(gitURLSuf)]
}