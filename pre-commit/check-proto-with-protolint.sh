#!/usr/bin/env bash

args=(
  "lint"
  "-fix"
)

for file_with_path in "$@"; do
  protolint "${args[@]}" "${file_with_path}"
done
