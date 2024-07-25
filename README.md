# File-Integrity-Monitor
A lightweight file integrity monitoring application written in Go.

# Flags
- **-dir \<*target path*\>:** Specify the directory you want to scan to build a baseline on
- **-baseline \<*outfile name*\>:** Specify the name of the baseline file. *_baseline\<current date\>* will be apended to the file name.
- **-compare \<*baseline file name*\>:** Specify the baseline file you want to compare against.

# How-To/Workflow
- Build application (go build -o FIM)
- Initially run the program and specify a root directory you want to monitor as well as the baseline file name. The application will scan all files and subdirectories in the root
- When you want to run a scan to check the integrity of the file, run the application with the "*-compare*" flag and specify the baseline file of the directory you want to scan. The application will check the root directory of the JSON file and automatically scan it. Output will be both printed to the console and written to a new json file.
