# @telefeedbot

Бот в телеграм для чтения ленты вконтакте и рсс лент
Посмотреть в работе можно тут: t.me/telefeedbot

# Description
http://telegra.ph/telefeedbot-05-12


# Как это работает

Бот состоит из трех частей:
 - Сервер базы данных boltsrv
 - Сервер телеграм tgsrv
 - Сервер публикации postsrv

Для хранения опубликованных ссылок используется redis. Можно запустить в докере с сохранением данных на диск хостмашины. Важно: запускть нужно от своего пользователя, а не от root. + нужно убедиться, что id пользователя больше 1000, т.к. "Setting the UID to 1000 ensures we will not run in permissions issues when mapping volumes from our computer to the running container, once 1000 is the first UID assigned to a non root user in Linux, at least in Debian and Ubuntu" 
https://github.com/mhart/alpine-node/issues/48#issuecomment-370171836
```
docker run --name tfredis --restart unless-stopped --user $(id -u) -p 127.0.0.1:6379:6379 -v $(pwd)/redisdata:/data -d redis redis-server --appendonly yes
```

При необходимости прописываем переменные:
```
export GOPATH=~/.go
export GOBIN=$GOPATH/bin
```

Переходим в директорию проекта, ставим зависимости:
```
go get -v -d -t ./...
```

Собираем сервера:

```
cd boltssrv
go install

cd tgsrv 
go install

cd postsrv 
go install
```

Переходим в $GOBIN.
Стартуем: ./boltsrv &

Поднимется http интерфейс на 5000 порту к базе данных (в качестве движка испоьзуется boltdb)
Возможно потребуется поднять лимит одновременно открытых соединений - ulimit

Можно запустить в докере. Dokerfile - boltsrv/Dockerfile.

```
docker run -d --name boltsrv --restart unless-stopped --user $(id -u) -p 127.0.0.1:5000:5000 -v $(pwd)/bolt.db:/app/bolt.db ds0102/boltsrv:unstable
```
bolt.db - файл данных базы, лучше, чтобы был на диске, потому что иначе - если закрашится контейнер - потеряете всю инфу о подписках.


Заводим в телеграм бота.
Кладем в корень файл telefeed.bot c токеном нашего бота

Стартуем сервер телеграм: ./tgsrv>>tgsrv.log &

После этого идем в бота и активируем его - пробуем подписаться на что нибудь, например шлем боту vc.ru
Вобщем как обычно в @telefeedbot


Затем нам надо запустить сервер парсинга и рассылки

Стартуем  так: 
```
./postsrv publics>>publics.log &
./postsrv feeds>>feeds.log &
```
Не забыть поднять редис предварительно - либо заменив редис на хранение в инмемори, например, если редис претит

Для вконтакте - получите токен. см https://github.com/recoilme/tf/tree/master/vkapi
