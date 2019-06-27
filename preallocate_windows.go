// +build windows
//
// Support for Windows was cloned from rclone (MIT)
// https://github.com/ncw/rclone/blob/master/backend/local/preallocate_windows.go
//
// Copyright (C) 2012 by Nick Craig-Wood http://www.craig-wood.com/nick/
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package preallocate

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type fileAllocationInformation struct {
	AllocationSize uint64
}

type fileFsSizeInformation struct {
	TotalAllocationUnits     uint64
	AvailableAllocationUnits uint64
	SectorsPerAllocationUnit uint32
	BytesPerSector           uint32
}

type ioStatusBlock struct {
	Status, Information uintptr
}

var (
	ntdll                        = windows.NewLazySystemDLL("ntdll.dll")
	ntQueryVolumeInformationFile = ntdll.NewProc("NtQueryVolumeInformationFile")
	ntSetInformationFile         = ntdll.NewProc("NtSetInformationFile")
)

func preallocFile(file *os.File, size int64) error {
	var (
		iosb       ioStatusBlock
		fsSizeInfo fileFsSizeInformation
		allocInfo  fileAllocationInformation
	)

	// Query info about the block sizes on the file system
	_, _, e1 := ntQueryVolumeInformationFile.Call(
		uintptr(file.Fd()),
		uintptr(unsafe.Pointer(&iosb)),
		uintptr(unsafe.Pointer(&fsSizeInfo)),
		uintptr(unsafe.Sizeof(fsSizeInfo)),
		uintptr(3), // FileFsSizeInformation
	)
	if e1 != nil && e1 != syscall.Errno(0) {
		return WriteSeeker(file, size)
	}

	// Calculate the allocation size
	clusterSize := uint64(fsSizeInfo.BytesPerSector) * uint64(fsSizeInfo.SectorsPerAllocationUnit)
	if clusterSize <= 0 {
		return WriteSeeker(file, size)
	}
	allocInfo.AllocationSize = (1 + uint64(size-1)/clusterSize) * clusterSize

	// Ask for the allocation
	_, _, e1 = ntSetInformationFile.Call(
		uintptr(file.Fd()),
		uintptr(unsafe.Pointer(&iosb)),
		uintptr(unsafe.Pointer(&allocInfo)),
		uintptr(unsafe.Sizeof(allocInfo)),
		uintptr(19), // FileAllocationInformation
	)
	if e1 != nil && e1 != syscall.Errno(0) {
		return WriteSeeker(file, size)
	}

	return nil
}
