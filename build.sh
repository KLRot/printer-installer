#!/usr/bin/env bash
set -euo pipefail
APP_NAME=printer-installer
OUT=dist
APPDIR=appimage/AppDir

rm -rf ${OUT}
mkdir -p ${OUT}
rm -rf ${APPDIR}/usr/bin
mkdir -p ${APPDIR}/usr/bin
mkdir -p ${APPDIR}/usr/share/icons

echo ">>> Building for linux/amd64"
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ${OUT}/${APP_NAME}-x86 main.go

echo ">>> Building for linux/arm64"
CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig go build -o ${OUT}/${APP_NAME}-arm64 main.go

cp ${OUT}/${APP_NAME}-x86 ${APPDIR}/usr/bin/${APP_NAME}-x86
cp ${OUT}/${APP_NAME}-arm64 ${APPDIR}/usr/bin/${APP_NAME}-arm64
chmod +x ${APPDIR}/usr/bin/*

cp appimage/printer-installer.desktop ${APPDIR}/printer-installer.desktop
cp assets/printer.png ${APPDIR}/usr/share/icons/printer.png
cp appimage/AppDir/AppRun ${APPDIR}/AppRun
chmod +x ${APPDIR}/AppRun

if ! command -v appimagetool >/dev/null 2>&1; then
  echo "appimagetool not found."
  exit 1
fi

APPIMAGE_OUT=${OUT}/PrinterInstaller.AppImage
appimagetool ${APPDIR} ${APPIMAGE_OUT}
echo "Done: ${APPIMAGE_OUT}"
