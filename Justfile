# Frostyard Site
# Go static site generator with Templ + Tailwind

default:
    just --list --unsorted

# Generate templ files and build the static site
build:
    templ generate
    go run ./cmd/frostyard build

# Start dev server with live reload
serve:
    templ generate
    go run ./cmd/frostyard serve

# Run all tests
test:
    go test ./... -v

# Generate templ Go code
generate:
    templ generate

# Create a new docs page
new-page path:
    go run ./cmd/frostyard new page {{ path }}

# Create a new blog post
new-post title:
    go run ./cmd/frostyard new post "{{ title }}"

# Run Pagefind to build search index (post-build)
search-index:
    pagefind --site dist

# Clean build artifacts
clean:
    rm -rf dist
