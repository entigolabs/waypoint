#!/bin/sh

command=${1:-migrate}

case "$command" in
  validate)
    output=$(tern status --migrations /migrations)
    status=$(echo "$output" | grep -i 'status:' | awk '{print $2 " " $3 " " $4}')
    echo "Migration status: $status"
    if [ "$status" = "up to date" ]; then
      echo "Migration is up to date"
      exit 0
    else
      echo "Migration is not up to date"
      exit 1
    fi
    ;;
  migrate)
    tern code install /core
    exit_code=$?
    if [ $exit_code -ne 0 ]; then
      exit $exit_code
    fi
    tern migrate --migrations /migrations
    exit_code=$?
    exit $exit_code
    ;;
  *)
    echo "Invalid argument: $command"
    echo "Usage: $0 {validate|migrate}"
    exit 1
    ;;
esac