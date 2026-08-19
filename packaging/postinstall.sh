#!/bin/sh
set -e

warn() {
    echo "atoman-mirror: $1"
}

if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload || warn "systemd reload failed"
    systemctl stop hubproxy >/dev/null 2>&1 || true
    systemctl disable hubproxy >/dev/null 2>&1 || true
    systemctl enable atoman-mirror >/dev/null 2>&1 || warn "systemd enable failed"

    if [ -d /run/systemd/system ]; then
        systemctl restart atoman-mirror || systemctl start atoman-mirror || {
            warn "service start failed, check: journalctl -u atoman-mirror"
        }
    fi
fi

if command -v rc-update >/dev/null 2>&1; then
    rc-update add hubproxy default >/dev/null 2>&1 || warn "OpenRC enable failed"
fi

if command -v rc-service >/dev/null 2>&1; then
    rc-service hubproxy restart || rc-service hubproxy start || {
        warn "service start failed, check: rc-service hubproxy status"
    }
fi
