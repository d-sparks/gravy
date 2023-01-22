#!/bin/bash

git grep 175.. | grep "\.go" | grep -o "175[0-9][0-9]" | sort | uniq
