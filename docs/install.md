# Install `pbauthor`

## Install the command line tool

Follow the instructions below for your operating system to install the
`pbauthor` command line tool. Be sure to follow the
[Verify your installation](#verify-your-installation) steps after following the
OS-specific instructions.

### macOS

Use [Homebrew](https://brew.sh/).

1.  [Install Homebrew](https://brew.sh/) (if not already installed).
2.  Install the command line tool.
    ```shell
    brew install pathbird/tap/pbauthor
    ```

### Windows

1.  [Install Scoop](https://scoop.sh/) (if not already installed). From
    PowerShell (**not** _Command Prompt_), run these commands to install Scoop.

    ```powershell
    # Allow PowerShell to execute arbitrary code (you may be prompted to allow this)
    Set-ExecutionPolicy RemoteSigned -scope CurrentUser

    # Download and execute the Scoop installer
    Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://get.scoop.sh')
    ```

2.  Install the command line tool.

    ```powershell
    # Install git (this is required to add the Pathbird bucket)
    scoop install git

    # Tell scoop where to find Pathbird's published packages
    scoop bucket add pathbird https://github.com/pathbird/scoop-bucket.git

    # Install the pbauthor tool
    scoop install pathbird/pbauthor
    ```

### Custom installation

If you're using Linux, or you want more control over your installation, you can
follow these steps.

1.  Download the binary from the
    [latest GitHub release](https://github.com/pathbird/pbauthor/releases).
2.  Add the executable to your `PATH`. Make sure to set the executable bit if
    necessary.

Alternatively, you can download and build the Go source code.

## Verify your installation

1.  Log in to Pathbird using the command line tool.
    ```shell
    # Authenticate with the Pathbird API
    # This should prompt you for your email and password
    pbauthor auth login
    ```
2.  Verify that the authentication worked.
    ```shell
    pbauthor auth status
    ```
    The output should look something like this:
    ```
    ✅ Authenticated (until Tue, 06 Apr 2021 00:38:46 UTC)
    ```

# Upgrade

The `pbauthor` command line is periodically updated. It may be necessary
to update your installed version in order to deal with breaking changes in the
Pathbird API or to access newer features.

## Upgrade on macOS (Homebrew)

Open a terminal and run this command.

```shell
brew upgrade pbauthor
```

## Upgrade on Windows (Scoop)

Open PowerShell and run this command.

```powershell
scoop update pbauthor
```

## Upgrade custom installation

Download a new binary as described in the
[Custom installation](#custom-installation) and replace the old binary with the
new binary.
