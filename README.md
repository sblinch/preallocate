# preallocate
[![GoDoc](https://godoc.org/gitlab.com/tslocum/preallocate?status.svg)](https://godoc.org/gitlab.com/tslocum/preallocate)
[![CI status](https://gitlab.com/tslocum/preallocate/badges/master/pipeline.svg)](https://gitlab.com/tslocum/preallocate/commits/master)
[![Donate](https://img.shields.io/liberapay/receives/rocketnine.space.svg?logo=liberapay)](https://liberapay.com/rocketnine.space)

File preallocation library

## Features

- Allocates files efficiently (via syscall) on the following platforms:
  - [Linux](http://man7.org/linux/man-pages/man2/fallocate.2.html)
  - [Windows](https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/ntifs/nf-ntifs-ntsetinformationfile)
- Falls back to writing null bytes

## Documentation

Docs are hosted on [godoc.org](https://godoc.org/gitlab.com/tslocum/preallocate).

## Support

Please share issues/suggestions [here](https://gitlab.com/~tslocum/preallocate/issues).
