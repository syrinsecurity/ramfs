# RAMFS

RAMFS is a in memory file system which has the capability of creating and removing directories and files. Files can be written to and read from but also read/write permissions can be set in order to protect the filesystem if desired.

## Install

```
go get github.com/syrinsecurity/ramfs
```

## Examples

```go
package main

import (
	"fmt"

	"github.com/syrinsecurity/ramfs"
)

func main() {

	vfs := ramfs.New()

	//Create directories
	vfs.Mkdir("/users")
	vfs.Mkdir("/users/root")
	vfs.Mkdir("/users/admin")

	fmt.Println("Listing of /users:")
	//List directories and files
	_, dirs, _ := vfs.Ls("/users")
	for _, dir := range dirs {
		fmt.Println("DIR:", dir.Name)
	}

	//Remote a directory. Note in order to delete a directory it must end in a trailing "/"
	vfs.Rm("/users/root/")

	fmt.Println("Listing of /users after deleting /users/root:")
	//List directories and files
	_, dirs, _ = vfs.Ls("/users")
	for _, dir := range dirs {
		fmt.Println("DIR:", dir.Name)
	}

	//Write or modify a files contents
	vfs.WriteFile("/readme.txt", []byte("Hi, you are reading the readme file in the root dir, weldone."))

	fmt.Println("Readme.txt contents:")
	data, _ := vfs.FileGetContents("/readme.txt")

	fmt.Println(string(data))
}

```
