#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu

LIBGIT2_VERSION=1.5.0
LIBGIT2_SHA256=8de872a0f201b33d9522b817c92e14edb4efad18dae95cf156cf240b2efff93e
LIBGIT2_URL=https://github.com/libgit2/libgit2/archive/refs/tags/v${LIBGIT2_VERSION}.tar.gz

current_version=$(pkg-config --modversion libgit2 2>/dev/null || true)
case "$current_version" in
    1.5.*)
        echo "libgit2 $current_version is already installed"
        exit 0
        ;;
esac

if [ "$(id -u)" -ne 0 ]; then
    echo "install-libgit2.sh must run as root" >&2
    exit 1
fi

build_dir=$(mktemp -d)
trap 'rm -rf "$build_dir"' EXIT HUP INT TERM

export DEBIAN_FRONTEND=noninteractive

# Debian releases are pulled from security.debian.org/deb.debian.org only
# while they have active (LTS) support. Once a release goes EOL, its
# packages move to archive.debian.org and the live mirrors 404. Repoint
# apt at the archive for releases we know are EOL so this keeps working
# without depending on the CI image being rebuilt on a newer base.
eol_codenames="bullseye"
codename=$(. /etc/os-release && echo "$VERSION_CODENAME")
case " $eol_codenames " in
    *" $codename "*)
        echo "Debian $codename is EOL; repointing apt sources at archive.debian.org"
        for f in /etc/apt/sources.list /etc/apt/sources.list.d/*.list; do
            [ -f "$f" ] || continue
            sed -i \
                -e "s|http://deb.debian.org/debian|http://archive.debian.org/debian|g" \
                -e "s|http://security.debian.org/debian-security|http://archive.debian.org/debian-security|g" \
                "$f"
        done
        echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/99no-check-valid-until
        ;;
esac

apt_update_and_install() {
    apt-get update
    apt-get install -y --no-install-recommends \
        ca-certificates \
        cmake \
        curl \
        gcc \
        libssh2-1-dev \
        libssl-dev \
        make \
        pkg-config \
        zlib1g-dev
}

n=0
until apt_update_and_install; do
    n=$((n + 1))
    if [ "$n" -ge 3 ]; then
        echo "apt-get update/install failed after $n attempts" >&2
        exit 1
    fi
    echo "apt-get update/install failed, retrying ($n/3)..." >&2
    apt-get clean
    sleep 5
done

archive="$build_dir/libgit2.tar.gz"
curl --fail --location --retry 3 "$LIBGIT2_URL" --output "$archive"
echo "$LIBGIT2_SHA256  $archive" | sha256sum --check --strict

mkdir "$build_dir/source"
tar -xzf "$archive" -C "$build_dir/source" --strip-components=1

multiarch=$(gcc -print-multiarch)
libdir=lib
if [ -n "$multiarch" ]; then
    libdir="lib/$multiarch"
fi

cmake \
    -S "$build_dir/source" \
    -B "$build_dir/build" \
    -DBUILD_SHARED_LIBS=ON \
    -DBUILD_TESTS=OFF \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_LIBDIR="$libdir" \
    -DCMAKE_INSTALL_PREFIX=/usr
cmake --build "$build_dir/build" --parallel "$(getconf _NPROCESSORS_ONLN)"
cmake --install "$build_dir/build"
ldconfig

export PKG_CONFIG_PATH="/usr/$libdir/pkgconfig${PKG_CONFIG_PATH:+:$PKG_CONFIG_PATH}"
installed_version=$(pkg-config --modversion libgit2)
if [ "$installed_version" != "$LIBGIT2_VERSION" ]; then
    echo "expected libgit2 $LIBGIT2_VERSION, found $installed_version" >&2
    exit 1
fi

if [ -n "${GITHUB_ENV:-}" ]; then
    echo "PKG_CONFIG_PATH=$PKG_CONFIG_PATH" >> "$GITHUB_ENV"
fi

rm -rf /var/lib/apt/lists/*
echo "installed libgit2 $installed_version"
