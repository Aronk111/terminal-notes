## If you wish to build the project yourself

Make sure you have [go](https://go.dev/dl/) installed!<br>
Navigate to where you want to have the project downloaded to and do:

`git clone https://github.com/Aronk111/terminal-notes/edit/main/README.md`<br>
then `cd` into the downloaded folder

### Windows:
```
mkdir bin
cd ./src
go build -o ../bin/notes.exe
cd ../bin
./notes.exe
```

### Macos *(afaik, I don't have macOS)*:
```
mkdir bin
cd ./src
go build -o ../bin/notes
cd ../bin
./notes
```

To have the *`notes`* command always available just add the **bin** folder to your system environment variables.
