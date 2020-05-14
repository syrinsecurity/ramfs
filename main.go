package main

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
)

var (
	//ErrorsDataSourceNotSupported is returned when a file has a data source what is not supported
	ErrorsDataSourceNotSupported = errors.New("ErrorsDataSourceNotSupported")

	//ErrorsDirectoryNotFound is returned when a directory can not be resolved, this may be due to it not existing
	ErrorsDirectoryNotFound = errors.New("directory not found")

	//ErrorsFileNotFound is returned when a file can not be resolved, this may be due to it not existing
	ErrorsFileNotFound = errors.New("file not found")

	//ErrorsFileSystemReadOnly is returned when you try modify a readonly file system
	ErrorsFileSystemReadOnly = errors.New("file system is readonly")

	//ErrorsCanNotReCreateRoot is returned when you try make the dir "/"
	ErrorsCanNotReCreateRoot = errors.New("you can not recreate the root folder")

	//ErrorsNoParentDirectory will be returned if you try make a directory without a parent
	ErrorsNoParentDirectory = errors.New("you can not create a directory with out a parent")

	//ErrorsNoWritePermission means either the direct or parent node is read only
	ErrorsNoWritePermission = errors.New("no write permission to modify that file/directory")

	//ErrorsNoReadPermission means either the direct or parent node has denied read access
	ErrorsNoReadPermission = errors.New("no read permission to view contents that file/directory")
)

const (
	OptionReadOnly    = 0
	OptionDisalowRead = 1
)

type Option int

//NewRamFS create a new memory based file system
func NewRamFS(options ...Option) RamFileSystem {
	fs := RamFileSystem{
		directories: make(map[string]*Directory),
	}

	var (
		read, write bool
	)

	write = true
	read = true

	for _, option := range options {
		switch option {
		case OptionReadOnly:
			write = false
			read = true
			break
		case OptionDisalowRead:
			read = false
			break
		}
	}

	fs.directories["/"] = &Directory{
		Name: "/",

		Read:  read,
		Write: write,

		Created:  time.Now().UnixNano(),
		Modified: time.Now().UnixNano(),
	}

	return fs
}

//Get 	/root/users/admin
//			/logs
//			/downloads

type RamFileSystem struct {
	directories map[string]*Directory
}

type Directory struct {
	Name string

	Directories []*Directory
	Files       []*File

	Read  bool
	Write bool

	Created  int64
	Modified int64
}

type File struct {
	Name string

	Content []byte

	Read  bool
	Write bool

	Created  int64
	Modified int64
}

func (fs *RamFileSystem) WriteFile(path string, content []byte) error {

	nodePath := strings.Split(cleanPath(path)[1:], "/")
	parent, ok := fs.directories[getParent(nodePath)]
	if ok != true {
		return ErrorsNoParentDirectory
	}

	if parent.Write == false {
		return ErrorsNoWritePermission
	}

	file := &File{
		Name: nodePath[len(nodePath)-1],

		Content: content,

		Read:  true,
		Write: true,

		Created:  time.Now().UnixNano(),
		Modified: time.Now().UnixNano(),
	}

	var rebuiltParentFiles []*File
	for _, parentFile := range parent.Files {
		if file.Name != parentFile.Name {
			rebuiltParentFiles = append(rebuiltParentFiles, parentFile)
		} else {
			if parentFile.Write == false {
				return ErrorsNoWritePermission
			}
		}
	}

	rebuiltParentFiles = append(rebuiltParentFiles, file)

	parent.Files = rebuiltParentFiles
	parent.Modified = time.Now().UnixNano()

	return nil

}

func (fs *RamFileSystem) FileGetContents(path string) ([]byte, error) {

	nodePath := strings.Split(cleanPath(path)[1:], "/")
	parent, ok := fs.directories[getParent(nodePath)]
	if ok != true {
		return nil, ErrorsNoParentDirectory
	}

	if parent.Read == false {
		return nil, ErrorsNoReadPermission
	}

	for _, file := range parent.Files {
		if file.Name == nodePath[len(nodePath)-1] {

			if file.Read == false {
				return nil, ErrorsNoReadPermission
			}

			return file.Content, nil
		}
	}

	return nil, ErrorsFileNotFound

}

func (fs *RamFileSystem) Mkdir(path string) error {

	if path == "/" {
		return ErrorsNoParentDirectory
	} else if len(path) == 0 {
		return ErrorsNoParentDirectory
	}

	nodePath := strings.Split(cleanPath(path)[1:], "/")

	parent, ok := fs.directories[getParent(nodePath)]
	if ok != true {
		return ErrorsNoParentDirectory
	}

	if parent.Write == false {
		return ErrorsNoWritePermission
	}

	dir := &Directory{
		Name: nodePath[len(nodePath)-1],

		Read:  true,
		Write: true,

		Created:  time.Now().UnixNano(),
		Modified: time.Now().UnixNano(),
	}

	var rebuiltParentDirectories []*Directory
	for _, parentSubDir := range parent.Directories {
		if dir.Name != parentSubDir.Name {
			rebuiltParentDirectories = append(rebuiltParentDirectories, parentSubDir)
		} else {
			if parentSubDir.Write == false {
				return ErrorsNoWritePermission
			}
		}
	}

	rebuiltParentDirectories = append(rebuiltParentDirectories, dir)

	parent.Directories = rebuiltParentDirectories
	parent.Modified = time.Now().UnixNano()

	fs.directories["/"+strings.Join(nodePath, "/")] = dir

	return nil
}

func (fs *RamFileSystem) Ls(path string) ([]*File, []*Directory, error) {
	dir, err := fs.getDir(path)
	if err != nil {
		return nil, nil, err
	}

	return dir.Files, dir.Directories, nil
}

func (fs *RamFileSystem) getDir(path string) (*Directory, error) {

	dir, ok := fs.directories[path]
	if ok != true {
		return nil, ErrorsDirectoryNotFound
	}

	return dir, nil
}

func getParent(nodePath []string) string {

	if len(nodePath) == 0 {
		return "/"
	}

	if len(nodePath) == 1 {
		return "/"
	}

	return "/" + strings.Join(nodePath[:len(nodePath)-1], "/")
}

func cleanPath(path string) string {
	return strings.Replace(filepath.Clean(path), "\\", "/", -1)
}
