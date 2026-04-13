#!/bin/sh
set -e

atlas schema apply \
  --auto-approve \
  --url "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@db:3306/${MYSQL_DATABASE}" \
  --to "file:///schema/schema.my.hcl"
