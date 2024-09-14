# SyncApp

SyncApp is a command-line tool designed for syncing files and directories using
efficient compression methods, like zstd. It integrates with Git by adding hooks
to automatically sync files on push and pull operations. This tool simplifies
managing large files or assets in Git repositories without committing them
directly.

## Features

- **Fast Compression**: Uses zstd compression to create efficient backups of
  directories.
- **Git Integration**: Automatically syncs files before Git pushes and after
  checkouts using Git hooks.
- **Customizable**: Easily specify which files and directories to include in
  sync using patterns.
- **Cross-Platform**: Works on Linux, macOS, and Windows.
- **Keep Latest**: Option to keep only the most recent archive for each branch.

## Installation

1. **Clone the repository**:

```bash
git clone https://github.com/yourusername/syncapp.git
cd syncapp
```

2. **Build the application**:

You need to have Go installed. You can install Go from
[here](https://golang.org/doc/install).

```bash
go build -o syncapp
```

For windows:

```bash
go build -o syncapp.exe
```

3. **Move the binary** to your system's `PATH` for easy access (optional but
   recommended):

```bash
mv syncapp /usr/local/bin/
```

## Usage

### Initialize SyncApp

To set up SyncApp in your Git project, run the following command. This will
create a `sync.yaml` file and add Git hooks to manage syncing on push and
checkout.

```bash
syncapp init
```

- sync.yaml: This file contains your sync configuration, such as file patterns
  and directories to sync.

### Syncing Files

You can manually trigger a sync of your files by running the following command:

```bash
syncapp sync
```

This will create a compressed zstd archive of the specified files based on the
patterns in `sync.yaml`.

### Pulling Synced Files

When you switch branches or perform a Git checkout, SyncApp will automatically
pull synced files from the cloud directory. You can manually pull files using:

```bash
syncapp pull
```

### Pushing Synced Files

Before pushing changes to a Git repository, SyncApp will automatically push
synced files. To manually trigger this process:

```bash
syncapp push
```

### Configuring `sync.yaml`

The `sync.yaml` file defines what files and directories should be included in
the sync process. By default, it includes patterns for common assets like `.jpg`
and `.png` files. You can customize it to fit your project's needs. The
`keep_latest` option determines whether only the most recent archive for each
branch should be kept.

Example `sync.yaml`:

```yaml
# sync.yaml - SyncApp configuration file
cloud_dir: "/path/to/cloud"
patterns:
  - "*.jpg"
  - "*.png"
  - "assets/*"
keep_latest: true
```

- `cloud_dir`: Specifies where the compressed files will be stored.
- `patterns`: List of file patterns to include in the sync. Use `*` as a
  wildcard.
- `keep_latest`: If set to `true`, only the most recent archive for each branch
  will be kept. If `false`, older archives will not be removed automatically.

## Git Integration

SyncApp automatically manages Git hooks:

- **Pre-push hook**: Runs `syncapp push` before pushing changes to the remote
  repository.
- **Post-checkout hook**: Runs `syncapp pull` after a Git checkout. These hooks
  are set up when you run syncapp init.

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository.
2. Create a new feature branch (git checkout -b feature/your-feature).
3. Commit your changes (git commit -am 'Add new feature').
4. Push to the branch (git push origin feature/your-feature).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for
details.
