🎬 Movie Tracker Telegram Bot
A Telegram bot for tracking movies you want to watch, built with Go.
Features

➕ Add movies to your watchlist
📋 View your movie list with watch status and ratings
✅ Mark movies as watched
⭐ Rate movies from 1 to 5
🎲 Get a random unwatched movie suggestion
🔍 Search movies by title
🗑️ Delete movies from your list

Tech Stack

Go
telebot.v3
JSON file storage

Setup

Clone the repository
Create a .env file:

BOT_TOKEN=your_telegram_bot_token

Run the bot:

bashgo run main.go
Usage
CommandDescription
/add <title>     Add a movie
/list            Show all movies
/watched <id>    Mark as watched
/rate <id> <1-5> Rate a movie
/delete <id>     Delete a movie
/search <title>  Search by title
