# website-builder

A static website builder with a CMS delivered as a single-page application (SPA) and integrated asset storage.

> ⚠️ **This project is unfinished and still at a very early stage.**
> Much of it is only sketched out or hard-coded, and interfaces and data formats may change at any time. Not intended for production use.

## Idea

Content is maintained in a browser-based CMS (page tree, content blocks, images). A separate build step turns it into a purely static website: HTML files in a directory structure that mirrors the page tree, plus all referenced assets. No CMS and no database are needed at runtime.

## Highlights

- **Two separate programs:** `spa-server` (CMS + REST API + asset storage) and `builder` (generates the static site). The CMS only needs to run while editing.
- **Templates are plain HTML.** No custom template format and no build toolchain: a template is an ordinary HTML file that uses CSS classes as placeholders. Content is bound using a Vue-like syntax (`v-for`, `v-if`, `:href`, `{{ }}`), rendered with [vuego](https://github.com/titpetric/vuego).
- **Integrated asset storage** (from [c8121/asset-storage](https://github.com/c8121/asset-storage)): images and files are picked in the CMS. At build time the builder copies them into the output's `assets/` directory automatically, named after the original file name (normalized, including German umlauts, with collision handling).
- **Internal links via `page://<ID>`:** links to other pages stay stable when pages are renamed or moved, and are resolved to real paths at build time.
- **Navigation generated from the page tree:** the main navigation is declared in the template and filled from the page hierarchy at build time.
- **Lean stack:** Go, Gin, SQLite, and an embedded Vue SPA without a Node build step. The CMS ships as a single binary (the SPA is embedded via `go:embed`, or can be loaded from a directory with `-spa`).

## Layout

```
cmd/
  spa_server/        Entry point of the CMS server
  builder/           Entry point of the site generator
internal/
  builder/           Template parser, asset loader, page tree, file name normalization
  entity/            Data model: Page, PageContent
  spa_server/        Routing of the SPA resources
shared/
  rest_api/          CMS REST API (pages, content)
  httproot/app/      CMS UI (Vue as ES modules)
  templates/default/ Example template
scripts/             build.sh / build.cmd
```

## Build and run

Requirement: Go (see `go.mod` for the version).

```bash
./scripts/build.sh        # or scripts\build.cmd on Windows
```

This produces `bin/spa-server` and `bin/builder`.

```bash
# 1. Start the CMS (default: http://127.0.0.1:9998)
bin/spa-server -data ~/website-builder

# 2. Edit content in the browser, then generate the site
mkdir -p /path/to/output
bin/builder -data ~/website-builder -tpl shared/templates/default -out /path/to/output
```

| Flag            | Meaning                                         | Default             |
|-----------------|-------------------------------------------------|---------------------|
| `-data`         | Data directory (database, asset storage)       | `~/website-builder` |
| `-out`          | Output directory (must exist)                   | `/tmp`              |
| `-tpl`          | Template directory (`builder` only)             | –                   |
| `-spa`          | Serve SPA from a directory instead of embedded  | embedded            |
| `-listen`       | Server address (`ip:port`)                      | `127.0.0.1:9998`    |
| `-cert`, `-key` | TLS certificate and key                         | –                   |

Currently the builder fetches assets from the running `spa-server` at `http://localhost:9998` (hard-coded, see below).

## Creating templates

A template is a directory containing an `index.html`. A working example lives in `shared/templates/default/`. The builder produces one output file per page from this `index.html`.

The builder relies on three CSS classes:

### 1. `.page-content` – where a page's content goes

The element with class `page-content` is emptied at build time and filled with the page's content blocks, in the order set in the CMS.

### 2. `.tpl.<ContentType>` – one template per content type

For every content block, the builder looks for an element with the classes `tpl` and the name of the content type (e.g. `tpl Text`) and uses its **innerHTML** as the template. Place these templates inside `.page-content`, so they don't end up in the output themselves.

The data of a block is a JSON object whose fields are directly available in the template:

| Content type | Fields               |
|--------------|----------------------|
| `Text`       | `text`               |
| `Image`      | `src`, `alt`, `href` |

(The CMS also contains a `Test` editor, which is meant for development only.)

Syntax (vuego, modeled on Vue):

```html
<div class="tpl Text">
    <p>{{ text }}</p>
</div>

<div class="tpl Image">
    <a v-if="href" :href="href"><img :src="src" :alt="alt"></a>
    <img v-else :src="src" :alt="alt">
</div>
```

If there is no template for a content type, the block is output as a `<pre>` with its raw content, so the problem stays visible.

**Assets:** `src` and `href` attributes pointing to `/assets/<hash>` are automatically copied from the asset storage into the output and rewritten. For values outside attributes there is also an `asset` filter (`{{ src | asset }}`).

**Page links:** `page://<ID>` in `href`/`src` is resolved to the target page's path.

### 3. `.main-nav` – navigation

The inner HTML of `.main-nav` is used as a template and filled with the page tree. `Children` is available; each item has `Name`, `Path` and its own `Children`. Recursive templates are not supported yet, so nesting depth has to be written out in the template (the default template shows three levels):

```html
<div class="main-nav">
    <ul>
        <li v-for="child in Children">
            <a :href="child.Path">{{ child.Name }}</a>
        </li>
    </ul>
</div>
```

### Minimal example

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>My site</title>
</head>
<body>
    <nav class="main-nav">
        <ul><li v-for="child in Children"><a :href="child.Path">{{ child.Name }}</a></li></ul>
    </nav>

    <main class="page-content">
        <div class="tpl Text"><p>{{ text }}</p></div>
        <div class="tpl Image"><img :src="src" :alt="alt"></div>
    </main>
</body>
</html>
```

## Current state and known gaps

Incomplete list:

- Only one template per build; per-page template selection is prepared but not implemented.
- Other files in the template directory (CSS, JS, images) are not yet copied into the output automatically (`CopyTemplate` exists but is never called).
- The home page is simply the first page in the tree.
- The builder's asset storage URL is hard-coded to `http://localhost:9998`.
- Only the content types `Text` and `Image` exist.
- The test in `test/` is outdated and currently does not compile.
- Error handling, validation and documentation are still thin.

## License

Not yet decided.
