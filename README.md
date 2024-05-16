# GophKeeper

## Установка и настройка СУБД postgress
apt install postgresql

psql (14.11 (Ubuntu 14.11-0ubuntu0.22.04.1))
Type "help" for help.

postgres=# ALTER USER postgres WITH PASSWORD 'postgres';

postgres=# CREATE DATABASE praktikum;

## Установка и настройка s3

### Linux
wget https://dl.min.io/server/minio/release/linux-amd64/minio_20240510014138.0.0_amd64.deb

dpkg -i minio_20240510014138.0.0_amd64.deb

MINIO_ROOT_USER=admin MINIO_ROOT_PASSWORD=password minio server /mnt/data --console-address ":9001"

### Windows

PS> Invoke-WebRequest -Uri "https://dl.min.io/server/minio/release/windows-amd64/minio.exe" -OutFile "C:\minio.exe"

PS> setx MINIO_ROOT_USER admin

PS> setx MINIO_ROOT_PASSWORD password

PS> C:\minio.exe server F:\Data --console-address ":9001"

### Создание Access Key
Сделать через WebUI

S3AccessKeyID     = "aHLytUVhTKOPMYD6nYA2"

S3SecretAccessKey = "F2Avh18pul7X8IsGhCTeWPnaQNhlOuda3iAYSO30"

## Сборка дистрибутива
make build

## Запуск сервера
goph-keeper-windows-server.exe

## Запуск клиента
goph-keeper-windows-client.exe

## Порядок работы

health

regiter user1 user1

logout

login user1 user1 

put obj1 1 data1

get obj1

put obj2 3 inFilename

get obj2 outFilename

exit

