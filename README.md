![some-dead-black](https://github.com/user-attachments/assets/7b524627-f437-483e-84f3-46acf368e0bc)

# sl

A CLI built for managing Pokémon `Soul Link` play throughs.

## Setup and Installation

Download the correct `.zip` file for your OS from the [Latest Release](https://github.com/ieedan/sl/releases/latest).

Extract the contents and place them somewhere on your computer.

Copy the path where the `sl` binary lives `ex: C:/soul-link` and set your `PATH` variable.

To test your installation open your terminal and run:

```bash
sl list
```

You should see the following output:

```

  │ Name │ Players │ Started │
  ├──────┼─────────┼─────────┤

```

If you can see the above output that means the CLI is working and has successfully created and migrated the database.

## Using sl

To see the available commands you can run:

```
sl --help
```

### Create a new game

To create a new game, type `sl new` followed by a name for the game.

```
sl new [game-name]
```

It will then prompt you to add `Trainers` (Players) and their starters.

### Return to an existing game

```
sl resume [game-name]
```

### List existing games

```
sl list
```

### Playing the game

```

  │ Route   │ Johnathy │ Jimothy  │
  ├─────────┼──────────┼──────────┤
  │ Starter │ Tepig    │ Oshawatt │

Waiting for command (catch, kill, end, quit, help)...
```

While playing the game you manage it through commands (`catch`, `kill`, `end`).

For help with these commands you can type `help`.

To exit the game you can type `quit`.

