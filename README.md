# Evil: The Simple Text Editor

I'm making Evil for myself, with all the features that I need and no
more. I won't claim that it'll be good for anyone else. I'm writing it
in Go becouse it needs to be relativly fast but still in a reasonably
high-level language.

## Must-have features:

* VIM keybindings (modal too)
* Split-window
* Emacs keybindings for window management:
    * <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>0</kbd>: Kill window
    * <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>3</kbd>: New vertical window
    * <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>2</kbd>: Now horizontal window
    * <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>o</kbd>: Next window
* Ability to open any repl command in a split window and have it work nicely
* Simple keyword syntax highlighting
* Basic in-file word tab completion
* As-you-type menu of completions from language plugin
* Syntax highlighting and language support plugins written in Lua
* Lua repl available in editor
* Auto parenthises, brace, and quote matching
* Line wrapping at set border
* On bottom bar:
    * Which language plugin is running
    * Which line and character you are on
	* What percentage of the way through your file are you.
	* File name
	* What mode you are in (NORMAL, INSERT, VISUAL)
* Sidebar with files in directory (clickable)

## Finished featrues:

- [ ] VIM keybindings (modal too)
- [ ] Split-window
- [ ] Emacs keybindings for window management:
    - [ ] <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>0</kbd>: Kill window
    - [ ] <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>3</kbd>: New vertical window
    - [ ] <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>2</kbd>: Now horizontal window
    - [ ] <kbd>ctrl</kbd>+<kbd>x</kbd> <kbd>o</kbd>: Next window
- [ ] Ability to open any repl command in a split window and have it work nicely
- [ ] Simple keyword syntax highlighting
- [ ] Basic in-file word tab completion
- [ ] As-you-type menu of completions from language plugin
- [ ] Syntax highlighting and language support plugins written in Lua
- [ ] Lua repl available in editor
- [ ] Auto parenthises, brace, and quote matching
- [ ] Line wrapping at set border
- [ ] On bottom bar:
    - [ ] Which language plugin is running
    - [ ] Which line and character you are on
	- [ ] What percentage of the way through your file are you.
	- [ ] File name
	- [ ] What mode you are in (NORMAL, INSERT, VISUAL)
- [ ] Sidebar with files in directory (clickable)
