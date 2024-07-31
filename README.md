# File-Integrity-Monitor
A lightweight file integrity monitoring application written in Go.

# How-To/Workflow
- Build application (go build -o FIM)
- Initially run the program and specify a root directory you want to monitor as well as the baseline file name. The application will scan all files and subdirectories in the root
- When you want to run a scan to check the integrity of all files in the original root directory, run the application with the "*-compare*" flag and specify the baseline file of the directory you want to scan. The application will check the root directory of the JSON file and automatically re-scan it and determine which files have been changed based on hash value, which files have been added, and also removed anywhere within the root directory. Output will be both printed to the console and written to a new json file.

# Flags
### Creating baseline file:
These two flags need to be ran together to generate the baseline file:
- **-dir \<*target path*\>:** Specify the directory you want to scan to build a baseline on
- **-baseline \<*outfile name*\>:** Specify the name of the baseline file. *_baseline\<current date\>* will be appended to the file name.

### Re-scan and see changes to baseline file:
- **-compare \<*baseline file name*\>:** Specify the baseline file you want to compare against.
