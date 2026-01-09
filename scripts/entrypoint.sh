#!/bin/sh
set -e

# Запускаем миграции
echo "Running migrations..."
./migrationsBuild

# Запускаем основное приложение
echo "Starting main application..."
exec ./main