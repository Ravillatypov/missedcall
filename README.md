Скрипт предназначен для уведомления о пропущенных вызовах
уведомления отправляются через SMS/Telegram

Сборка
go get https://github.com/Ravillatypov/missedcall
cd $GOPATH/src/github.com/Ravillatypov/missedcall
go build main.go

Установка

cp $GOPATH/bin/missedcall /usr/local/bin/
touch /var/log/missedcall.log
chown asterisk /var/log/missedcall.log
echo "* * * * * asterisk /usr/local/bin/missedcall >> /var/log/missedcall.log 2&>1" >> /etc/crontab

Настройка

