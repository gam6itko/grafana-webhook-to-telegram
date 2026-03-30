Необходимо написать сервис-адаптер преобразующий grafana webhook payload в в telegram api sendMessage.



Нужно поднять HTTP server, который биндится на HTTP_SERVER_LISTEN_ADDR.
HTTP_SERVER_LISTEN_ADDR - ENV переменная, по умолчанию 127.0.0.1:8080.

Имеется handler `/api/<bot_name>/<chat_id>`. 
- bot_name - динамическая переменна.
- chat_id- динамическая переменная.

На входе у нас будет POST или PUT запрос c body подобного вида.

```json
{
  "receiver": "test",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "TestAlert",
        "grafana_folder": "Test Folder",
        "instance": "Grafana"
      },
      "annotations": {
        "summary": "Notification test"
      },
      "startsAt": "2026-03-30T09:15:43.485120482Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "?orgId=1",
      "fingerprint": "326ea703b01f6100",
      "silenceURL": "https://grafana.bestchange.dev/alerting/silence/new?alertmanager=grafana&matcher=alertname%3DTestAlert&matcher=grafana_folder%3DTest+Folder&matcher=instance%3DGrafana&orgId=1",
      "dashboardURL": "https://grafana.bestchange.dev/d/dashboard_uid?from=1774858543485&orgId=1&to=1774862143486",
      "panelURL": "https://grafana.bestchange.dev/d/dashboard_uid?from=1774858543485&orgId=1&to=1774862143486&viewPanel=1",
      "values": {
        "B": 22,
        "C": 1
      },
      "valueString": "[ var='B' labels={__name__=go_threads, instance=host.docker.internal:3000, job=grafana} value=22 ], [ var='C' labels={__name__=go_threads, instance=host.docker.internal:3000, job=grafana} value=1 ]",
      "orgId": 1
    }
  ],
  "groupLabels": {
    "alertname": "TestAlert",
    "grafana_folder": "Test Folder",
    "instance": "Grafana"
  },
  "commonLabels": {
    "alertname": "TestAlert",
    "grafana_folder": "Test Folder",
    "instance": "Grafana"
  },
  "commonAnnotations": {
    "summary": "Notification test"
  },
  "externalURL": "https://grafana.bestchange.dev/",
  "version": "1",
  "groupKey": "test-326ea703b01f6100-1774862143",
  "truncatedAlerts": 0,
  "orgId": 1,
  "title": "[FIRING:1] TestAlert Test Folder Grafana ",
  "state": "alerting",
  "message": "**Firing**\n\nValue: B=22, C=1\nLabels:\n - alertname = TestAlert\n - grafana_folder = Test Folder\n - instance = Grafana\nAnnotations:\n - summary = Notification test\nSource: ?orgId=1\nSilence: https://grafana.bestchange.dev/alerting/silence/new?alertmanager=grafana&matcher=alertname%3DTestAlert&matcher=grafana_folder%3DTest+Folder&matcher=instance%3DGrafana&orgId=1\nDashboard: https://grafana.bestchange.dev/d/dashboard_uid?from=1774858543485&orgId=1&to=1774862143486\nPanel: https://grafana.bestchange.dev/d/dashboard_uid?from=1774858543485&orgId=1&to=1774862143486&viewPanel=1\n"
}
```

Описание формата: https://grafana.com/docs/grafana/latest/alerting/configure-notifications/manage-contact-points/integrations/webhook-notifier/

Необходимо отправить POST сообщение вида:

```json
{
  "chat_id": "то что пришло в url-path сhat_id",
  "text": "тут данные из grafana message"
}
```


На адрес `TELEGRAM_API_HOST/bot<bot_api_key>/sendMessage`

TELEGRAM_API_HOST - ENV переменная. По-умолчанию: `https://api.telegram.org`

bot_api_key - Берётся из сервиса APIKeyStorage.

## APIKeyStorage

```go
type APIKeyStorage interface {
    Get(ctx context.Context, name string) (token string, ok bool)
}
```

### Реализация APIKeyENVStorage

Нужно слепить переменную вид `BOT_API_KEY_<name>` и отдать её наружу.

Пример:
    Если запросили `storaage.Get(ctx, "foo")`, то нужно отдать значение переменной `BOT_API_KEY_FOO`.

Если такой переменной нет, то возвращаем ok=false.

В вызывающем handler должен быть записан лог с ошибкой (используй: zap), а в ответе должен присутствовать код 404. `{"error": "bot with name foo not found"}`

## logs

на каждый запрос необходимо писать info log о новом входящем сообщении. Где нужно указать bont_name, grafana.title, grafana.status
