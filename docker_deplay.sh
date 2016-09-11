#!/usr/bin/env bash

REV=$(git log --pretty=format:'%h' -n 1)

echo "ic build -t mcp:$REV ."

cf ic build -t mcp:$REV  .

