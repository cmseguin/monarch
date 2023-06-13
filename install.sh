#!/bin/bash
latest_release_url() {
  echo "https://github.com/cmseguin/monarch/releases/latest"
}

get_source_tarball_url() {
  local version=$1
  echo "https://github.com/cmseguin/monarch/archive/refs/tags/$1.tar.gz"
}

get_lastest_release() {
  local version=$(curl --show-error --location --fail $(latest_release_url) | grep -Eo "tag/[0-9]\\.[0-9]\\.[0-9](-[0-9]+)?" | grep -Eo "[^/]+$" | head -n 1)
  echo $version
}

download_source_tarball() {
  local install_dir=$1
  local version=$2
  local url=$(get_source_tarball_url $version)
  local tmp_dir="$install_dir/.monarch/tmp"
  local file="$tmp_dir/monarch-$version.tar.gz"

  curl --show-error --location --fail --output "$file" --write-out "$file" $url
}

extract_source_tarball() {
  local install_dir=$1
  local version=$2
  local tmp_dir="$install_dir/.monarch/tmp"
  local file="$tmp_dir/monarch-$version.tar.gz"

  tar -xzf $file -C $tmp_dir
  local dirname=$(basename $file .tar.gz)
  echo "$tmp_dir/$dirname"
}

clean_tmp_dir() {
  local install_dir=$1
  local tmp_dir="$install_dir/.monarch/tmp"

  rm -rf "$tmp_dir"
}

create_install_dir() {
  local install_dir=$1

  if [ ! -d "$install_dir/.monarch" ]; then
    mkdir -p "$install_dir/.monarch"
  fi

  if [ ! -d "$install_dir/.monarch/bin" ]; then
    mkdir -p "$install_dir/.monarch/bin"
  fi

  if [ ! -d "$install_dir/.monarch/tmp" ]; then
    mkdir -p "$install_dir/.monarch/tmp"
  fi
}

find_profile_file() {
  local profile_file

  if [ -f "$HOME/.zshrc" ]; then
    profile_file="$HOME/.zshrc"
  elif [ -f "$HOME/.bashrc" ]; then
    profile_file="$HOME/.bashrc"
  elif [ -f "$HOME/.bash_profile" ]; then
    profile_file="$HOME/.bash_profile"
  elif [ -f "$HOME/.profile" ]; then
    profile_file="$HOME/.profile"
  fi

  echo $profile_file
}

set_path_in_profile() {
  local install_dir=$1
  local profile=$2

  if [ -f "$profile" ]; then
    if ! grep -q "export PATH=\"\$PATH:$install_dir/.monarch/bin\"" "$profile"; then
      echo "export PATH=\"\$PATH:$install_dir/.monarch/bin\"" >> "$profile"
    fi
  fi
}

install_monarch() {
  local install_dir=$1
  local version=$(get_lastest_release)

  create_install_dir $install_dir
  download_source_tarball $install_dir $version
  local dirname=$(extract_source_tarball $install_dir $version)

  go build -o "$install_dir/.monarch/bin/monarch" "$dirname/main.go"

  clean_tmp_dir $install_dir
  local profile=$(find_profile_file)
  set_path_in_profile $install_dir $profile

  source $profile
}

install_monarch $HOME
