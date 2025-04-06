# Prayer Times Cli

Cli tool to show today's prayer times

![img](docs/prayers-command-preview.png)

Features
- Show prayer times to current day (or provide a specific day if you want)
- Show time left till next prayer


## Installation
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/MABD-dev/prayer-times-cli
    cd prayer-times-cli
    ```
2.  **Ensure Go is installed:**
    * Make sure Go `1.24.1` or later is installed on your machine.
    * You can check using:
        ```bash
        go version
        ``` 
3.  **Build the project:**
    ```bash
    go build
    ```
    This will generate the `prayer-times-cli` executable.

    > Optional: I like to shorten the executabel file name to `prayers`. you can do that by
    ```sh
    mv prayer-times-cli prayers
    ```

4.  **Add to your PATH (Linux/macOS):**
    * Move the `prayer-times-cli` executable to a directory in your `PATH` (e.g., `~/bin` or `/usr/local/bin`).
    * Alternatively, add the directory containing the executable to your `PATH` environment variable.


## Usage
NOTE: use `prayer-times-cli` if you did not rename it to `prayers`

```sh
prayers # will show prayer times of the day
```
```sh
prayers -y 2025 -m 5 -d 11 

# Or
prayers --year 2025 --month 5 --day 11 
```
By default year, month and day are today's dates, but you can override any of them to values you like. 
> NOTE: datas in future years might not work


## Roadmap
Check [issues](https://github.com/MABD-dev/prayer-times-cli/issues)

<br/><br/><br/>

# Special Thanks

This project was possible because of this [ibad-al-rahman/prayer-times](https://github.com/ibad-al-rahman/prayer-times) repo. That provided an easy way to generate and fetch prayer times in simple api request

