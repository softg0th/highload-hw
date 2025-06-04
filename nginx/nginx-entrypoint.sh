#!/bin/sh
echo "[nginx] waiting for tasks.filter to resolve..."
until getent hosts tasks.filter; do
  sleep 1
done
echo "[nginx] DNS OK. Starting nginx..."
exec nginx -g 'daemon off;'
