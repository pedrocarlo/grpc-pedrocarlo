package repl

import (
	"fmt"
	"grpc-pedrocarlo/pkg/client"
	filesync "grpc-pedrocarlo/pkg/file"
	"grpc-pedrocarlo/pkg/utils"
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

func usage(commandMap *CommandMap, w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, commandMap.completer.Tree("    "))
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

func (commandMap *CommandMap) buildCompleter(c *client.FileClient) {
	completer := readline.PrefixCompleter{}
	children := make([]readline.PrefixCompleterInterface, 0)
	listCwd := func(string) []string {
		return listFiles(c, []string{"."})
	}
	for name := range commandMap.commands {
		children = append(children, readline.PcItem(name, readline.PcItemDynamic(listCwd)))
	}
	completer.SetChildren(children)
	commandMap.completer = &completer
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
	commands["mkdir"] = Command{
		f:    Mkdir,
		name: "mkdir",
		desc: "Make remote directory",
	}
	commands["rm"] = Command{
		f:    RemoveFile,
		name: "rm",
		desc: "Remove file from server",
	}
	commands["rmdir"] = Command{
		f:    RemoveDir,
		name: "rmdir",
		desc: "Remove empty dir from server",
	}
	commands["cd"] = Command{
		f:    ChangeDir,
		name: "cd",
		desc: "Change directory",
	}
	return commandMap
}

// ****************** REPL FUNCTIONS **************************
// See these for an example of how to get a REPL with history
// like in the IP/TCP reference
// This REPL relies on a go module for "readline".  To add it to your project,
// run:  "github.com/chzyer/readline"
func Repl(c *client.FileClient) {
	commandMap := initializeCommands()
	commandMap.buildCompleter(c)
	ChangeDir(c, []string{"/"})

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
			usage(commandMap, os.Stdout)
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
		AutoComplete:      commandMap.completer,
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
	// TODO CHANGE SPLITLIST
	split_path := strings.Split(folder, string(os.PathSeparator))
	if len(split_path) > 0 {
		if split_path[0] == "." {
			split_path[0] = c.Curr_dir
		} else if split_path[0] == ".." {
			split_path[0] = filepath.Dir(c.Curr_dir)
		} else if split_path[0] == "" {
			split_path[0] = "/"
		}
		// else {
		// 	split_path[0] = c.Curr_dir + split_path[0]
		// }
	}
	folder = filepath.Join(split_path...)
	if folder == "." {
		folder = c.Curr_dir
	} else if folder == ".." {
		if filepath.Dir(folder) == "." {
			folder = c.Curr_dir
		} else {
			folder = filepath.Dir(folder)
		}
	} else if !strings.HasPrefix(folder, "/") {
		folder = c.Curr_dir + folder
	}
	return folder
}

func listFiles(c *client.FileClient, args []string) []string {
	out := make([]string, 0)
	if len(args) < 1 {
		fmt.Println("usage: ls <remote_folder>")
		return out
	}
	folder := translateFolderClient(c, args[0])
	files, err := c.GetFileList(folder)
	if err != nil {
		fmt.Println(err)
		return out
	}
	for _, file := range files {
		if file.IsDir {
			out = append(out, file.Filename+"/")
		} else {
			out = append(out, file.Filename)
		}
	}
	return out
}

func ListFiles(c *client.FileClient, args []string) {
	out := listFiles(c, args)
	out_str := ""
	for _, name := range out {
		out_str += fmt.Sprintf("%10s", name)
	}
	fmt.Printf("%s\n", out_str)
}

func Mkdir(c *client.FileClient, args []string) {
	if len(args) < 1 {
		fmt.Println("usage: mdkir <remote_folder>")
		return
	}
	folder := translateFolderClient(c, args[0])
	utils.Log_trace(folder)
	_, err := c.Mkdir(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RemoveFile(c *client.FileClient, args []string) {
	if len(args) < 2 {
		fmt.Println("usage: rm <remote_filename> <remote_folder>")
		return
	}
	filename, folder := args[0], translateFolderClient(c, args[1])
	err := c.RemoveFile(folder, filename)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func RemoveDir(c *client.FileClient, args []string) {
	if len(args) < 1 {
		fmt.Println("usage: rmdir <remote_folder>")
		return
	}
	folder := translateFolderClient(c, args[0])
	err := c.RemoveDir(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Tab autocomplete only completes for current folder at the moment
func ChangeDir(c *client.FileClient, args []string) {
	if len(args) < 1 {
		fmt.Println("usage: cd <remote_folder>")
		return
	}
	folder := translateFolderClient(c, args[0])
	// Query its parent folder, to see if it either errors or if it is inside
	out, err := c.GetFileList(filepath.Dir(folder))
	if err != nil {
		fmt.Println(err)
		return
	}
	if folder == "/" {
		c.Curr_dir = folder
		return
	}
	for _, file := range out {
		if file.IsDir && filepath.Join(file.Folder, file.Filename) == folder {
			// Attempt at Update cache
			files, err := c.GetFileList(folder)
			if err != nil {
				fmt.Println(err)
				return
			}
			c.Curr_dir = folder
			c.Curr_dir_files = make(map[string]*filesync.FileMetadata)
			for _, file := range files {
				c.Curr_dir_files[file.Filename] = file
			}
			// TODO Restart Cache Timer
			return
		}
	}
	fmt.Println(fmt.Errorf("folder not found"))
}
