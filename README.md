
# Steam ID Availability Checker

This is a command-line tool written in Go that allows you to check the availability of Steam IDs.
This tool also has a session feature that allows you to save your progress so you dont have to recheck usernames you already checked previously.

This is also my first program in Go, and something i decided to do for fun so the code might not be the best. Feel free to report issues or suggest changes.
## Features

- **Check Steam ID availability**: Check if a Steam ID is available or already claimed.
- **Session management**: Create new sessions or resume existing ones.
- **Progress tracking**: The tool remembers where it left off in case of interruptions.

## Installation

You can download the latest release of the Steam ID Checker from the [releases page](https://github.com/ytax/steam-id-checker/releases) on GitHub. Simply go to the page, choose the latest version, and download the zip.


## Build Instructions

1. Clone the repository:

   ```bash
   git clone https://github.com/ytax/steam-id-checker.git
   ```

2. Navigate to the project folder:

   ```bash
   cd steam-id-checker
   ```

3. Build the project:

   ```bash
   go build
   ```

4. Run the program:

   ```bash
   ./steamidcheck.exe
   ```

## Usage

When you run the tool, you'll be greeted with a menu to select your desired action. You can:

1. **Start a New Session**: Create a new session and begin checking Steam IDs from a `targets.txt` file.
2. **Resume an Existing Session**: Choose an existing session to continue checking IDs.
3. **Exit**: Exit the program.

The program will check each Steam ID from your `targets.txt` file. If an ID is available, it will be saved to the `output.txt` file. The program will also track your progress, so you can stop and resume at any time.

### Example Session Workflow

1. **Start a New Session**:
   - If no session exists, the tool will create a new session automatically and read from the `targets.txt` file.
   - By default, this file contains some random semi-og usernames, you can replace the content of this file with whatever you want the software to check.
   - I also HEAVILY recommend you to run your list through a randomizer so that you arent checking the targets in alphabetic order.
2. **Resume a Session**:
   - If sessions exist, you can select one to resume from where you left off.
   - You will be shown the available sessions, and you can write it's name to select it.

## Known Issues

- The tool may take time depending on the number of IDs being checked. I havent implemented threading to this tool yet.
- If you delete the targets.txt file the tool will stop working because it wont be able to read the targets.
- Sometimes the tool will give false-positives because of shadow-banned accounts.

## To-do List

- [ ] Add steam API Key support, this will prevent the false-positives regarding the shadow-banned accounts.
- [ ] Add threading.
- [ ] Add an option to generate targets (3c, 3l, 4c and 4l usernames).
- [ ] Add discord webhook support, this might not be a good idea because of how often valid IDs are found.
- [ ] Add a turbo to quickly claim usernames. Proxy support might also be needed for this.
- [ ] Add some comments because i was too lazy to do it. I think the code is very easy to understand though.

## Credits

This tool was created by [ytax](https://github.com/ytax).

Feel free to contribute or open issues if you encounter any bugs or have suggestions for new features.

---

## Preview of the interfance

![image](https://github.com/user-attachments/assets/98e9d95d-2f21-4137-a919-1346ace4dc90)

---

Enjoy checking your Steam IDs and have fun!
