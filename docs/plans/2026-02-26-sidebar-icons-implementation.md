# Sidebar Icons Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Render Heroicon SVGs next to section titles in the sidebar navigation.

**Architecture:** Add `Icon string` field to the sidebar component data, pass it through from content sections, and render via a new Templ icon component that maps icon names to inline SVG markup.

**Tech Stack:** Go, Templ, Tailwind CSS, Heroicons (16x16 solid)

---

### Task 1: Add Icon field to SidebarSection and pass it through

**Files:**
- Modify: `templates/components/sidebar.templ:3-8` (SidebarSection struct)
- Modify: `internal/render/render.go:90-94` (convertSection function)

**Step 1: Add Icon field to SidebarSection struct**

In `templates/components/sidebar.templ`, add `Icon` to the struct:

```go
type SidebarSection struct {
	Title       string
	Path        string
	Icon        string
	Pages       []SidebarLink
	Subsections []SidebarSection
}
```

**Step 2: Pass Icon through in convertSection**

In `internal/render/render.go`, update the `convertSection` function at line 91-94:

```go
ss := components.SidebarSection{
    Title: sec.Title,
    Path:  sec.Path,
    Icon:  sec.Icon,
}
```

**Step 3: Run tests to verify nothing breaks**

Run: `go test ./... -v`
Expected: All existing tests pass (the new field is just a string, no behavior change yet).

**Step 4: Commit**

```
git add templates/components/sidebar.templ internal/render/render.go
git commit -m "feat: pass icon field through to sidebar component"
```

---

### Task 2: Create the Icon templ component

**Files:**
- Create: `templates/components/icon.templ`

**Step 1: Create icon.templ with Heroicon SVGs**

Create `templates/components/icon.templ` with the `Icon` component. The SVGs use `fill="currentColor"` so they inherit text color from their parent. Size is `w-4 h-4` (16x16) with `shrink-0` to prevent flex shrinking.

```templ
package components

templ Icon(name string) {
	if name == "server" {
		<svg class="w-4 h-4 shrink-0" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
			<path d="M3.66489 3.58798C3.85973 2.66249 4.67621 2 5.62199 2H10.3764C11.3222 2 12.1387 2.66249 12.3335 3.58798L13.432 8.80576C12.9947 8.60929 12.5097 8.5 11.9992 8.5H3.9992C3.4887 8.5 3.00373 8.60929 2.56641 8.80576L3.66489 3.58798Z"></path>
			<path fill-rule="evenodd" clip-rule="evenodd" d="M4 10C2.89543 10 2 10.8954 2 12C2 13.1046 2.89543 14 4 14H12C13.1046 14 14 13.1046 14 12C14 10.8954 13.1046 10 12 10H4ZM12 12.75C12.4142 12.75 12.75 12.4142 12.75 12C12.75 11.5858 12.4142 11.25 12 11.25C11.5858 11.25 11.25 11.5858 11.25 12C11.25 12.4142 11.5858 12.75 12 12.75ZM9.75 12C9.75 12.4142 9.41421 12.75 9 12.75C8.58579 12.75 8.25 12.4142 8.25 12C8.25 11.5858 8.58579 11.25 9 11.25C9.41421 11.25 9.75 11.5858 9.75 12Z"></path>
		</svg>
	} else if name == "wrench" {
		<svg class="w-4 h-4 shrink-0" viewBox="0 0 16 16" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
			<path fill-rule="evenodd" clip-rule="evenodd" d="M11.5 8C13.433 8 15 6.433 15 4.5C15 4.16126 14.9519 3.83377 14.8621 3.52398C14.768 3.19913 14.365 3.13464 14.126 3.37385L12.0993 5.40148C11.921 5.57986 11.6367 5.59872 11.4459 5.43382C11.1321 5.16268 10.8384 4.86896 10.5672 4.55521C10.4024 4.36442 10.4212 4.08015 10.5995 3.90184L12.6271 1.87425C12.8662 1.63516 12.8017 1.23238 12.477 1.13817C12.1669 1.04822 11.8391 1 11.5 1C9.567 1 8 2.567 8 4.5C8 4.52183 8.0002 4.54361 8.0006 4.56535C8.01871 5.55435 7.86784 6.65835 7.08704 7.26564L1.77778 11.3951C1.28703 11.7768 1 12.3636 1 12.9853C1 14.098 1.90199 15 3.01466 15C3.63637 15 4.22325 14.713 4.60494 14.2222L8.73436 8.91296C9.34165 8.13216 10.4457 7.98129 11.4347 7.9994C11.4564 7.9998 11.4782 8 11.5 8ZM3 13.75C3.41421 13.75 3.75 13.4142 3.75 13C3.75 12.5858 3.41421 12.25 3 12.25C2.58579 12.25 2.25 12.5858 2.25 13C2.25 13.4142 2.58579 13.75 3 13.75Z"></path>
		</svg>
	}
}
```

