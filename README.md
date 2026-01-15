# ironroll

[![Build and Deploy](https://github.com/mtzvd/ironroll/actions/workflows/deploy.yml/badge.svg)](https://github.com/mtzvd/ironroll/actions/workflows/deploy.yml)

Dice roller for [Ironsworn](https://www.ironswornrpg.com/) tabletop RPG. Supports Telegram inline bot, Discord slash command, and HTTP API.

## Ironsworn Dice Mechanic

An Ironsworn roll consists of:
- One action die (1d6)
- Two challenge dice (1d10, 1d10)
- An optional integer modifier

The action score (action die + modifier) is compared against each challenge die:
- Beat both challenge dice: **Strong Hit**
- Beat one challenge die: **Weak Hit**
- Beat neither: **Miss**

When both challenge dice show the same value:
- Strong Hit with doubles: **Strong Hit with a Match**
- Miss with doubles: **Miss with a Match**

## Usage

### Telegram

Use the bot inline in any chat:

```
@ironrollbot +2
@ironrollbot 0
@ironrollbot -1
```

### Discord

Use the slash command:

```
/ironroll modifier:2
/ironroll modifier:0
/ironroll modifier:-1
```

### HTTP API

```bash
curl "https://your-host/roll?m=2"
```

Response:

```json
{
  "action_die": 4,
  "modifier": 2,
  "challenge_dice": [3, 7],
  "total": 6,
  "outcome": "Weak Hit"
}
```

## Installation

```bash
go install github.com/mtzvd/ironroll/cmd/ironroll@latest
```

Or build from source:

```bash
git clone https://github.com/mtzvd/ironroll.git
cd ironroll
go build -o ironroll ./cmd/ironroll
```

## Configuration

Set environment variables or create a `.env` file:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
DISCORD_BOT_TOKEN=your_discord_bot_token
PORT=8080
```

Both bot tokens are optional. The service will start with only the configured adapters.

## Running

```bash
./ironroll
```

Or with [Task](https://taskfile.dev/):

```bash
task run
```

## Project Structure

```
ironroll/
├── cmd/ironroll/      # Application entry point
├── core/roll/         # Pure dice logic (no external dependencies)
├── adapters/
│   ├── telegram/      # Telegram inline bot
│   ├── discord/       # Discord slash command
│   └── httpapi/       # HTTP API handler
├── ratelimit/         # In-memory rate limiter
└── util/
    ├── env/           # .env file loader
    ├── logging/       # Colorized slog handler
    └── random/        # Shared RNG
```

## License

MIT
