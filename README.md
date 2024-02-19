# Telegram-бот для отправки файлов в облачное хранилище

## Сценарий использования:
1. Пользователь пишет боту /start (или нажимает на соответствующую кнопку) в ответ бот присылает ссылку для авторизации на yandex-диске.
2. Пользователь переходит по полученной ссылке в браузер и там авторизуется, в случае успешной авторизации пользователь видит соответствующее сообщение, в случае ошибки - видит ошибку. После успешной авторизации бот получает доступ только к своей папке на yandex-диске пользователя, просматривать и изменять другие файлы пользователя он не может.
3. Далее пользователь переходит обратно в Telegram и присылает боту файл с изображением - бот сохраняет его в облачном хранилище, предварительно создав там папку с названием в виде текущей даты. Все изображения от текущей даты попадают в эту папку. В текущей версии возможна отправка только изображений.

## Реализация
#### Кратко:      
Весь проект представляет собой кластер Kubernetes с несколькими сервисами:
   - сервис бота (для взаимодействия с telegram),
   - сервис авторизации (для взаимодействия с oauth.yandex.ru),
   - сервис работы с хранилищем (для взаимодействия с cloud-api.yandex.net),
   - брокер сообщений RabbitMQ (для внутреннего взаимодействия внутри кластера kubernetes),
   - сервер раздачи изображений (для тестирования).
   - prometheus (для сбора и анализа метрик)

<strong>Сервис бота</strong> - получает файлы изображений (в виде url-ссылок) от пользователя и отправляет их в брокер сообщений (в виде url-ссылок); в ответ на /start отдаёт пользователю строку авторизации.     
<strong>Сервис авторизации</strong> - формирует строку авторизации, хранит токены полученные от oauth.yandex.ru в памяти (ОЗУ) с привязкой к ID telegram-пользователя.       
<strong>Сервис работы с хранилищем</strong> - получает из брокера сообщений url-ссылки на изображения, скачивает изображения по ссылкам и отправляет их в хранилище cloud-api.yandex.net.     
<strong>Брокер сообщений RabbitMQ</strong> - используется для повышения отказоустойчивости и горизонтального масштабирования(к одному сервису бота подключено несколько сервисов работы с хранилищем).    
<strong>Сервер раздачи изображений</strong> - реализован на основе nginx, используется для раздачи "статики". Используется только в тестах, на нём хранятся изображения для тестирования.     
<strong>Prometheus</strong> - выполняет сбор метрик с сервиса бота и сервиса работы с хранилищем для последующего анализа результатов нагрузочного тестирования.     

#### Подробно:  
https://docs.google.com/document/d/1D6Lg1bm3j7GlibDTnGk1mUr28Xeg496MhJE61NpWVpw
## Используемые технологии из курса: 
- kubernetes для управления docker-контейнерами и управления количеством запущенных инстансов сервиса; 
- авторизация OAuth 2.0; 
- Prometheus для сбора и анализа метрик; 
- RabbitMQ для балансировки нагрузки через очередь сообщений; 
- RabbitMQ для реализации межсервисного обмена - Event Collaboration паттерна;
- Postman для функционального и нагрузочного тестирования.

## Установка
### На prodaction-сервер:
1. перед применением манифестов необходимо установить rabbitmq командой:    
sudo helm install mq-csysbot oci://registry-1.docker.io/bitnamicharts/rabbitmq --set auth.username='guest',auth.password='guest'
2. установить prometheus командами:
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install stack prometheus-community/kube-prometheus-stack -f <абсолютный путь до файла ./deployments/prometeus_helm/prometheus.yaml>
3. заполнить поля в configmap.yaml манифесте(tgbottoken - токен telegram-бота, ydiskid - ID приложения выданный yandex-ом для приложения работающего с yandex-диском, ydisksecret - ключ приложения выданный yandex-ом для приложения работающего с yandex-диском, storagefolder - путь к папке приложения на yandex-диске в формате "disk:/Приложения/<имя приложения>")
4. применить манифесты из папки /deployments/kubernetes/ командой kubectl apply -f .
5. настроить проброс портов с хоста в кластер кубернетес командой для сервиса авторизации:    
sudo kubectl port-forward svc/auth-service-service 8099:8099 --address 192.168.49.1, где 192.168.49.1 ip-адрес хоста 
6. настроить проброс портов с хоста в кластер кубернетес для работы с дашбордом prometheus
sudo kubectl port-forward service/prometheus-operated  9090

### Для запуска функциональных тестов:
1. установить rabbitmq как указано выше
2. заполнить два поля в configmap.yaml манифесте (debugtoken - токен авторизации выданный yandex-ом для работы с yandex-диском, storagefolder - путь к папке приложения на yandex-диске в формате "disk:/Приложения/<имя приложения>")
3. применить манифесты из папки /deployments/kubernetes/ командой kubectl apply -f .
4. запустить коллекцию тестов из папки /tests/ командой newman run CSYSBot.postman_collection.json 

### Для запуска нагрузочных тестов:
1. установить rabbitmq как указано выше
2. установить prometheus как указано выше
3. поля в configmap.yaml манифесте оставить пустыми
4. применить манифесты из папки /deployments/kubernetes/ командой kubectl apply -f .
5. открыть коллекцию тестов ./tests/CSYSBotHighload.postman_collection.json в Postman и запустить в режиме "оценки производительности"
6. настроить проброс портов с хоста в кластер кубернетес для работы с дашбордом prometheus
sudo kubectl port-forward service/prometheus-operated  9090
7. выполнить анализ метрик
