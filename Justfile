# Frostyard Site Justfile
# MkDocs documentation site with git submodules

set dotenv-load := true

python := "/home/linuxbrew/.linuxbrew/bin/python3"
venv_dir := ".venv"
venv_python := venv_dir / "bin/python"
venv_pip := venv_dir / "bin/pip"

default:
    just --list --unsorted

# === Virtual Environment ===

# Create the virtual environment
venv-create:
    {{ python }} -m venv {{ venv_dir }}
    @echo "Virtual environment created. Run 'source {{ venv_dir }}/bin/activate' to activate."

# Install dependencies into the virtual environment
venv-install:
    {{ venv_pip }} install -r requirements.txt

# Recreate venv and install dependencies
venv-setup: venv-create venv-install

# Remove the virtual environment
venv-clean:
    rm -rf {{ venv_dir }}

# === MkDocs Site ===

# Serve the site locally with live reload
serve:
    {{ venv_python }} -m mkdocs serve

# Serve on all interfaces (useful for containers/VMs)
serve-public:
    {{ venv_python }} -m mkdocs serve -a 0.0.0.0:8000

# Build the static site
build:
    {{ venv_python }} -m mkdocs build

# Build with strict mode (treat warnings as errors)
build-strict:
    {{ venv_python }} -m mkdocs build --strict

# Clean the build output
build-clean:
    rm -rf site

# === Git Submodules ===

# Initialize all submodules
submodules-init:
    git submodule update --init --recursive

# Update all submodules to latest commit on their tracked branch
submodules-update:
    git submodule update --remote --merge

# Pull latest changes for all submodules
submodules-pull:
    git submodule foreach git pull origin main

# Show status of all submodules
submodules-status:
    git submodule status

# Sync submodule URLs from .gitmodules
submodules-sync:
    git submodule sync --recursive

# Deinitialize all submodules (removes working tree)
submodules-deinit:
    git submodule deinit --all -f

# === Combined Targets ===

# Full setup: init submodules, create venv, install deps
setup: submodules-init venv-setup

# Clean everything
clean: build-clean venv-clean

publish: build
    gh-pages -d site -m "Publish site"