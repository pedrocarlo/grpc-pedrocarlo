package repl

import (
	"fmt"
	"grpc-pedrocarlo/pkg/client"
	filesync "grpc-pedrocarlo/pkg/file"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type Command struct {
	f     func(*client.FileClient, []string)
	name  string
	usage string
}

type CommandMap struct {
	commands  *map[string]Command
	completer *readline.PrefixCompleter
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("mode"),
	readline.PcItem("login"),
	readline.PcItem("setprompt"),
	readline.PcItem("setpassword"),
	readline.PcItem("bye"),
	readline.PcItem("help"),
	readline.PcItem("sleep"),
)

// ****************** REPL FUNCTIONS **************************
// See these for an example of how to get a REPL with history
// like in the IP/TCP reference
// This REPL relies on a go module for "readline".  To add it to your project,
// run:  "github.com/chzyer/readline"
func Repl() {
	l := ReplInitialize()
	defer l.Close()
	for {
		line, done := ReplGetLine(l)
		if done {
			break
		}
		if line == "" {
			usage(os.Stdout)
		}

	}
}

// Initialize the repl
func ReplInitialize() *readline.Instance {
	l, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       "./client_files/tmp/history.tmp",
		InterruptPrompt:   "^C",
		HistorySearchFold: true,
		AutoComplete:      completer,
	})
	if err != nil {
		panic(err)
	}
	return l
}

// Get a line from the repl
// To keep the example clean, we abstract this into a helper.
// For better error handling, you may just want to do this in the loop that reads a line
func ReplGetLine(repl *readline.Instance) (string, bool) {
	line, err := repl.Readline()
	if err == readline.ErrInterrupt {
		return "", true
	} else if err == io.EOF {
		return "", true
	}

	line = strings.TrimSpace(line)

	return line, false
}

func initializeCommands() {
	commandMap := &CommandMap{}
	commands := make(map[string]Command)
}

func UploadFile(c *client.FileClient, args []string) {
	if len(args) < 2 {
		fmt.Println("usage: upload <filepath> <remote_folder>")
		return
	}
	filepath, folder := args[0], args[1]
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = c.UploadFile(file, folder)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DownloadFile(c *client.FileClient, args []string) {
	if len(args) < 1 {
		fmt.Println("usage: download <remote_filename> [<remote_dir>]")
		return
	}
	var folder string = c.Curr_dir
	if len(args) == 2 {
		folder = args[1]
	}
	filename := args[0]
	var file_meta *filesync.FileMetadata = nil
	if folder != c.Curr_dir {
		// Get File List for dir
		files, err := c.GetFileList(folder)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, f := range files {
			if f.Filename == filename {
				file_meta = f
				break
			}
		}
		// Error here did not find file
		if file_meta == nil {
			fmt.Errorf("cannot find file %s in folder %s\n", filename, folder)
			return
		}
	} else {
		placeholder_meta, ok := c.Curr_dir_files[filename]
		if !ok {
			fmt.Errorf("cannot find file %s in folder %s\n", filename, folder)
			return
		}
		file_meta = placeholder_meta
	}
	err := c.DownloadFile(file_meta)
	if err != nil {
		fmt.Println(err)
		return
	}
}
