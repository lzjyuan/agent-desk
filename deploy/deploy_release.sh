#!/usr/bin/env bash
set -euo pipefail

root="${1:?remote root required}"
archive="${2:?archive required}"
expected_hash="${3:?checksum required}"
release="${4:?release id required}"
service_user="${5:?service user required}"

[[ $EUID -eq 0 ]] || { echo 'Run as root (or passwordless sudo).' >&2; exit 1; }
[[ $root =~ ^/[A-Za-z0-9._/-]+$ && $root != / && $release =~ ^[0-9]{17}$ ]] || exit 1
[[ /$root/ != *'/../'* && /$root/ != *'/./'* ]] || exit 1
[[ $service_user =~ ^[A-Za-z_][A-Za-z0-9_-]*$ ]] || exit 1
id "$service_user" >/dev/null
service_group="$(id -gn "$service_user")"
echo "$expected_hash  $archive" | sha256sum -c -

install -d -m 755 "$root" "$root/releases"
install -d -o "$service_user" -g "$service_group" -m 750 "$root/shared" "$root/shared/config" "$root/shared/data"
config="$root/shared/config/config.yaml"
[[ -s $config ]] || { echo "Missing $config. Create it from config/config.example.yaml before deploying." >&2; exit 1; }
destination="$root/releases/$release"
[[ ! -e $destination ]] || { echo "Release already exists: $destination" >&2; exit 1; }
install -d -m 755 "$destination"
tar -xzf "$archive" -C "$destination" --no-same-owner
[[ -f $destination/agent-service ]] || exit 1
chmod 755 "$destination/agent-service"

unit=/etc/systemd/system/agent-service.service
if [[ ! -f $unit ]]; then
  cat > "$unit" <<EOF
[Unit]
Description=Agent Service
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
User=$service_user
WorkingDirectory=$root/shared
ExecStart=$root/current/agent-service -config $config
Restart=on-failure
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable agent-service.service
fi

previous=''
if [[ -L $root/current ]]; then previous="$(readlink "$root/current")"; fi
ln -s "$destination" "$root/current.next"
mv -Tf "$root/current.next" "$root/current"
if ! systemctl restart agent-service.service; then
  healthy=false
else
  healthy=false
  for _ in {1..20}; do
    if systemctl is-active --quiet agent-service.service && curl -fsS --max-time 2 http://127.0.0.1:8083/api/config >/dev/null; then
      healthy=true
      break
    fi
    sleep 1
  done
fi
if [[ $healthy != true ]]; then
  echo 'Health check failed; restoring previous release.' >&2
  if [[ -n $previous ]]; then
    ln -s "$previous" "$root/current.next"
    mv -Tf "$root/current.next" "$root/current"
    systemctl restart agent-service.service
  else
    systemctl stop agent-service.service || true
    rm -f "$root/current"
  fi
  exit 1
fi
echo "agent-service $release is healthy"
