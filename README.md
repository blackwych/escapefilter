# escapefilter - ANSI escape code filter

Read the input, process characters including ANSI escape codes and print the result in plain text.


## Usage

```
escapefilter [INFILE...]
```

### Options

* `-h`, `--help`:

  Print usage and exit.

* `-v`, `--version`:

  Print version information and exit.


### Arguments

* `INFILE`:

  Path to input file. Standard input will be used if no files are specified or `INFILE` is `-`.


## Supported ANSI escape codes

Any other unsupported escape code is just ignored.

Tab stops are fixed at multiples of 8, and cannot be changed currently.


### C0 control codes

Codepoint | Abbr. | Name            | Effect
----------|-------|-----------------|--------
U+0008    | BS    | Backspace       | Moves the cursor left. ("backword wrap" is not supported.)
U+0009    | HT    | Horizontal Tab  | Moves the cursor right to the next tabstop.
U+000A    | LF    | Line Feed       | (Behaves as CR+LF) Moves the cursor to the beginning of the next line.
U+000B    | VT    | Vertical Tab    | Moves the cursor to the next line, keeping its column.
U+000D    | CR    | Carriage Return | Moves the cursor the beginning of the line.
U+001B    | ESC   | Escape          | Starts escape sequences.


### C1 control codes

Totally unsupported.


### Escape sequences

Code  | Abbr. | Name                        | Effect
------|-------|-----------------------------|--------
ESC [ | CSI   | Control Sequence Introducer | Starts control sequences.


### Control sequences

Code            | Abbr. | Name                                 | Effect
----------------|-------|--------------------------------------|--------
CSI *n* A       | CUU   | Cusror Up                            | Moves the cursor to *n* line(s) up. *n* defaults to 1.
CSI *n* B       | CUD   | Cursor Down                          | Moves the cursor to *n* lines down. *n* defaults to 1.
CSI *n* C       | CUF   | Cusror Forward                       | Moves the cursor to *n* column(s) forward. *n* defaults to 1.
CSI *n* D       | CUB   | Cursor Backward                      | Moves the cursor to *n* column(s) backward. *n* defaults to 1.
CSI *n* E       | CNL   | Cursor Next Line                     | Moves the cursor to the beginning of *n* line(s) down. *n* defaults to 1.
CSI *n* F       | CPL   | Cursor Previous Line                 | Moves the cursor to the beginning of *n* line(s) up. *n* defaults to 1.
CSI *n* G       | CHA   | Cursor Horizontal Absolute           | Moves the cursor to column *n*. *n* defaults to 1.
CSI *m* ; *n* H | CUP   | Cursor Position                      | Moves the cursor to row *m* column *n*. *m* and *n* defaults to 1.
CSI *n* I       | CHT   | Cursor Horizontal Forward Tabulation | Moves the cursor *n* tab(s) forward. *n* defaults to 1.
CSI *n* J       | ED    | Erase in Display                     | [*n* = 0] Erases characters from the cursor to the end of the screen.<br>[*n* = 1] Erases characters from the beginning of the screen to the cursor.<br>[*n* = 2] Erases all characters in the screen.<br>*n* defaults to 0.
CSI *n* K       | EL    | Erase in Line                        | [*n* = 0] Erases characters from the cursor to the end of the line.<br>[*n* = 1] Erases characters from the beginning of the line to the cursor.<br>[*n* = 2] Erases all characters in the line.<br>*n* defaults to 0.
CSI *n* Z       | CBT   | Cursor Backward Tabulation           | Moves the cursor *n* tab(s) backward. *n* defaults to 1.
