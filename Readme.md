# Readme

In this challenge, you are provided with the source code for a distributed MIDI-based music player.
The system consists of two main components:
    Conductor: The conductor is responsible for sending musical notes to the musicians.
    Musician: Musicians receive the notes from the conductor and play them in sequence.

The current solution cannot play the music correctly. Your challenge is to change whatever you need to make the music play correctly.

We are going to test the solution using the queen bohemian rhapsody music `queen.mid`

1. run in bash: midi synthesizer: example: `fluidsynth`
2. run in bash: `make run`
3. you can send a music using: `./build/cli play beeth.mid`

To stop you can use `make stop`

Note: You may need to run this in a unix environment and you may need to install the following depedency: `libportmidi-dev`