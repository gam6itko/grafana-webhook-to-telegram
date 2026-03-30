# grafana-webhook-to-telegram

Accepts a Grafana webhook and forwards the text to Telegram.

## Usage

1. Copy `.env.example` to `.env` and set variables:
   - `BOT_API_KEY_<name>` — bot token from [@BotFather](https://t.me/BotFather). The name in the URL `/api/<name>/...` matches the variable suffix: for `mybot` use `BOT_API_KEY_MYBOT`; hyphens in the URL name become underscores in the variable name (`my-bot` → `BOT_API_KEY_MY_BOT`).
   - `HTTP_SERVER_LISTEN_ADDR` — listen address (Docker image defaults to `0.0.0.0:8080`).
   - `TELEGRAM_API_HOST` — usually `https://api.telegram.org` (optional to change).

2. Grafana **Contact point** (HTTP) URL:

   `POST` or `PUT` to `http://<host>:<port>/api/<bot_name>/<chat_id>`

   `chat_id` is the chat or channel ID (e.g. from [@userinfobot](https://t.me/userinfobot) or the Bot API).

3. Request body is Grafana webhook JSON (`message`, `title`, `status`). Telegram receives `message`, or `title` if `message` is empty.

## Docker build and run

Build the image from the repository root:

```bash
docker build -t grafana-webhook-to-telegram .
```

Run with environment variables:

```bash
docker run --rm -p 8080:8080 \
  -e BOT_API_KEY_MYBOT='your_token' \
  grafana-webhook-to-telegram
```

Override the listen address if needed:

```bash
docker run --rm -p 8080:8080 \
  -e HTTP_SERVER_LISTEN_ADDR=0.0.0.0:8080 \
  -e BOT_API_KEY_MYBOT='your_token' \
  grafana-webhook-to-telegram
```

Locally without Docker: `go run ./cmd` (`.env` is not loaded automatically — export variables or use something like `env $(grep -v '^#' .env | xargs) go run ./cmd`).
