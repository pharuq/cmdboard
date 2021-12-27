# cmdboard

`cmdboard` is a simple command line tool for storing shell commands. 
This is especially useful when you cannot follow commands from history, such as when accessing a container.

# Demo


# Requirement
Nothing.

# Installation
# Usage
First run the `init` command.

```sh
cmdboard init
```

Use the `add` command to add a command.
If you want to specify a directory, Use the `-d` option.

```sh
cmdboard add "register command" -d "specify directory"
```

Open and select the list of saved commands.

```sh
cmdboard
```

|  Key   |  Action  |
| ----   | ---- |
|  j     |  Down  |
|  k     |  Up  |
|  Enter | Select / Expand directory |
|  d     | Delete command |
