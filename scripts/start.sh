#! /bin/bash

chmod +x bot

# ps aux | grep bot

if pgrep bot >> /dev/null; then
  echo "Бот уже запущен"
else
  nohup ./bot &
fi
