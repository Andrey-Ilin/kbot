# kbot
DevOps application from scratch

test working bot [Link To Bot](t.me/AndriiIlin_bot)

## Features
- Telegram bot integration
- Task automation
- Monitoring and alerting

## Prerequisites
- [Go](https://golang.org/doc/install) (version 1.16 or later)

## Installation

### 1. Clone the repository
```sh
git clone https://github.com/Andrey-Ilin/kbot.git
cd kbot
```

### 2. Install dependencies
```
go mod download
```

### 3. Build
```
go build -ldflags "-X="github.com/andrey-ilin/kbot/cmd.appVersion={version}
```
### 4. Export Tele Token
```
export TELE_TOKEN={your-tele-token}
```

### 5. Start
```
./kbot start
```

## Commands
- **Start the bot**: 
```
./kbot start
```
- **Get bot info**:   
```
./kbot start
./kbot
```
- **Version**
```
./kbot version
```




