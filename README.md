# RDS Migrator

A lightweight Go tool for migrating Redis data between different Redis instances, with support for Dragonfly and Redis databases.

## 🚀 Features

- **Cross-Platform Migration**: Migrate data from Dragonfly to Redis or between Redis instances
- **TTL Preservation**: Maintains original key expiration times during migration
- **Bulk Migration**: Efficiently migrates all keys using Redis SCAN operations
- **Error Handling**: Graceful error handling with detailed logging
- **Lightweight**: Single binary with minimal dependencies

## 📁 Project Structure

```
rds-migrator/
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── Makefile            # Build automation
├── migrator.go         # Main migration logic
├── README.md           # This documentation
├── env.example         # Environment variables template
└── config.yaml.example # YAML configuration template
```

## 📋 Prerequisites

- Go 1.23.0 or higher
- Access to source and target Redis/Dragonfly instances
- Network connectivity between source and target servers

## 🛠️ Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/atadzan/rds-migrator.git
cd rds-migrator
```

2. Install dependencies:
```bash
go mod download
```

3. Build the binary:
```bash
go build -o migrator ./migrator.go
```

### Using Makefile

For Linux AMD64 builds:
```bash
make build
```

## ⚙️ Configuration

The migrator supports multiple configuration methods. Choose one of the following approaches:

### Method 1: Environment Variables (Recommended)

Set the following environment variables before running the migrator:

```bash
# Source Redis Configuration
export SOURCE_ADDR="192.168.1.137:6379"
export SOURCE_PASSWORD="your_source_password"
export SOURCE_DB="2"

# Target Redis Configuration
export TARGET_ADDR="192.168.1.226:6379"
export TARGET_PASSWORD="your_target_password"
export TARGET_DB="2"
```

Or create a `.env` file (copy from `env.example`):

```bash
cp env.example .env
# Edit .env with your actual connection details
```

### Method 2: YAML Configuration File

Copy the example configuration file:

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your actual connection details
```

### Configuration Parameters

| Parameter | Description | Required | Default |
|-----------|-------------|----------|---------|
| `SOURCE_ADDR` | Source Redis host:port | Yes | - |
| `SOURCE_PASSWORD` | Source Redis password | Yes | - |
| `SOURCE_DB` | Source Redis database number | No | 0 |
| `TARGET_ADDR` | Target Redis host:port | Yes | - |
| `TARGET_PASSWORD` | Target Redis password | Yes | - |
| `TARGET_DB` | Target Redis database number | No | 0 |

## 🚀 Usage

### Loading Environment Variables

If you're using a `.env` file, you can load it using a tool like `dotenv`:

```bash
# Install dotenv if you don't have it
go install github.com/joho/godotenv/cmd/godotenv@latest

# Run with dotenv
dotenv ./migrator
```

Or source the environment variables manually:

```bash
source .env
./migrator
```

### Basic Migration

```bash
./migrator
```

The tool will:
1. Connect to the source Redis/Dragonfly instance
2. Scan all available keys
3. Migrate each key with its value and TTL to the target instance
4. Display progress and any errors encountered

### Migration Process

1. **Connection**: Establishes connections to both source and target instances
2. **Scanning**: Uses Redis SCAN to iterate through all keys (batch size: 500)
3. **Data Extraction**: Dumps each key's value and TTL from source
4. **Data Transfer**: Restores each key to target with original TTL
5. **Completion**: Reports migration status

## 📊 Supported Operations

- **Key Migration**: All key types (strings, hashes, lists, sets, sorted sets)
- **TTL Preservation**: Maintains original expiration times
- **Batch Processing**: Efficient handling of large datasets
- **Error Recovery**: Continues migration on individual key failures

## 🔧 Customization

### Batch Size

Modify the scan batch size in the code:
```go
iter := src.Scan(ctx, 0, "*", 500).Iterator() // Change 500 to desired batch size
```

### Key Pattern Filtering

To migrate only specific keys, modify the scan pattern:
```go
iter := src.Scan(ctx, 0, "user:*", 500).Iterator() // Only migrate keys starting with "user:"
```

## 📝 Logging

The tool provides real-time feedback:
- Migration progress
- Individual key errors
- Overall completion status

## ⚠️ Important Notes

- **Data Overwrite**: Uses `RestoreReplace` which will overwrite existing keys
- **Network Stability**: Ensure stable network connection during migration
- **Memory Usage**: Large datasets may require sufficient memory
- **Downtime**: Consider running during low-traffic periods

## 🐛 Troubleshooting

### Common Issues

1. **Connection Failed**
   - Verify host, port, and password
   - Check network connectivity
   - Ensure Redis is running and accessible

2. **Authentication Error**
   - Verify password and username (if applicable)
   - Check Redis ACL settings

3. **Migration Errors**
   - Check available memory on target
   - Verify target Redis version compatibility
   - Review error logs for specific key issues


## 🙏 Acknowledgments

- Built with [go-redis](https://github.com/redis/go-redis) library
- Supports both Redis and Dragonfly databases

---

**Note**: This tool is designed for data migration purposes. Always backup your data before running migrations in production environments.
