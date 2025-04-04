package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"resty.dev/v3"
)

var scrapper = resty.New().SetTimeout(10 * time.Second)

// silly terminal colors
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func checkID(id string) bool {
	url := "https://steamcommunity.com/id/" + id
	resp, err := scrapper.R().Get(url)
	if err != nil {
		fmt.Println(Red+"Request error:"+Reset, err)
		pauseTerminal()
		return false
	}
	body := resp.String()
	return !strings.Contains(body, "The specified profile could not be found.")
}

func pauseTerminal() {
	fmt.Println("\nPress Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getAllSessions() ([]string, string) {
	sessionsDir := "sessions"
	files, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, ""
	}

	var sessions []string
	var latestSession string
	maxSession := 0

	for _, file := range files {
		if file.IsDir() {
			sessionName := file.Name()
			if num, err := strconv.Atoi(strings.TrimPrefix(sessionName, "SESSION_")); err == nil {
				if num > maxSession {
					maxSession = num
					latestSession = sessionName
				}
				sessions = append(sessions, sessionName)
			}
		}
	}
	return sessions, latestSession
}

func getSessionPath(sessionName string) string {
	return filepath.Join("sessions", sessionName)
}

func readTargets(filename string) ([]string, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var ids []string
	scanner := bufio.NewScanner(file)
	progress := 0

	if scanner.Scan() {
		firstLine := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(firstLine, "Progress:") {
			p := strings.TrimSpace(strings.TrimPrefix(firstLine, "Progress:"))
			if n, err := strconv.Atoi(p); err == nil {
				progress = n
			}
		} else if firstLine != "" {
			ids = append(ids, firstLine)
		}
	}

	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id != "" {
			ids = append(ids, id)
		}
	}

	return ids, progress, scanner.Err()
}

