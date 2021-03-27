# lidi
lidi - a little English-Polish dictionary, MySQL app with REST API and command-line client in Golang

MySQL database definition:
`./database/database.sql`

English-Polish dictionary, 9900 words:
`words.csv`

Serwer:
`go build -o lidi-server ./cmd/web`

Client:
`go build -o lidi-client ./cmd/client`

English word list: [The BNC/COCA headword lists](https://www.wgtn.ac.nz/lals/resources/paul-nations-resources/vocabulary-lists).
Translation: Google Translate, printed dictionaries, memory.