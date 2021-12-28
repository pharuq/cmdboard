# cmdboard

`cmdboard` is a simple command line tool for storing shell commands. 
This is especially useful when you cannot follow commands from history, such as when accessing a container.

# Demo


# Requirement
Nothing.

# Installation
[Downloads](https://github.com/pharuq/cmdboard/releases)

# Usage
1. (optional)Specify the save destination of the file for saving the command with the environment variable `CMDBOARD_STORED_FILE_PATH`. If not specified, it will be saved in $HOME.
Please be aware of the following points.
- `~` cannot be used
- Please specify the file name with the extension `.json`

```sh
echo export 'CMDBOARD_STORED_FILE_PATH=$HOME/.cmdboard.json' >> ~/.bashrc
```

2. Run the `init` command.

```sh
cmdboard init
```

3. Use the `add` command to add a command.
If you want to specify a directory, Use the `-d` option.

```sh
cmdboard add "register command" -d "specify directory"
```

4. Open and select the list of saved commands.
If you want to copy the selected command to the clipboard, add the `-c` option.

```sh
cmdboard -c
```

|  Key   |  Action  |
| ----   | ---- |
|  j     |  Down  |
|  k     |  Up  |
|  Enter | Select / Expand directory |
|  d     | Delete command |
