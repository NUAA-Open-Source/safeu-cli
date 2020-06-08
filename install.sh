#!/bin/bash
#
# This script should be run via curl:
#   sh -c "$(curl -fsSL https://raw.githubusercontent.com/NUAA-Open-Source/safeu-cli/master/install.sh)"
# or wget:
#   sh -c "$(wget -qO- https://raw.githubusercontent.com/NUAA-Open-Source/safeu-cli/master/install.sh)"
#
# As an alternative, you can first download the install script and run it afterwards:
#   wget https://raw.githubusercontent.com/NUAA-Open-Source/safeu-cli/master/install.sh
#   sh install.sh

SAFEU_RELEASE="https://github.com/NUAA-Open-Source/safeu-cli/releases/download/v1.0.0-alpha/safeu-linux-x64"
SAFEU_CN_RELEASE="https://triplez-public-1251926021.cos.ap-shanghai.myqcloud.com/safeu-cli/v1.0.0-alpha/safeu-linux-x64"
BIN_DIR=/usr/local/bin
BIN_FILENAME=safeu.tmp
SAFEU_CMD=safeu
IS_LOCAL=0


error() {
	echo ${RED}"Error: $@"${RESET} >&2
}

setup_color() {
	# Only use colors if connected to a terminal
	if [ -t 1 ]; then
		RED=$(printf '\033[31m')
		GREEN=$(printf '\033[32m')
		YELLOW=$(printf '\033[33m')
		BLUE=$(printf '\033[34m')
		BOLD=$(printf '\033[1m')
		RESET=$(printf '\033[m')
	else
		RED=""
		GREEN=""
		YELLOW=""
		BLUE=""
		BOLD=""
		RESET=""
	fi
}

download_safeu_cli() {
    if [ "$1" = "cn" ]; then
        wget -cO ${BIN_FILENAME} ${SAFEU_CN_RELEASE} || {
            error "cannot download safeu-cli by ${SAFEU_CN_RELEASE}"
            exit 1
        }
    else
        wget -cO ${BIN_FILENAME} ${SAFEU_RELEASE} || {
            error "cannot download safeu-cli by ${SAFEU_RELEASE}"
            exit 1
        }
    fi

}

install_scope() {
    printf "${YELLOW}Install safeu command tool globally (require sudo permission later) ? [Y/N, Defualt: Y]: "
    read isGlobal

    if [ "$isGlobal" = "n" ]  || [ "$isGlobal" = "N" ] ; then
        BIN_DIR=${HOME}/.local/bin
        IS_LOCAL=1
    else
        BIN_DIR=/usr/local/bin
    fi
}

install_safeu_cli() {
    if [ ${IS_LOCAL} -eq 1 ]; then
        install -Dm755 "${BIN_FILENAME}" "${BIN_DIR}/${SAFEU_CMD}" || {
            error "install the safeu-cli tool failed"
            exit 1
        }
    else
        sudo install -Dm755 "${BIN_FILENAME}" "${BIN_DIR}/${SAFEU_CMD}" || {
            error "install the safeu-cli tool failed"
            exit 1
        }
    fi
}

post_install() {
    rm -f ${BIN_FILENAME}
    printf ${GREEN}

    cat <<-'EOF'
         ____         __      _   _    ____ _     ___ 
        / ___|  __ _ / _| ___| | | |  / ___| |   |_ _|
        \___ \ / _` | |_ / _ \ | | | | |   | |    | | 
         ___) | (_| |  _|  __/ |_| | | |___| |___ | | 
        |____/ \__,_|_|  \___|\___/   \____|_____|___|  ....is now installed!

        Now you can upload and download files via "safeu" command !
        If you have further questions, you can find support in here: 
                https://github.com/NUAA-Open-Source/safeu-cli/issues
        
EOF
    printf "        Current installed safeu-cli version: $(safeu version)\n"
    printf ${RESET}
}

main() {
    setup_color
    # download_safeu_cli $1
    install_scope
    install_safeu_cli
    post_install
}

main "$@"
