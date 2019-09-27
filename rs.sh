#!/bin/bash

rsync -vvv -r --exclude=dist --exclude=.DS_Store --exclude=.vscode --exclude=.git --exclude=node_modules  -i /Users/dark/.ssh/google_compute_engine  . dark@34.68.11.70:swarm-chaos
