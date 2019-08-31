# Nintaco Go API - Tetris Bot Example

### About

[The Nintaco NES/Famicom emulator](https://nintaco.com/) provides [a Go API](https://github.com/meatfighter/nintaco-go-api) that enables externally running programs to control the emulator at a very granular level. This example is an AI that plays tetris.

### Launch

1. Start Nintaco and open Nintendo Tetris.
2. Open the Start Program Server window via Tools | Start Program Server...
3. Press Start Server.
4. From the command-line, launch this Go program.
5. From the GAME TYPE menu, select A-TYPE. Press Start to advance to the next menu.
6. From the A-TYPE menu, select any starting LEVEL. Finally, press Start to let the AI take over.

As an experiment, press the Stop Server button on the Program Server Controls window in the middle of a Tetris game. Let a bunch of pieces drop to grow the stack. Then, press Start Server to see if the AI can clean up the mess.

Alternatively, from the GAME TYPE menu, select B-TYPE. From the B-TYPE menu, select LEVEL 9 and HEIGHT 5, the most difficult combination. Press Start to see if the AI is up to the task. The Tetris Bot actually provides a convenient way to view all the B-TYPE endings.

To make it play faster, from the command-line, break the program (Ctrl+C). Then launch it again, but this time with an argument: `fast`.