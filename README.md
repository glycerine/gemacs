gemacs: an emacs written in Go
==============================

`gemacs` is derived from `godit`, https://github.com/nsf/godit,
and adapted to use https://github.com/gdamore/tcell for portability.

Release v1.0.2 improves file auto-completion with tab and enter.
Ctrl-x Ctr-w does write file, as in traditional emacs.
View operations is rebound to Ctrl-x Ctrl-v.

Quoting from the original godit readme, with the program name updated
to avoid confusion.

Screenshots:

 * http://nosmileface.ru/images/godit-linux1.png
 * http://nosmileface.ru/images/godit-linux2.png


Gemacs is an emacs-ish lightweight text editor.

Gemacs uses many of the emacs key
bindings and operates using a notion of "micromodes". It's easier to explain
what a micromode is by a simple example. Let's take the keyboard macros feature
from both emacs and gemacs. You can start recording a macro using `C-x (` key
combination and then when you're ready to start repeating it, you do the
following: `C-x e (e...)`. Not only `C-x e` ends the recording of a macro, it
executes the macro once and enters a micromode, where typing `e` again, will
repeat that action. But as soon as some other key was pressed you quit this
micromode and everything is back to normal again. The idea of micromode is used
in gemacs a lot.

# List of keybindings (updated for gemacs, * means new)

~~~

Basic things:
  C-g              - Universal cancel button
  C-x C-c          - Quit from the gemacs
  C-x C-s          - Save file [prompt maybe]
  C-x S            - Save file (raw) [prompt maybe]
  C-x M-s          - Save file as [prompt]
  C-x M-S          - Save file as (raw) [prompt]
  C-x C-w          - Save file as (raw) [prompt] *
  C-x C-f          - Open file
  M-g              - Go to line [prompt]
  C-/              - Undo
  C-x C-/ (C-/...) - Redo

View/buffer operations:
  C-x C-v          - View operations mode *
  C-x 0            - Kill active view
  C-x 1            - Kill all views but active
  C-x 2            - Split active view vertically
  C-x 3            - Split active view horizontally
  C-x o            - Make a sibling view active
  C-x b            - Switch buffer in the active view [prompt]
  C-x k            - Kill buffer in the active view

View operations mode:
  v                - Split active view vertically
  h                - Split active view horizontally
  k                - Kill active view
  C-f, <right>     - Expand/shrink active view to the right
  C-b, <left>      - Expand/shrink active view to the left
  C-n, <down>      - Expand/shrink active view to the bottom
  C-p, <up>        - Expand/shrink active view to the top
  1, 2, 3, 4, ...  - Select view

Cursor/view movement and text editing:
  C-f, <right>     - Move cursor one character forward
  M-f              - Move cursor one word forward
  C-b, <left>      - Move cursor one character backward
  M-b              - Move cursor one word backward
  C-n, <down>      - Move cursor to the next line
  C-p, <up>        - Move cursor to the previous line
  C-e, <end>       - Move cursor to the end of line
  C-a, <home>      - Move cursor to the beginning of the line
  C-v, <pgdn>      - Move view forward (half of the screen)
  M-v, <pgup>      - Move view backward (half of the screen)
  C-l              - Center view on line containing cursor
  C-s              - Search forward [interactive prompt]
  C-r              - Search backward [interactive prompt]
  C-j              - Insert a newline character and autoindent
  <enter>          - Insert a newline character
  <backspace>      - Delete one character backwards
  C-d, <delete>    - Delete one character in-place
  M-d              - Kill word
  M-<backspace>    - Kill word backwards
  C-k              - Kill line
  M-u              - Convert the following word to upper case
  M-l              - Convert the following word to lower case
  M-c              - Capitalize the following word
  <any other key>  - Insert character

Mark and region operations:
  C-<space>        - Set mark
  C-x C-x          - Swap cursor and mark locations
  C-x > (>...)     - Indent region (lines between the cursor and the mark)
  C-x < (<...)     - Deindent region (lines between the cursor and the mark)
  C-x C-r          - Search & replace (within region) [prompt]
  C-x C-u          - Convert the region to upper case
  C-x C-l          - Convert the region to lower case
  C-w              - Kill region (between the cursor and the mark)
  M-w              - Copy region (between the cursor and the mark)
  C-y              - Yank (aka Paste) previously killed/copied text
  M-q              - Fill region (lines between the cursor and the mark) [prompt]

Advanced:
  M-/              - Local words autocompletion
  C-x C-a          - Invoke buffer specific autocompletion menu [menu]
  C-x (            - Start keyboard macro recording
  C-x )            - Stop keyboard macro recording
  C-x e (e...)     - Stop keyboard macro recording and execute it
  C-x =            - Info about character under the cursor
  C-x !            - Filter region through an external command [prompt]

~~~

License: MIT License.
