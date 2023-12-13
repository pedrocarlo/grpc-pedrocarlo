# grpc-pedrocarlo

## GRPC File Sharing

## Description

This is a GRPC File Sharing Client and Server for the Final Project of CS1680 course I have been auditing at Brown. Special and huge thank you to Nick DeMarinis for allowing me to be present in your lectures and for answering my barrage of questions. I have learned a lot with you!

This project aims to replicate in some ways the behavior of a FTP server by allowing a client to connect to the GRPC server and Upload, Download, Remove Files and Directories and more you would expect from a file system on the "cloud".

The motivation for this project was to understand more the benefits of GRPC and some of the intricacies of creating and managing a file server. Also, it seems a better alternative to send files in an API instead of json, as you can have a continuos stream of messages flowing through it. I have not dabbled yet with websockets, but GRPC can seems to be able to replicate the socket behavior well enough with the added error handling. 

## Install

Clone the project

PREREQUISITES:

- Go Version = 1.21

Do not know if any other version works

HTTPS: 
```shell
git clone https://github.com/pedrocarlo/grpc-pedrocarlo.git
```
SET GO PATH:

```shell
export PATH="$PATH:$(go env GOPATH)/bin"
```

BUILD FROM SOURCE:

```shell
make server
```
```shell
make client
```

## USAGE

- Initialize GRPC Server
```shell
./server
```

This will initialize the server by default at 127.0.0.1:7070. This can be modified in the cmd/main.go file in line 20

- Initialize Client
```shell
./client
```

## Commands
In all commands you can always use relative paths or absolute paths

- ### Cd 
    - ```cd <remote_folder>```
    - Changes the directory the client is in

- ### Mkdir 
    - ```mkdir <remote_folder>```
    - Creates a folder in the server

- ### Rmdir 
    - ```rmdir <remote_folder>```
    - Removes empty directory from server

- ### Rm 
    - ```rm <remote_filename> <remote_folder>```
    - Creates a folder in the server
    
- ### Ls 
    - ```ls <remote_folder>```
    - List files from a directory

- ### Download 
    - ```download <remote_filename> [<remote_folder>]```
    - Download a file from the remote directory. Remote folder can be omitted to select the current client directory. Files are download to the ./client_files on the directory the binary is located on. Did not create logic for the folder to be created automatically, so if an errors occurs create a client_files folder with a downloads folder and a tmp folder. 

- ### Upload 
    - ```upload <filepath> <remote_folder>```
    - Upload a file from your local machine to remote folder. Your server should have the following folder structure in the location the binary is created -> server_files with files folder and a tmp folder. Files uploaded to the server are stored in ./server_files/files/ .

## Improvements

- Implement caching of file listing on the client so you do not call the server everytime to know the files in that directory
- Auto-create the server_files and client_files folder
- Create automated tests for relative pathing and for for each RPC service. 

## Conclusion

- Overall, I had a great experience with GRPC will definitely have it in the back of my mind when creating a new project that uses an API. Thanks again Nick DeMarinis for allowing me to be in lectures!