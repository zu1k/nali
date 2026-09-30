#!/bin/sh
# Build Debian packages from the Linux binaries in $BINDIR.
# Usage: build_deb.sh <version>   (run `make deb` instead of calling it directly)
set -eu

NAME=${NAME:-nali}
BINDIR=${BINDIR:-bin}
VERSION=${1:?version required}

# v0.8.1 -> 0.8.1, v0.8.1-14-gabc123 -> 0.8.1+14.gabc123
DEB_VERSION=$(echo "$VERSION" | sed 's/^v//; s/-/+/; s/-/./g')
case "$DEB_VERSION" in
	[0-9]*) ;;
	*) DEB_VERSION="0.0.0+$DEB_VERSION" ;;
esac

# debian architecture : Makefile target
for pair in amd64:linux-amd64 i386:linux-386 armhf:linux-armv7 arm64:linux-armv8; do
	arch=${pair%%:*}
	target=${pair#*:}
	bin="$BINDIR/$NAME-$target"
	[ -f "$bin" ] || { echo "missing $bin, run 'make $target' first" >&2; exit 1; }

	root=$(mktemp -d)
	chmod 0755 "$root"
	install -D -m 0755 "$bin" "$root/usr/bin/$NAME"
	install -D -m 0644 LICENSE "$root/usr/share/doc/$NAME/copyright"
	install -D -m 0644 README_en.md "$root/usr/share/doc/$NAME/README.md"
	mkdir -p "$root/DEBIAN"
	cat > "$root/DEBIAN/control" <<CONTROL
Package: $NAME
Version: $DEB_VERSION
Section: net
Priority: optional
Architecture: $arch
Maintainer: zu1k <i@zu1k.com>
Homepage: https://github.com/zu1k/nali
Description: offline tool for querying IP geographic information and CDN provider
 Nali annotates IP addresses and CDN domain names in any text, e.g. the output
 of dig, nslookup or traceroute, with their location and CDN provider.
CONTROL
	dpkg-deb --root-owner-group --build "$root" "$BINDIR/${NAME}_${DEB_VERSION}_${arch}.deb" >/dev/null
	rm -rf "$root"
	echo "$BINDIR/${NAME}_${DEB_VERSION}_${arch}.deb"
done