func updateProgress(filename string, progress int, ids []string) error {
	lines := []string{fmt.Sprintf("Progress: %d", progress)}
	lines = append(lines, ids...)
	return os.WriteFile(filename, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func showSplash() {
	fmt.Println(Blue + `
------------------------------------------------------------------------------------------------------------------------

███████╗████████╗███████╗ █████╗ ███╗   ███╗    ██╗██████╗      ██████╗██╗  ██╗███████╗ ██████╗██╗  ██╗███████╗██████╗ 
██╔════╝╚══██╔══╝██╔════╝██╔══██╗████╗ ████║    ██║██╔══██╗    ██╔════╝██║  ██║██╔════╝██╔════╝██║ ██╔╝██╔════╝██╔══██╗
███████╗   ██║   █████╗  ███████║██╔████╔██║    ██║██║  ██║    ██║     ███████║█████╗  ██║     █████╔╝ █████╗  ██████╔╝
╚════██║   ██║   ██╔══╝  ██╔══██║██║╚██╔╝██║    ██║██║  ██║    ██║     ██╔══██║██╔══╝  ██║     ██╔═██╗ ██╔══╝  ██╔══██╗
███████║   ██║   ███████╗██║  ██║██║ ╚═╝ ██║    ██║██████╔╝    ╚██████╗██║  ██║███████╗╚██████╗██║  ██╗███████╗██║  ██║
╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝    ╚═╝╚═════╝      ╚═════╝╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
                                                                                                                        
------------------------------------------------------------------------------------------------------------------------` + Cyan + `
STEAM ID AVAILABILITY CHECKER — by ytax - https://oguser.com/clarke

Send suggestions or report bugs at:` + Blue + ` https://github.com/ytax/steam-id-checker` + Cyan + `

This software will check for IDs inside "` + Blue + `targets.txt` + Cyan + `" feel free to replace the content of this file with a list of IDs
you want to check!

By default targets.txt is loaded with some shitty semi-og IDs, this tool is very good for finding 3c, 3l, 4c and 4l IDs.

I also recommend you to run the targets.txt file through a randomizer so you're not checking 
the usernames in alphabetic order` + Blue + `
------------------------------------------------------------------------------------------------------------------------` + Reset)
}

func main() {

	showSplash()

	sessions, latestSession := getAllSessions()

	fmt.Println(Cyan + "\n-> Existing sessions:" + Reset)
	for _, session := range sessions {
		if session == latestSession {
			fmt.Printf("  - %s"+Blue+" (LATEST SESSION)\n"+Reset, session)
		} else {
			fmt.Printf("  - %s\n", session)
		}
	}

	fmt.Println(Cyan + `
+-----------------------+
|` + Blue + ` 1. Start New Session` + Cyan + `  |
|` + Blue + ` 2. Resume Session` + Cyan + `     |
|` + Blue + ` 3. Exit` + Cyan + `               |
+-----------------------+` + Reset)
	fmt.Print(Cyan + "\n-> Choose an option" + Reset + ": ")

	var choice string
	fmt.Scanln(&choice)

	var sessionPath, targetsPath, outputPath string
	var isNewSession bool

	switch choice {
	case "1":
		newSessionName := "SESSION_" + strconv.Itoa(len(sessions)+1)
		sessionPath = getSessionPath(newSessionName)
		targetsPath = filepath.Join(sessionPath, "targets.txt")
		outputPath = filepath.Join(sessionPath, "output.txt")
		isNewSession = true
	case "2":
		fmt.Print(Cyan + "-> Enter the session name (e.g. SESSION_1): " + Reset)
		var chosenSession string
		fmt.Scanln(&chosenSession)
		sessionPath = getSessionPath(chosenSession)
		targetsPath = filepath.Join(sessionPath, "targets.txt")
		outputPath = filepath.Join(sessionPath, "output.txt")
		isNewSession = false
	case "3":
		fmt.Println(Red + "Exiting program. Hope you found some good IDs!" + Reset)
		os.Exit(0)
	default:
		fmt.Println(Red + "Invalid choice. Please restart the program." + Reset)
		return
	}

	if isNewSession {
		if err := os.MkdirAll(sessionPath, os.ModePerm); err != nil {
			fmt.Println(Red+"Error creating session directory (this is really weird, make sure controlled folder access isnt blocking the program or that you arent running from a place where the program doesnt have permission to write.):"+Reset, err)
			pauseTerminal()
			return
		}

		input, err := os.ReadFile("targets.txt")
		if err != nil {
			fmt.Println(Red+"Error reading targets.txt (this is really weird, make sure controlled folder access isnt blocking the program or that you arent running from a place where the program doesnt have permission to write.):"+Reset, err)
			pauseTerminal()
			return
		}
		if err := os.WriteFile(targetsPath, input, 0644); err != nil {
			fmt.Println(Red+"Error copying targets.txt (this is really weird, make sure controlled folder access isnt blocking the program or that you arent running from a place where the program doesnt have permission to write.):"+Reset, err)
			pauseTerminal()
			return
		}
	}

	ids, progress, err := readTargets(targetsPath)
	if err != nil {
		fmt.Println(Red+"Error reading targets file (this is really weird, make sure controlled folder access isnt blocking the program or that you arent running from a place where the program doesnt have permission to write.):"+Reset, err)
		pauseTerminal()
		return
	}

	if progress >= len(ids) {
		fmt.Println(Green + "All IDs have already been checked!" + Reset)
		pauseTerminal()
		return
	}

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(Red+"Error creating/opening output file (this is really weird, make sure controlled folder access isnt blocking the program or that you arent running from a place where the program doesnt have permission to write.):"+Reset, err)
		pauseTerminal()
		return
	}
	defer file.Close()

	fmt.Println(Cyan + "\nChecking Steam IDs...\n" + Reset)

	for i := progress; i < len(ids); i++ {
		id := ids[i]
		if !checkID(id) {
			fmt.Printf(Green+"Available: %s\n"+Reset, id)
			file.WriteString(id + "\n")
		} else {
			fmt.Printf(Red+"Not available: %s\n"+Reset, id)
		}

		if err := updateProgress(targetsPath, i+1, ids); err != nil {
			fmt.Println(Red+"Failed to update progress:"+Reset, err)
		}
	}

	fmt.Println(Green + "\nCheck completed. Available IDs saved to " + outputPath + Reset)
	pauseTerminal()
}
