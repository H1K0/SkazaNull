#!/bin/bash

if [[ $EUID -ne 0 ]]; then
	echo "Ты ридмишку не читал что ли? Сказано же русским языком: \"ПОД РУТОМ запускать\"!" >&2
	exit 1
fi

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

echo "Идёт установка..." &&

echo "Чудодейство конфигурации..." &&
	mkdir -p /etc/skazanull &&
	cp $SCRIPT_DIR/web/web.conf.yml /etc/skazanull/web.conf.yml &&
	chmod 640 /etc/skazanull/web.conf.yml &&
	chown www-data:www-data /etc/skazanull/web.conf.yml &&

echo "Сборка..." &&
	go build -C $SCRIPT_DIR/web -o $SCRIPT_DIR/bin/skazanull $SCRIPT_DIR/web/cmd/main.go &&
	mkdir -p /opt/skazanull/bin &&
	cp $SCRIPT_DIR/bin/skazanull /opt/skazanull/bin/skazanull &&
	chmod 755 /opt/skazanull/bin/skazanull &&

echo "Установка сервиса в systemctl..." &&
	cp $SCRIPT_DIR/web/snw.service /etc/systemd/system &&
	chmod 644 /etc/systemd/system/snw.service &&
	systemctl daemon-reload &&
	systemctl start snw &&

echo "СказаНулл успешно установлен."
