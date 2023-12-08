package repl

import (
	"fmt"
	"grpc-pedrocarlo/pkg/client"
	filesync "grpc-pedrocarlo/pkg/file"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

type Command struct {
	f    func(*client.FileClient, []string)
	name string
	desc string
}

type CommandMap struct {
	commands  map[string]Command
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
func Repl(c *client.FileClient) {
	commandMap := initializeCommands()
	l := ReplInitialize(commandMap)
	defer l.Close()
	for {
		line, done := ReplGetLine(l)
		if done {
			break
		}
		args := strings.Split(line, " ")
		name := args[0]
		if line == "" {
			usage(os.Stdout)
		} else {
			command, ok := commandMap.commands[name]
			if ok {
				command.f(c, args[1:])
			} else {
				fmt.Printf("Command '%s' not found\n", args[0])
			}
		}

	}
}

// Initialize the repl
func ReplInitialize(commandMap *CommandMap) *readline.Instance {
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

func initializeCommands() *CommandMap {
	commandMap := &CommandMap{}
	commands := make(map[string]Command)
	commandMap.commands = commands
	commands["upload"] = Command{
		f:    UploadFile,
		name: "upload",
		desc: "Uploads a file from the hosts machine to the server",
	}
	commands["download"] = Command{
		f:    DownloadFile,
		name: "download",
		desc: "Downloads a file from a folder on the server to the client_files/downloads/ folder",
	}
	commands["ls"] = Command{
		f:    ListFiles,
		name: "ls",
		desc: "List Files from remote folder",
	}
	return commandMap
}

func UploadFile(c *client.FileClient, args []string) {
	if len(args) < 2 {
		fmt.Println("usage: upload <filepath> <remote_folder>")
		return
	}
	filepath, folder := args[0], translateFolderClient(c, args[1])
	fmt.Println("folder:", folder)
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
	folder = translateFolderClient(c, folder)
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
			fmt.Println(fmt.Errorf("cannot find file %s in folder %s", filename, folder))
			return
		}
	} else {
		placeholder_meta, ok := c.Curr_dir_files[filename]
		if !ok {
			fmt.Println(fmt.Errorf("cannot find file %s in folder %s", filename, folder))
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

func translateFolderClient(c *client.FileClient, folder string) string {
	split_path := filepath.SplitList(filepath.Join(folder, ""))
	if len(split_path) > 0 {
		if split_path[0] == "." {
			split_path[0] = c.Curr_dir
		} else if split_path[0] == ".." {
			split_path[0] = filepath.Dir(c.Curr_dir)
		} else {
			split_path[0] = c.Curr_dir + split_path[0]
		}
	}
	folder = filepath.Join(split_path...)
	if folder == "." {
		folder = c.Curr_dir
	}
	//
	if folder == ".." {
		folder = c.Curr_dir
	}
	return folder
}

func ListFiles(c *client.FileClient, args []string) {
	if len(args) < 1 {
		fmt.Println("usage: ls <remote_folder>")
		return
	}
	folder := translateFolderClient(c, args[0])
	files, err := c.GetFileList(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	out_str := ""
	for _, file := range files {
		if file.Filename == "" {
			base := filepath.Base(file.Folder)
			out_str += fmt.Sprintf("%10s/", base)
		} else {
			out_str += fmt.Sprintf("%10s", file.Filename)
		}
	}
	fmt.Printf("%s\n", out_str)
}
