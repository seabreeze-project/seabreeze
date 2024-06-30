# Seabreeze

Seabreeze is a really simple container orchestration tool with superpowers, optimized for the management of web applications.

It serves as a "poor man's Kubernetes" while also providing a comprehensive set of utilities for containerized environments.

**Please note that Seabreeze is currently under heavy development. We are open to and deeply thankful for constructive feedback, innovative ideas, and any kind of code contributions.** Your input can help shape the future of this project and make it more useful for the community. Please join the discussions in our [issue tracker](https://github.com/seabreeze-project/seabreeze/issues).

## Features

- 🚀 **Project Management:** Easily create and manage multiple compose projects.
- 📜 **Script Runner:** Run predefined commands within containers and on the host.
- 🐚 **Shell Mode (WIP):** Run Seabreeze commands in an interactive shell-like interface.

### Planned Features

- 🔀 **VCS Integration:** Clone and update source code that is used for service containers.
- ⏰ **Tasks (Cronjobs):** Automatically run commands in containers based on a defined schedule.
- 🔐 **Docker API Proxy:** Control and restrict access to Docker endpoints.
- 📚 **Recipes:** Automatically set up projects using predefined workflows.

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

## History

Seabreeze is a project that has evolved organically over time. It began as a collection of scripts I developed to simplify the management of containerized applications. These scripts were eventually rewritten in Go, and refined to work together seamlessly. Over this period, I accumulated several unimplemented ideas that aligned perfectly with the scope of Seabreeze. These ideas are now finding a home within this ecosystem, with the potential to be implemented in the future.

The name "Seabreeze" itself has a bit of history too. It was inspired by an older, somewhat similar tool I had created. This project name carries forward the essence of that earlier tool, now refined and expanded into a more powerful and versatile solution.

— [secondtruth](https://github.com/secondtruth)

## Contributing

We welcome contributions! Feel free to fork the repository and submit a pull request for review.

We appreciate all kinds of contributions, whether they are new code, bug reports and fixes, documentation improvements, or fresh ideas. We encourage you to open an issue to discuss your ideas or to report any bugs you encounter.

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgements

- Docker API Proxy is based on [docker-proxy-acl](https://github.com/titpetric/docker-proxy-acl)
- Shell Mode is made possible by [go-prompt](https://github.com/c-bata/go-prompt) and [go-shellwords](https://github.com/mattn/go-shellwords)
