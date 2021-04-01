# lidi
lidi - a little English-Polish dictionary, MySQL database + REST API server and command-line client in Golang (with text to speech thanks to Google API).

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
    lidi [en|pl|speak|add|update|delete]

  Subcommands: 
    en       Translate from English to Polish
    pl       Translate from Polish to English
    speak    Say in English (Google API is used)
    add      Add new item to dictionary (English=Polish) returns ID
    update   Update item in dictionary (ID English=Polish)
    delete   Delete item in dictionary (English=Polish)

  Flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -s --server    Dictionary server address (default: http://localhost:8080)
    -p --speak     Speak English after translate (en|pl commands)
    -i --id        show record id


  Example:
    1. Server start: ./lidi-server
    2. Client query: ./lidi-client en house -p
                     ./lidi-client add tree=drzefo
                     ./lidi-client update 16354 tree=drzewo  
                     ./lidi-client delete 16354
```

Speak command (and -p flag) use [htgo-tts](https://github.com/hegedustibor/htgo-tts) lib,
htgo-tts needs mplayer. 

Note: server and client tested on Linux only.

English word list: [The BNC/COCA headword lists](https://www.wgtn.ac.nz/lals/resources/paul-nations-resources/vocabulary-lists).
Translation: Google Translate, printed dictionaries, memory.