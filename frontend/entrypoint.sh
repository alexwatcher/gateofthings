#!/usr/bin/env sh
set -Ex

find /app/.next \( -type d -name .git -prune \) -o -type f -print0 | xargs -0 sed -i "s#APP_NEXT_PUBLIC_API_URL#$API_URL#g"

exec "$@"
