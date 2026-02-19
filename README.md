# Platformer

A fun 2D platformer game! Jump across platforms, dodge gaps, and reach the finish line across 3 levels. You play as a rolling ball with eyes (or a triangle, or a hexagon -- your choice!). 

## How to play

- **Spacebar** -- Jump
- **A** or **Left Arrow** -- Move left
- **D** or **Right Arrow** -- Move right
- **Tab** -- Change your character's shape

## How to set up and run the game on your computer

Follow these steps one at a time. If a step works, move on to the next one.

### Step 1: Install Go (the programming language)

You only need to do this once. Open the **Terminal** app on your computer (ask a parent or teacher to help you find it), then type this and press Enter:

```
brew install go
```

Wait for it to finish. When you see your blinking cursor again, type this to make sure it worked:

```
go version
```

You should see something like `go version go1.25.7`. If you see that, you're good! If it says "command not found", ask a parent or teacher for help.



### Step 2: Create a "Code" directory

On MacOS, create a new directory named "Code" in the user's root directory. Very important - DO NOT put the "Code" folder in either "Downloads" or "Documents". 

```
cd ../  # change to parent directory 
pwd   # shows your current path
ls -al   # lists files at your current path
mkdir Code   # makes a "Code" folder in your current path
cd Code  # change directory to "Code" (only when listed by `ls -al` command)  
```

Once you have created a "Code" directory in your user's root folder, `/users/YOUR-USERNAME/Code/`, move onto next step.

### Step 3: Download the game code

Type this in Terminal (replace the URL with the real one):

```
git clone https://github.com/AAAHHHXXX/platform-game-one.git
```

Then go into the game folder:

```
cd platform-game-one
```

If you already have the folder on your computer, just open Terminal and type `cd ` (with a space after it), then drag the game folder into the Terminal window and press Enter.

### Step 4: Download the game's helper files

Type this and press Enter:

```
go mod tidy
```

This downloads some extra code the game needs. It might take a minute. Wait until you see the blinking cursor again.

### Step 4: Run the game!

Type this and press Enter:

```
go run ./cmd/game
```

A window should pop up with the game! Use **Spacebar** to jump and **A**/**D** or the **Arrow Keys** to move. Try to reach the gold box at the end of each level.

### If something goes wrong

- **"command not found: go"** -- Go isn't installed yet. Go back to Step 1.
- **"command not found: git"** -- You need to install Git. Type `brew install git` and try again.
- **"command not found: brew"** -- You need to install Homebrew first. Ask a parent or teacher to go to https://brew.sh and follow the instructions there.
- **The game window doesn't appear** -- Make sure you're in the right folder. Type `ls` and press Enter. You should see files like `go.mod` and a `cmd` folder. If you don't, go back to Step 2.
- **Still stuck?** -- Ask a parent, teacher, or someone who knows about computers for help. Show them this README file!

## Features

- 3 levels that get harder as you go
- 3 character shapes to choose from (press Tab to switch)
- If you fall off the bottom, you come right back to the start of that level
- Reach the gold goal at the end of each level to move to the next one
- Beat all 3 levels to win!
