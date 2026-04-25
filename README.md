# ACMP - Assetto Corsa Mod Manager

A Windows desktop application for managing mods and profiles for Assetto Corsa. Built with Wails, Go, and Svelte 5.

## Features

- **Mod Management**: Scan and manage mods from a designated folder
- **Profile System**: Create, edit, and organize mod profiles
- **Symlink Activation**: Activate/deactivate mods using symbolic links
- **Automatic Profiles**: Auto-generate profiles based on folder structure
- **SQLite Storage**: Local database for mods, profiles, and settings

## Requirements

- Go 1.26+
- Node.js 20.19+ or 22.12+
- Wails CLI v2.11.0+
- **Administrator privileges** (required for symlink operations)

## Important: Running as Administrator

**Always run the application as Administrator** to ensure proper access to create and delete symbolic links. Without administrator privileges, the mod activation/deactivation feature will not work.

## Build Instructions

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install

# Build production binary
cd ..
wails build
```

The built executable will be in the `bin` directory.

## Usage Instructions

1. **Run as Administrator** - Right-click the executable and select "Run as administrator"
2. Configure paths in Settings:
   - **Mods Path**: Select your mods folder (containing tracks/cars subfolders)
   - **Assetto Corsa Path**: Select your Assetto Corsa installation directory
3. Enable Automatic Profiles if you want auto-generated profiles based on folder structure
4. Create profiles and add mods to them
5. Activate/deactivate profiles to enable/disable mods via symlinks

## Project Structure

```
├── frontend/               # Svelte 5 frontend
│   ├── src/
│   │   ├── components/    # UI components
│   │   ├── pages/         # Page components
│   │   └── lib/           # UI library components
│   └── package.json
├── app.go                 # Main application logic
├── main.go                # Application entry point
├── settings.go            # Settings management
├── database/              # Database operations
│   ├── database.go        # SQLite operations
│   ├── scanner.go         # Mod filesystem scanner
│   └── auto_profiles.go   # Auto profile generation
├── models/                # Data models
└── symlink/               # Symlink management
```

## Tech Stack

- **Backend**: Go with Wails v2
- **Frontend**: Svelte 5 with TypeScript
- **Styling**: Tailwind CSS v4
- **UI Components**: shadcn-svelte (bits-ui)
- **Build Tool**: Vite 8
- **Database**: SQLite (modernc.org/sqlite)

## License

MIT License