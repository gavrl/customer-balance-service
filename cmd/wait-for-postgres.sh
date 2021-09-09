#!/bin/sh
# wait-for-postgres.sh

set -e

while getopts h:u:p:d:c: flag
do
    case "${flag}" in
        h) host=${OPTARG};;
        u) username=${OPTARG};;
        p) password=${OPTARG};;
        d) db=${OPTARG};;
        c) com=${OPTARG};;
    esac
done

until PGPASSWORD="$password" psql -h "$host" -U "$username" "$db" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $com