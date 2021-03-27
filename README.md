# lidi
lidi - a little English-Polish dictionary, MySQL app with REST API and command-line client in Golang

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
    lidi [en|pl] word

  Subcommands: 
    en   translate from English to Polish
    pl   translate from Polish to English

  Positional Variables: 
    word   word to translate (Required)

  Flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    --s --server    dictionary server address (default: http://localhost:8080)
```

English word list: [The BNC/COCA headword lists](https://www.wgtn.ac.nz/lals/resources/paul-nations-resources/vocabulary-lists).
Translation: Google Translate, printed dictionaries, memory.