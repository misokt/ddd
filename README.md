# dump-discord-data

## usage:

dump all messages from default path:
```
$ dump-discord-data.exe
```

dump all messages from specified year from specified path:
```
$ dump-discord-data.exe --path /path/to/messages/ --year 2022
```

dump all messages excluding 1 channel:
```
$ dump-discord-data.exe --exclude --channels 123123123123123
```

dump all messages excluding multiple channels:
```
$ dump-discord-data.exe --exclude --channels 234234234234234,345345345345345
```

dump all messages only from specified channels:
```
$ dump-discord-data.exe --include --channels 456456456456456,567567567567567
```

dump all messages only from specified channels from a specified year from a specified path:
```
$ dump-discord-data.exe --include --channels 678678678678678,789789789789789 --year 2023 --path /path/to/messages/
```

## directory structure

- where the executable should be placed to use default path:
```
.
└── package/
    ├── dump-discord-data.exe
    ├── messages/
    └── README.txt
```
- if specifying the path, include "messages" directory in the path
- it is okay to have multiple directories under package. program looks specfically for "messages" directory

## (re)-written in go because:

- easier to distribute a compiled binary that won't require additional installations
- the python script hung on me the first time I used it
- heard complaint the python script took too long with larger amount of message files (10.000+)
- uses flags. if you're going to make a CLI app, take advantage of flags
- my own. so, I can blindly trust it now and any subsequent future changes

## exec time difference:

[(dumping all messages) top is the python script, bottom is this go binary.](screenshots/python-vs-go-exec-time.png)

## credits:

to & inspired by the python script: [ishnz/bulk_deletion_helper](https://github.com/ishnz/bulk_deletion_helper)
