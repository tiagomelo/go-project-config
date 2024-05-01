// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package cfg

import (
	"io/fs"
	"os"
)

// File defines the interface for file operations.
type File interface {
	WriteString(s string) (n int, err error)
	Name() string
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
}

// osFile struct implements the File interface using the os.File type.
type osFile struct {
	*os.File
}

func (f osFile) WriteString(s string) (n int, err error) {
	return f.File.WriteString(s)
}

func (f osFile) Name() string {
	return f.File.Name()
}

func (f osFile) Read(p []byte) (n int, err error) {
	return f.File.Read(p)
}

func (f osFile) Write(p []byte) (n int, err error) {
	return f.File.Write(p)
}

func (f osFile) Close() error {
	return f.File.Close()
}

// fileSystem interface abstracts the file system operations. This allows
// for easier testing by mocking file system interactions. It includes
// methods for opening, reading, writing files, and checking their status.
type fileSystem interface {
	Open(name string) (File, error)
	Create(name string) (File, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	IsNotExist(err error) bool
	Mkdir(dirName string) error
}

// osFileSystem struct implements the fileSystem interface using
// the standard library's os package. This is the real implementation
// that interacts with the actual file system.
type osFileSystem struct{}

func (osFileSystem) Open(name string) (File, error) {
	file, err := os.Open(name)
	return osFile{file}, err
}

func (osFileSystem) Create(name string) (File, error) {
	return os.Create(name)
}

func (osFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (osFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (osFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (osFileSystem) IsDir(fi fs.FileInfo) bool {
	return fi.IsDir()
}

func (osFileSystem) Mkdir(dirName string) error {
	return os.Mkdir(dirName, os.ModePerm)
}
