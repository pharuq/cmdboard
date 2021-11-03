# cmdboard

`cmdboard` is a simple command line tool for storing shell commands. 
It is especially useful in the following use cases:
- Record a series of operations as a group like a runbook.
- Select a command from the list when you cannot follow the command from history, such as when accessing a container.
- When you want to make a note of the command

# Demo
![cmdboard_demo](https://user-images.githubusercontent.com/33982301/147541938-7965d784-9dec-4cdb-a823-00f143456c3b.gif)

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
If you want to add a comment, add the `-c` option.

```sh
cmdboard add "register command" -d "specify directory" -c "write description"
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
