# lidi
lidi - a little English-Polish dictionary, MySQL database + REST API server and command-line client in Golang (with text to speech thanks Google API).

MySQL database definition:
`./database/database.sql`

English-Polish dictionary, 9900 words:
`words.csv`

Serwer:
`go build -o lidi-server ./cmd/web`

```
Usage of ./lidi-server:
  -addr string
    	HTTP network address (default ":8080")
  -dsn string
    	MySQL data source name (default "web:pass@/dictionary")
```

Client:
`go build -o lidi-client ./cmd/client`

```
Usage:
    lidi-client [en|pl] word

  Subcommands: 
    en      translate from English to Polish
    pl      translate from Polish to English
    speak   say in English (Google API and mplayer is used)

  Positional Variables: 
    word   word to translate (Required)

  Flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -s --server    dictionary server address (default: http://localhost:8080)
    -p --speak     speak English after translate

  Example:
    1. lidi-server
    2. lidi-client en house -p 
```

Speak command (and -p flag) use [htgo-tts](https://github.com/hegedustibor/htgo-tts) lib,
htgo-tts needs mplayer. 

Note: server and client tested on Linux only.

English word list: [The BNC/COCA headword lists](https://www.wgtn.ac.nz/lals/resources/paul-nations-resources/vocabulary-lists).
Translation: Google Translate, printed dictionaries, memory.