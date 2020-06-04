# Nuwa
Command line for Safeu (https://safeu.a2os.club)


## Usage

### Upload

#### upload one file
```bash
safeu upload filename
```

#### upload more than one file

```bash
safeu upload filename1 filename2 filename3
```

#### Full detail of upload command

```bash
$ ./safeu upload --help          
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

#### Set Password/Download Count/Expired Time
just in next release 😃
#### Download
just in next release 😃
## Demo
[![asciicast](https://asciinema.org/a/4G28AuGu92QSRG5NqjD4yhbSI.svg)](https://asciinema.org/a/4G28AuGu92QSRG5NqjD4yhbSI)
## Compile

```bash
# build binary for your OS 
# need go version > 1.13
make build

# build binary for Linux
make linux-build
```

