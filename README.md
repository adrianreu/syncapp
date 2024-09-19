<!-- Center the SVG logo -->
<div align="center">
  <img src="logo.svg" alt="Syncwave Logo" width="25%"/>
  <h1>Syncwave</h1>
</div>

Syncwave is a command-line tool designed for syncing files and directories using
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
- **`.gitignore` Syncing**: Automatically updates `.gitignore` to exclude files
  that are synced, preventing large assets from being accidentally committed.
- **Cloud Storage**: Use Syncwave with cloud services like Google Drive, MEGA, or
  your own network drive for syncing archives. Note: Access management is handled by
  the respective cloud provider.


## Installation

1. **Clone the repository**:

```bash
git clone https://github.com/yourusername/syncwave.git
cd syncwave
```

2. **Build the application**:

You need to have Go installed. You can install Go from
[here](https://golang.org/doc/install).

```bash
go build -o syncwave
```

For windows:

```bash
go build -o syncwave.exe
```

3. **Move the binary** to your system's `PATH` for easy access (optional but
   recommended):

```bash
mv syncwave /usr/local/bin/
```

## Usage

### Initialize Syncwave

To set up Syncwave in your Git project, run the following command. This will
create a `sync.yaml` file and add Git hooks to manage syncing on push and
checkout.

```bash
syncwave init
```

- `sync.yaml`: This file contains your sync configuration, such as file patterns
  and directories to sync.

### Syncing Files

You can manually trigger a sync of your files by running the following command:

```bash
syncwave sync
```

This will create a compressed zstd archive of the specified files based on the patterns in `sync.yaml`. Files that match the patterns and any large files (if you've set a max file size) will automatically be added to both the `.gitignore` and the archive, ensuring they are not committed to your Git repository by accident.


### Pulling Synced Files

When you switch branches or perform a Git checkout, Syncwave will automatically
pull synced files from the cloud directory. You can manually pull files using:

```bash
syncwave pull
```

### Pushing Synced Files

Before pushing changes to a Git repository, Syncwave will automatically push
synced files. To manually trigger this process:

```bash
syncwave push
```

### Configuring `sync.yaml`

The `sync.yaml` file defines what files and directories should be included in
the sync process. By default, it includes patterns for common assets like `.jpg`
and `.png` files. You can customize it to fit your project's needs. The
`keep_latest` option determines whether only the most recent archive for each
branch should be kept.

Additionally, files that match your patterns or exceed the specified size limit are automatically added to `.gitignore`, preventing large assets from being committed to your repository.

Example `sync.yaml`:

```yaml
# sync.yaml - Syncwave configuration file
cloud_dir: "/path/to/cloud"
patterns:
  - "*.jpg"
  - "*.png"
  - "assets/*"
keep_latest: true
max_file_size: 10485760  # Maximum file size in bytes (e.g., 10 MB)
```

- `cloud_dir`: Specifies where the compressed files will be stored.
- `patterns`: List of file patterns to include in the sync. Use `*` as a
  wildcard.
- `keep_latest`: If set to `true`, only the most recent archive for each branch
  will be kept. If `false`, older archives will not be removed automatically.
- `max_file_size` (NOT READY YET): Files larger than this size (in bytes) will automatically be added to .gitignore and included in the archive. This option is optional and can be left undefined.

## Cloud Storage Integration

You can specify a `cloud_dir` that points to cloud-synced folders like those from Google Drive, MEGA, or your own network-attached storage (NAS). Note that Syncwave itself doesn’t handle authentication or file transfer to cloud providers — you’ll need to manage that separately via the respective cloud service's app or sync mechanism.

For example, you can point `cloud_dir` to a Google Drive or MEGA-synced folder on your local system:
```yaml
cloud_dir: "/path/to/Google Drive/syncwave"
```
This will ensure that archives are stored in your cloud provider's folder, and they will be synced automatically by the cloud service’s app.

## Git Integration

Syncwave automatically manages Git hooks:

- **Pre-push hook**: Runs `syncwave push` before pushing changes to the remote
  repository.
- **Post-checkout hook**: Runs `syncwave pull` after a Git checkout. These hooks
  are set up when you run ``syncwave`` init.

## Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository.
2. Create a new feature branch (``git checkout -b feature/your-feature``).
3. Commit your changes (``git commit -am 'Add new feature'``).
4. Push to the branch (``git push origin feature/your-feature``).
5. Create a new Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for
details.
