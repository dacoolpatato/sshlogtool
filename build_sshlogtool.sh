#!/bin/bash
set -e

APPNAME="sshlogtool"
VERSION="1.0"

ARCHS=("amd64" "arm64")
DEB_BASE="./debbuild"

# Clean old builds
rm -rf $DEB_BASE
mkdir -p $DEB_BASE

# Check Go installation
if ! command -v go &> /dev/null
then
    echo "Go is required but not installed. Please install Go first."
    exit 1
fi

for ARCH in "${ARCHS[@]}"; do
    echo "Building for arch: $ARCH"

    # Build binary statically
    GOOS=linux GOARCH=$ARCH CGO_ENABLED=0 go build -ldflags="-s -w" -o $APPNAME ./sshlogtool.go

    # Prepare deb folder structure
    DEB_DIR="$DEB_BASE/${APPNAME}_${VERSION}_${ARCH}"
    mkdir -p $DEB_DIR/DEBIAN
    mkdir -p $DEB_DIR/usr/bin
    mkdir -p $DEB_DIR/etc/systemd/system
    mkdir -p $DEB_DIR/var/log

    # Copy binary
    cp $APPNAME $DEB_DIR/usr/bin/

    # Create empty log file
    touch $DEB_DIR/var/log/${APPNAME}.log
    chmod 644 $DEB_DIR/var/log/${APPNAME}.log

    # Create control file
    cat > $DEB_DIR/DEBIAN/control <<EOF
Package: $APPNAME
Version: $VERSION
Section: utils
Priority: optional
Architecture: $ARCH
Maintainer: Dario <dario@example.com>
Description: SSH login monitor and logger tool for Linux servers and desktops
EOF

    # Create systemd service file
    cat > $DEB_DIR/etc/systemd/system/${APPNAME}.service <<EOF
[Unit]
Description=SSH Login Monitor Service - sshlogtool
After=network.target

[Service]
ExecStart=/usr/bin/$APPNAME -watch
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target
EOF

    # Build the deb package
    dpkg-deb --build $DEB_DIR

    echo "Created package: ${DEB_DIR}.deb"
done

# Cleanup binary
rm $APPNAME
