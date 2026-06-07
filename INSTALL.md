# Installing language toolchains for DevAscent

DevAscent ships with **no bundled language runtimes** — it stays small and platform-independent, and you install only the language(s) you want to play. When you pick a language, DevAscent checks whether its toolchain is installed and actually works; if not, it points you here.

You can read every primer and lesson without installing anything — a toolchain is only needed to *run and grade* your own code in that language.

<!-- Generated from internal/content/data/install/*.yaml — edit those, not this file. -->

## Python

> DevAscent needs Python 3.8 or newer. After installing, restart DevAscent (or press [R] to re-check).

### Windows

Download: https://www.python.org/downloads/windows/

1. Download the latest Python 3 installer from python.org.
2. Run it and CHECK "Add python.exe to PATH" on the first screen.
3. Finish the install, then restart DevAscent.

Verify: `python --version`

### macOS

Download: https://www.python.org/downloads/macos/

1. Install via Homebrew - "brew install python" - or download the macOS installer from python.org.
2. Restart your terminal so PATH updates.

Verify: `python3 --version`

### Linux

Download: https://www.python.org/downloads/source/

1. Install with your package manager, e.g. "sudo apt install python3" (Debian/Ubuntu) or "sudo dnf install python3" (Fedora).

Verify: `python3 --version`

## JavaScript

> JavaScript runs on Node.js. Installing Node also gives you npm (needed for TypeScript).

### Windows

Download: https://nodejs.org/en/download

1. Download the Node.js LTS Windows installer from nodejs.org.
2. Run it with the default options (it adds Node to PATH).
3. Restart DevAscent.

Verify: `node --version`

### macOS

Download: https://nodejs.org/en/download

1. Install via Homebrew - "brew install node" - or download the macOS LTS installer from nodejs.org.
2. Restart your terminal.

Verify: `node --version`

### Linux

Download: https://nodejs.org/en/download/package-manager

1. Install via your package manager or nvm. With nvm - "nvm install --lts".
2. Or - "sudo apt install nodejs npm" (Debian/Ubuntu).

Verify: `node --version`

## TypeScript

> TypeScript needs Node.js first (for npm and to run the compiled output). Then install the TypeScript compiler (tsc) globally.

### Windows

Download: https://www.typescriptlang.org/download

1. Install Node.js first (see the JavaScript guide).
2. Open a terminal and run - "npm install -g typescript".
3. Restart DevAscent.

Verify: `tsc --version`

### macOS

Download: https://www.typescriptlang.org/download

1. Install Node.js first (see the JavaScript guide).
2. Run - "npm install -g typescript".
3. Restart your terminal.

Verify: `tsc --version`

### Linux

Download: https://www.typescriptlang.org/download

1. Install Node.js first (see the JavaScript guide).
2. Run - "npm install -g typescript" (you may need sudo, or configure a user npm prefix).

Verify: `tsc --version`

## Go

> The official Go distribution includes everything you need.

### Windows

Download: https://go.dev/dl/

1. Download the Windows MSI installer from go.dev/dl.
2. Run it (it adds Go to PATH automatically).
3. Restart DevAscent.

Verify: `go version`

### macOS

Download: https://go.dev/dl/

1. Download the macOS package from go.dev/dl, or "brew install go".
2. Restart your terminal.

Verify: `go version`

### Linux

Download: https://go.dev/doc/install

1. Download the Linux tarball and extract it to /usr/local per go.dev/doc/install.
2. Add /usr/local/go/bin to your PATH.

Verify: `go version`

## Java

> You need a full JDK (which includes the javac compiler), NOT just a JRE. Eclipse Temurin is a free, well-supported build.

### Windows

Download: https://adoptium.net/temurin/releases/

1. Download the latest Temurin JDK .msi for Windows from adoptium.net.
2. Run the installer and enable "Add to PATH" / "Set JAVA_HOME".
3. Restart DevAscent.

Verify: `javac -version`

### macOS

Download: https://adoptium.net/temurin/releases/

1. Install the Temurin JDK .pkg from adoptium.net, or "brew install --cask temurin".
2. Restart your terminal.

Verify: `javac -version`

### Linux

Download: https://adoptium.net/installation/

1. Install a JDK via your package manager, e.g. "sudo apt install default-jdk" (Debian/Ubuntu), or use the Adoptium APT/YUM repo.

Verify: `javac -version`

## C#

> You need the .NET SDK (which includes the compiler), NOT just the .NET Runtime. "dotnet --list-sdks" must show at least one entry.

### Windows

Download: https://dotnet.microsoft.com/download

1. Download the .NET SDK (not the Runtime) installer from dotnet.microsoft.com/download.
2. Run it, then restart DevAscent.

Verify: `dotnet --list-sdks`

### macOS

Download: https://dotnet.microsoft.com/download

1. Download the .NET SDK .pkg from dotnet.microsoft.com/download, or "brew install --cask dotnet-sdk".
2. Restart your terminal.

Verify: `dotnet --list-sdks`

### Linux

Download: https://learn.microsoft.com/dotnet/core/install/linux

1. Install the .NET SDK via your distro's package feed (see Microsoft's Linux install docs).

Verify: `dotnet --list-sdks`

## C++

> C++ is currently REFERENCE-ONLY in DevAscent — you can read its primers and topics without installing anything, and graded play isn't supported yet. Installing a compiler (below) is optional; it only lets you compile the C++ examples yourself.

### Windows

Download: https://www.msys2.org/

1. Install MSYS2 from msys2.org, then in its terminal run - "pacman -S mingw-w64-ucrt-x86_64-gcc".
2. Add the MSYS2 ucrt64\bin folder to your PATH (it contains g++).
3. Restart DevAscent. (Alternatively, install "Visual Studio Build Tools" with the C++ workload.)

Verify: `g++ --version`

### macOS

Download: https://developer.apple.com/xcode/

1. Run "xcode-select --install" to get Apple's clang/clang++ command-line tools.
2. Restart your terminal.

Verify: `clang++ --version`

### Linux

Download: https://gcc.gnu.org/

1. Install GCC/G++, e.g. "sudo apt install build-essential" (Debian/Ubuntu) or "sudo dnf install gcc-c++" (Fedora).

Verify: `g++ --version`

## Rust

> rustup installs the Rust compiler (rustc) and Cargo, and configures a default toolchain.

### Windows

Download: https://rustup.rs/

1. Download and run rustup-init.exe from rustup.rs.
2. Choose the default install (option 1). It may prompt you to install the MSVC build tools - accept.
3. Restart DevAscent.

Verify: `rustc --version`

### macOS

Download: https://rustup.rs/

1. Run the rustup command from rustup.rs in your terminal.
2. Run "source $HOME/.cargo/env" or restart your terminal.

Verify: `rustc --version`

### Linux

Download: https://rustup.rs/

1. Run the rustup command from rustup.rs in your terminal.
2. Run "source $HOME/.cargo/env" or restart your terminal.

Verify: `rustc --version`

