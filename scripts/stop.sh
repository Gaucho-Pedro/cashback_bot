#! /bin/bash

if pgrep bot >> /dev/null; then
  pkill bot
  echo "Бот остановлен"
fi
