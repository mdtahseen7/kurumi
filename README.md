# Kurumi-Downloader (forked from govd)

A Telegram bot for downloading media from various platforms.

This project was created after the discontinuation of the popular UVD bot and draws significant inspiration from [yt-dlp](https://github.com/yt-dlp/yt-dlp).

---

## Features

- Download media from various platforms  
- Download videos, photos, and audio  
- Inline mode support  
- Group chat support with customizable settings  
- Media caption support  

---

## Dependencies

- `ffmpeg >= 6.1.1`  
- `libheif >= 1.19.7`  
- `pkg-config`  
- MySQL database  

---

## Installation

```bash
git clone https://github.com/govdbot/govd.git
cd govd
# Edit the .env file with your bot token and database credentials
sh build.sh
```

---

## Installation with Docker

First, build the image using the Dockerfile:

```bash
docker build -t govd-bot .
```

Next, update the `.env` file to ensure the database properties match the environment variables defined for the MariaDB service in the `docker-compose.yml` file.  
While the default values are acceptable, **it is recommended** to change the `MYSQL_PASSWORD` in `docker-compose.yml` for enhanced security, and update the `DB_PASSWORD` in `.env` accordingly.

The following line in the `.env` file **must** be set exactly as shown:

```env
DB_HOST=db
```

Finally, run the following to start all services:

```bash
docker compose up -d
```

---

## Environment Variables

| Variable                | Description                                        | Default                                |
|-------------------------|----------------------------------------------------|----------------------------------------|
| `DB_HOST`              | Database host                                      | `localhost`                            |
| `DB_PORT`              | Database port                                      | `3306`                                 |
| `DB_NAME`              | Database name                                      | `govd`                                 |
| `DB_USER`              | Database user                                      | `govd`                                 |
| `DB_PASSWORD`          | Database password                                  | `password`                             |
| `BOT_API_URL`*         | Telegram bot API URL                               | `https://api.telegram.org`             |
| `BOT_TOKEN`            | Telegram bot token                                 | `12345678:ABC-DEF1234ghIkl-zyx57W2P0s` |
| `CONCURRENT_UPDATES`  | Max concurrent updates handled by the bot          | `50`                                   |
| `LOG_DISPATCHER_ERRORS` | Log dispatcher errors                              | `0`                                    |
| `DOWNLOADS_DIR`       | Directory for downloaded files                     | `downloads`                            |
| `HTTP_PROXY`          | HTTP proxy (optional)                              |                                        |
| `HTTPS_PROXY`         | HTTPS proxy (optional)                             |                                        |
| `NO_PROXY`            | No-proxy domains (optional)                        |                                        |
| `REPO_URL`            | Project repository URL                             | `https://github.com/govdbot/govd`      |
| `PROFILER_PORT`       | Port for profiler HTTP server (pprof)              | `0` _(disabled)_                       |

**Note:**  
To avoid file size limits, it is recommended to host your own Telegram Bot API server. A public instance is currently running under a Bot API fork: [tdlight-telegram-bot-api](https://github.com/tdlight-team/tdlight-telegram-bot-api), but you can also use the official Bot API client.

---

## Cookies

Some extractors require cookies for downloading. To use them, insert a `.txt` file in the `cookies` folder (in Netscape format).

---

## TODO

- [ ] Add more extractors
- [ ] Switch to Sonic JSON parser
- [ ] Switch to native libav
- [ ] Add tests
- [ ] Add Dockerfile and Compose
- [ ] Improve error handling
- [ ] Add support for Telegram webhooks
- [ ] Switch to PostgreSQL (?)
- [ ] Better API (?)
- [ ] Better documentation with multiple README files