**Step 2: Generate templ code**

Run: `templ generate`
Expected: Generates `icon_templ.go` with no errors.

**Step 3: Run tests**

Run: `go test ./... -v`
Expected: All tests pass.

**Step 4: Commit**

```
git add templates/components/icon.templ templates/components/icon_templ.go
git commit -m "feat: add Heroicon templ component for sidebar icons"
```

---

### Task 3: Render icons in sidebar section titles

**Files:**
- Modify: `templates/components/sidebar.templ:25-38` (sidebarSection template)

**Step 1: Update section title rendering to include icons**

In `templates/components/sidebar.templ`, update the `sidebarSection` template. Both the `<a>` (linked sections) and `<h4>` (unlinked sections) variants get flex layout with the icon:

```templ
templ sidebarSection(section SidebarSection, currentPath string) {
	<div class="mb-6">
		if section.Path != "" {
			<a
				href={ templ.SafeURL(section.Path) }
				class="flex items-center gap-1.5 text-sm font-semibold text-slate-800 dark:text-slate-200 hover:text-slate-900 dark:hover:text-white mb-2"
			>
				if section.Icon != "" {
					@Icon(section.Icon)
				}
				{ section.Title }
			</a>
		} else {
			<h4 class="flex items-center gap-1.5 text-sm font-semibold text-slate-800 dark:text-slate-200 mb-2">
				if section.Icon != "" {
					@Icon(section.Icon)
				}
				{ section.Title }
			</h4>
		}
		if len(section.Pages) > 0 {
			<ul class="border-l border-slate-200 dark:border-slate-700 ml-1 space-y-1">
				for _, page := range section.Pages {
					<li>
						if currentPath == page.Path {
							<a
								href={ templ.SafeURL(page.Path) }
								class="block pl-4 py-1 text-sm text-sky-600 dark:text-sky-400 border-l-2 border-sky-600 dark:border-sky-400 -ml-px font-medium"
							>
								{ page.Title }
							</a>
						} else {
							<a
								href={ templ.SafeURL(page.Path) }
								class="block pl-4 py-1 text-sm text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 border-l-2 border-transparent hover:border-slate-400 dark:hover:border-slate-500 -ml-px transition-colors"
							>
								{ page.Title }
							</a>
						}
					</li>
				}
			</ul>
		}
		if len(section.Subsections) > 0 {
			<div class="ml-3 mt-2">
				for _, sub := range section.Subsections {
					@sidebarSection(sub, currentPath)
				}
			</div>
		}
	</div>
}
```

Key changes from original:
- Line 30 (the `<a>` tag): `class="block ..."` → `class="flex items-center gap-1.5 ..."`
- Lines 31-33: Added `if section.Icon != "" { @Icon(section.Icon) }`
- Line 35 (the `<h4>` tag): same flex class change
- Lines 37-39: Same icon conditional

**Step 2: Generate templ code**

Run: `templ generate`
Expected: Generates updated `sidebar_templ.go` with no errors.

**Step 3: Build the site and verify**

Run: `just build`
Expected: Site builds successfully. Icons appear in the HTML output for sections that have them.

**Step 4: Verify icons in output HTML**

Run: `grep -l "viewBox.*0 0 16 16" dist/docs/images/index.html` (or similar)
Expected: The server SVG appears in pages under the images section.

**Step 5: Commit**

```
git add templates/components/sidebar.templ templates/components/sidebar_templ.go
git commit -m "feat: render Heroicons in sidebar section titles"
```
