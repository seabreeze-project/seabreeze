# Seabreeze

Seabreeze is a really simple container orchestration tool with superpowers, optimized for the management of web applications.

It serves as a "poor man's Kubernetes" while also providing a comprehensive set of utilities for containerized environments.

## Features

- **Project Management:** Easily create and manage multiple compose projects.
- **Script Runner:** Run predefined commands within containers and on the host.
- **Shell Mode (WIP):** Run Seabreeze commands in an interactive shell-like interface.

### Planned Features

- **VCS Integration:** Clone and update source code that is used for service containers.
- **Cron System:** Schedule commands to run automatically in containers based on a defined schedule.
- **Docker API Proxy:** Control and restrict access to Docker endpoints.
- **Recipes:** Automatically set up projects using predefined workflows.

## Build

### Linux / macOS

To build the binary on Linux or macOS, run:

```bash
go build -o bin/seabreeze
```

### Windows

To build the binary on Windows, run:

```powershell
go build -o bin\seabreeze.exe
```

## Contributing

We welcome and appreciate contributions! Please fork the repository and submit a pull request for review.

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgements

- Docker API Proxy is based on [docker-proxy-acl](https://github.com/titpetric/docker-proxy-acl)
- Shell Mode is made possible by [go-prompt](https://github.com/c-bata/go-prompt) and [go-shellwords](https://github.com/mattn/go-shellwords)
