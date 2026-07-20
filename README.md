# GoAdventure

### Table of Contents
- [About](#about)
- [Getting Started](#getting-started)
  - [Requirements](#requirements)
  - [Install](#install)

## About
Go Adventure is a text-based adventure game. This uses my previous personal project [Eorzean Adventure](https://github.com/poupardm-GhostWrath/EorzeanAdventure) as inspiration. I tried to see if it was possible to create this game in Go with an SQL database.

## Getting Started

#### Requirements
- Go 1.26.0
- Docker and Docker-Compose (or alternative: podman and podman-compose)

#### Install
1. Clone repository:
```
git clone https://github.com/poupardm-GhostWrath/GoAdventure
```
2. Change into the project directory:
```
cd GoAdventure
```
3. Copy .env.example
```
cp .env.example .env
```
4. Load .env (If wanted, read through .env file before loading)

**bash**
```
source .env
```

**zsh**
```
export $(grep -v '^#' .env | xargs)
```

**fish**

Need to add this function to ~/.config/fish/functions/load_env.fish
```
function load_env
    for line in (cat $argv | string match -vre '^#' | string match -vre '^$')
        set arr (string split -n -m 1 = $line)
        if test (count $arr) -eq 2
            set -gx $arr[1] $arr[2]
        end
    end
end
```
Then run
```
load_env .env
```
5. Start sql database container
```
docker-compose up -d
```
6. Download Dependencies
```
go mod download
```
7. Change Execute permissions for start-game.sh
```
chmod +x start-game.sh
```
8. Run Game
```
./start-game.sh
```