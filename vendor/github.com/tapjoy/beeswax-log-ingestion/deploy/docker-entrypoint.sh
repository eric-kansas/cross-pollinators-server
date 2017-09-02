#!/usr/bin/env bash

set -eux

# BEGIN FUNCTIONS #

function prepare_for_deployment() {
  ###
  # Everything in this script function is project-agnostic. It can be moved to other projects
  # wholesale. There is nothing load bearing in here, only a dependency on a
  # project makefile
  ###
  # Enable shell output glob searches
  shopt -s extglob
  # set the sha for our revision check
  echo `git rev-parse HEAD` > shafile &&\
  #make deps &&\
  make build &&\

  echo "Removing everything but the successfully built binary, the Procfile, and the .git folder..."
  rm -r $(ls -A | grep -v "${PROJECT_NAME:-$(basename $PWD)}$" | grep -v "\.git$" | grep -v "Procfile")
}

function check_dependencies() {
  echo "check_dependencies done"
}

function setup_container() {
  make deps &&\
 echo "setup_container done"
}

function project_exec() {
  check_dependencies &&\
  setup_container &&\
  echo "Container setup finished, running \`$@\`"
  rm -f pids/*
  exec $@
}

# END FUNCTIONS #

"$@"
