# Go Adventure

### Table of Contents
- [Description](#description)
- [Motivation](#motivation)
- [Installation and Usage](#installation-and-usage)
    - [Requirements](#-requirements)
    - [Install](#-installation)
    - [Quick Start](#-quick-start)
    - [Usage](#-usage)
- [Contributing](#-contributing)

## Description

Go Adventure is a text-based adventure game that is played in the terminal.

## Motivation

While I was learning on [Boot.Dev](https://www.boot.dev), we had to create a personal project. I built a text-based adventure game based in the world of Eorzea called [Eorzean Adventure](https://github.com/poupardm-GhostWrath/EorzeanAdventure). I used Python for the project. And during that time, I learned about Go and SQL and was wondering if I could redo that project in Go and SQL. I wanted to see how far I could push myself to re-create that game and improve upon it.

## Installation and Usage

### 🧩 Requirements
- Go 1.26.0
- Docker and Docker-Compose (or alternative like podman and podman-compose)

### ⚙️ Installation
1. Clone repo and Change Directory
```
git clone https://github.com/poupardm-GhostWrath/GoAdventure && cd GoAdventure
```
2. Rename .env.example and load variables into terminal
```
mv .env.example .env  #Renames the file

bash: source .env
zsh: export $(grep -v '^#' .env | xargs)
fish (Why so complicated): #Need to add the following function to ~/.config/fish/load_env.fish

function load_env
    for line in (cat $argv | string match -vre '^#' | string match -vre '^$')
        set arr (string split -n -m 1 = $line)
        if test (count $arr) -eq 2
            set -gx $arr[1] $arr[2]
        end
    end
end

Then: load_env .env
```
3. Start up SQL Database
```
docker-compose up -d
```
4. Download dependencies
```
go mod download
```
5. Change execution permission
```
chmod +x start-game.sh
```
### 🚀 Quick Start

Starting the game (This will build the game and start it):
```
./start-game.sh
```

### 📖 Usage

Available commands:
- **'exit'** or **'quit'** - Exit the game
- **'help'** - Displays the help menu
- **'stat'** - Displays your stats
- **'inv'** or 'inventory' - Displays your inventory
- **'look'** - Displays information about your surrounding
- **'move *\<direction>*'** - Moves you to the location in the direction specified.
- **'store'** - Enters the store that are available in towns.

## 🤝 Contributing

I would love your help! Contribute by forking the repo and opening pull requests. Please avoid spamming the pull request as I as can only review them when I get a chance.
