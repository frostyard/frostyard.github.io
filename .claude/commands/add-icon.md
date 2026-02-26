Add a new Heroicon to `templates/components/icon.templ` for use in sidebar section titles.

## Arguments

$ARGUMENTS should be a Heroicon name (e.g. "book-open", "globe-alt", "photo").

## Steps

1. Fetch the 16x16 solid SVG from https://raw.githubusercontent.com/tailwindlabs/heroicons/master/src/16/solid/$ARGUMENTS.svg
2. Add a new `case` to the switch statement in `templates/components/icon.templ`
3. In the SVG, replace `fill="black"` with `fill="currentColor"` and set `class="w-4 h-4 shrink-0"`
4. Remove the `width="16" height="16"` attributes (the class handles sizing)
5. Run `templ generate`
6. Run `go test ./... -v`
7. Show the user where to use it: set `icon: "$ARGUMENTS"` in a content `_index.md` frontmatter

## Pattern

Each case in the switch looks like:

```templ
case "icon-name":
    <svg class="w-4 h-4 shrink-0" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
        <path ...></path>
    </svg>
```

Add the new case before the closing `}` of the switch block.
