#!/bin/bash
home=$(cd $(dirname $0);pwd)
chmod +x *.sh
cp odin.service /usr/lib/systemd/system/
sed -i "s#PATH#${home}#g" /usr/lib/systemd/system/odin.service
systemctl enable odin
systemctl start odin
