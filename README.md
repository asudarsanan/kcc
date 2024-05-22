# KCC (Kube Context Changer)

KCC is a CLI tool that helps you easily switch between Kubernetes cluster contexts in your `kubeconfig` file. It provides an interactive selector interface for choosing the desired context.

## Features

- List all available Kubernetes contexts.
- Interactive CLI for selecting and switching contexts.
- Lightweight and easy to use.

## Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/asudarsanan/kcc.git
    cd kcc
    ```

2. **Build the project:**

   Ensure you have [Go](https://golang.org/dl/) installed on your machine. Then run:

    ```sh
    go build -o kcc
    ```

3. **Move the executable to a directory in your PATH:**

    ```sh
    # On macOS and Linux:
    sudo mv kcc /usr/local/bin/
   # or
   `sudo install kcc /usr/local/bin/kcc`

    # On Windows:
    move kcc.exe C:\path\to\directory\in\PATH\
    ```

4. **Ensure the directory is in your PATH:**

   ### macOS and Linux

   Add the following line to your `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc` file:

    ```sh
    export PATH=$PATH:/usr/local/bin
    ```

   Then source the file:

    ```sh
    source ~/.bashrc   # or ~/.bash_profile or ~/.zshrc
    ```

   ### Windows

    1. Open the Start Search, type in "env", and choose "Edit the system environment variables".
    2. In the System Properties window, click on the "Environment Variables" button.
    3. In the Environment Variables window, select the `Path` variable in the "System variables" section and click "Edit".
    4. Click "New" and add the path to the directory where you moved `kcc.exe`, e.g., `C:\path\to\directory\in\PATH`.
    5. Click "OK" to close all windows.

## Usage

After building the project and adding it to your `PATH`, you can run the `kcc` executable to interactively select and switch between your Kubernetes contexts.

```sh
kcc
```
## Example

### Running `kcc`:

```sh
Select Kubernetes cluster context:
  ❯ context-1
    context-2
    context-3
```
Select the context you want to switch to, and press Enter.
The tool will switch to the selected context and output:

```shell
Switched to context context-1
```
## Development
To contribute to this project, follow these steps:

Fork the repository:
Click the `"Fork"` button at the top right corner of this repository's GitHub page.

Clone your fork:
```shell
git clone https://github.com/asudarsanan/kcc.git
cd kcc
```
Create a new branch:
```shell
git checkout -b feature-branch-name
```
Make your changes and commit them:
```shell
git add .
git commit -m "Description of changes"
```
Push to your fork and create a pull request:
```shell
git push origin feature-branch-name
```

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements
- [Promptui](https://github.com/manifoldco/promptui) for the interactive CLI prompt.
- [yaml.v2](https://pkg.go.dev/gopkg.in/yaml.v2) for YAML parsing.

## Contact
If you have any questions, feel free to open an issue or contact the project maintainer at [X](https://twitter.com/ReallyG8Site).