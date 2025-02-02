# Distributed MIDI Music Player - Gorxestra

## **Challenge Overview**

In this challenge, you are provided with the source code for a distributed MIDI-based music player.  
The system consists of two main components:
- **Conductor:** Responsible for sending musical notes to the musicians.
- **Musician:** Musicians receive notes from the conductor and play them in sequence.

The original solution was unable to play the music correctly due to performance bottlenecks and synchronization issues. The challenge required modifying the code to ensure that music plays without delays or desynchronization.

We will test the solution using Queen's *Bohemian Rhapsody* (`queen.mid`).

### **Steps to Run the System**
1. Launch a MIDI synthesizer (e.g., `fluidsynth`) in the terminal:
   ```bash
   fluidsynth
   ```
2. Run the program:
   ```bash
   make run
   ```
3. Play a MIDI file using the CLI:
   ```bash
   ./build/cli play beeth.mid
   ```

4. Stop playback:
   ```bash
   make stop
   ```

### **System Requirements**
- Run in a UNIX environment.
- Install the following dependency:  
  ```bash
  sudo apt-get install libportmidi-dev
  ```

---

## **My Solution**

### **Changes Implemented**
I made significant modifications to the `baton.go` file to solve the issues outlined in the challenge.

1. **Parallel Note Processing with Go Routines:**  
   The original implementation used a single channel to sequentially send notes to all musicians, causing delays and incorrect timing during playback.  
   - I implemented a new architecture where each musician (representing an instrument with tracks and notes) has a dedicated Go routine that reads notes from its own channel.
   - This parallel processing ensures that multiple notes are played concurrently, improving synchronization and reducing lag.

2. **Dedicated Channel System:**  
   - For each musician, a separate channel is initialized.
   - The conductor now manages these channels through the baton, sending notes to each musician independently and concurrently.
   
3. **Pause and Resume Functionality:**  
   - I added the ability to pause and resume playback:
     - Use `make pause` or send a `SIGUSR1` signal to pause.
     - Use `make resume` or send a `SIGUSR2` signal to resume.

4. **Signal Handling:**  
   - The system listens for `SIGUSR1` (pause) and `SIGUSR2` (resume) signals for real-time control over music playback.

### **Technical Explanation**
- Each musician is represented by a Go routine that reads from a dedicated channel.
- Notes are sent in parallel from the conductor, ensuring instruments play in sync.
- The `baton.go` file includes a signal handler to manage pause and resume commands without interrupting the main playback loop.
- The `play` method initializes channels for tracks, creates Go routines for musicians, and handles synchronization through a `sync.WaitGroup`.

---

## **Testing**
The modified solution was tested using Queen's *Bohemian Rhapsody* MIDI file (`queen.mid`). The system now handles complex tracks with multiple instruments accurately, ensuring proper timing and playback.
