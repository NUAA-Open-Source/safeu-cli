<!-- omit in toc -->
# SafeU CLI

![Install](https://github.com/NUAA-Open-Source/safeu-cli/workflows/Install/badge.svg) | ![Go](https://github.com/NUAA-Open-Source/safeu-cli/workflows/Go/badge.svg)

A command line tool for SafeU (https://safeu.a2os.club).

- [Install](#install)
  - [China mainland optimized](#china-mainland-optimized)
- [Usage](#usage)
  - [Upload](#upload)
    - [Upload one file](#upload-one-file)
    - [Upload more than one file](#upload-more-than-one-file)
    - [Set Password / Download Count / Expired Time](#set-password--download-count--expired-time)
    - [Full detail of upload command](#full-detail-of-upload-command)
  - [Download](#download)
    - [Full detail of get command](#full-detail-of-get-command)
- [Demo](#demo)
- [Compile](#compile)
- [Known issues](#known-issues)
- [License](#license)

## Install

> If you are in China mainland, the install method in [China mainland optimized](#china-mainland-optimized) is a better choice. 

> NOTICE: The following methods would download a pre-compiled safeu-cli binary file which is ONLY for 64-bit Linux. If you are using a different architecture or OS, just check the [Compile](#compile) section below to build your own binary package.

`safeu-cli` is installed by running one of the following commands in your terminal. You can install this via the command-line with either `curl` or `wget`.

via curl:

```bash
$ sh -c "$(curl -fsSL https://raw.githubusercontent.com/NUAA-Open-Source/safeu-cli/master/install.sh)"
```

via wget:

```bash
$ sh -c "$(wget -qO- https://raw.githubusercontent.com/NUAA-Open-Source/safeu-cli/master/install.sh)"
```

Congratulations, you have successfully installed the `safeu-cli` tool :tada:

### China mainland optimized

via curl:

```bash
$ sh -c "$(curl -fsSL https://gitee.com/A2OS/safeu-cli/raw/master/install.sh) cn"
```

via wget:

```bash
$ sh -c "$(wget -qO- https://gitee.com/A2OS/safeu-cli/raw/master/install.sh) cn"
```

## Usage

### Upload

#### Upload one file
```bash
$ safeu upload filename
```

#### Upload more than one file

```bash
$ safeu upload filename1 filename2 filename3
```

#### Set Password / Download Count / Expired Time

Ref to [Full deteail of upload command](#full-detail-of-upload-command).

Examples for this section will be supplement lately.

#### Full detail of upload command

```bash
$ safeu upload --help          
Send and Share file by this command.
SafeU is responsible for ensuring upload speed and file safety

Usage:
  safeu upload [flags]

Flags:
  -d, --downcount int     specific down count
  -e, --expiretime int    specific expire time
  -p, --password string   specific password
  -r, --recode string     specific recode
  -h, --help              help for upload
```

### Download
```bash
$ safeu get your_recode
```

#### Full detail of get command

```bash
$ safeu get --help   
Download file(s) by this command.
SafeU is responsible for ensuring download speed and file safety :)

Usage:
  safeu get [flags]

Flags:
  -d, --dir string        download to specific directory
  -p, --password string   specific password
      --print             print the file URL directly, then you can 
                          download the file by other download tools 
                          (e.g. wget, aria2).
  -h, --help              help for get
```

## Demo

[![asciicast](https://asciinema.org/a/iZgrbrUpli4kxOQlOBco9jamH.svg)](https://asciinema.org/a/iZgrbrUpli4kxOQlOBco9jamH)

## Compile

```bash
# build binary for your OS 
# need go version > 1.13
make build

# build binary for Linux
make linux-build
```

## Known issues

## License

This project is open-sourced by [Apache 2.0](./LICENSE).
