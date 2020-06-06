<!-- omit in toc -->
# Nuwa
Command line for Safeu (https://safeu.a2os.club)

- [Install](#install)
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
- [License](#license)

## Install

```bash
$ wget -cO safeu https://github.com/arcosx/Nuwa/releases/download/v0.1-beta/safeu
$ chmod a+x safeu
$ sudo cp safeu /usr/local/bin/safeu
```

Congratulations, you successfully install the `safeu-cli` tool.

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

[![asciicast](https://asciinema.org/a/RkBqLlYnMA3bzXddJsbBRw5Ae.svg)](https://asciinema.org/a/RkBqLlYnMA3bzXddJsbBRw5Ae)

## Compile

```bash
# build binary for your OS 
# need go version > 1.13
make build

# build binary for Linux
make linux-build
```

## License

This project is open-sourced by [Apache 2.0](./LICENSE).