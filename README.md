# website-builder

A static website builder with a CMS delivered as a single-page application (SPA) and integrated asset storage.

> ⚠️ **This project is unfinished and still at an early stage.**
> Much of it is only sketched out or hard-coded, and interfaces and data formats may change at any time. Not intended for production use.

## Idea

Content is maintained in a browser-based CMS (page tree, content blocks, images, site settings). A separate build step turns it into a purely static website: HTML files in a directory structure that mirrors the page tree, plus all referenced assets. No CMS and no database are needed at runtime.

## Highlights

- **Separate programs:** `spa-server` (CMS + REST API + asset storage), `builder` (generates the static site) and `site-server` (a tiny static file server to preview the result). The CMS only needs to run while editing and building.
- **Templates are plain HTML.** No custom template format and no build toolchain: a template is an ordinary HTML file that uses CSS classes to mark where things go. Data is bound using a Vue-like syntax (`v-for`, `v-if`, `:href`, `v-html`, `{{ }}`), rendered with [vuego](https://github.com/titpetric/vuego).
- **Site settings:** key/value pairs maintained in the CMS (e.g. `Title`, `SubTitle`) can be used as `{{ Name }}` placeholders anywhere in a template.
- **Content types:** rich text (`Html`), plain text (`Text`) and images (`Image`). A template decides how each type is rendered.
- **Integrated asset storage** (from [c8121/asset-storage](https://github.com/c8121/asset-storage)): images and files are picked in the CMS. At build time they are copied into the output's `assets/` directory automatically, named after the original file name (normalized, including German umlauts, with collision handling).
- **Post-processing:** asset URLs are resolved, internal links via `page://<ID>` stay stable when pages are renamed or moved, and an optional table of contents is generated from the headings of each page.
- **Navigation generated from the page tree**, declared in the template.
- **Ordering:** pages, content blocks and site settings can be reordered in the CMS; the order is used for navigation and output.
- **Lean stack:** Go, Gin, SQLite, and an embedded Vue SPA without a Node build step. The CMS ships as a single binary (the SPA is embedded via `go:embed`, or can be loaded from a directory with `-spa`).

## Layout

```
cmd/
  spa_server/        Entry point of the CMS server
  builder/           Entry point of the site generator
  site_server/       Static file server for previewing the generated site
internal/
  builder/           Template parser, post-processors, asset loader, page tree, file name normalization
  entity/            Data model: Page, PageContent, SiteConfig
  spa_server/        Routing of the SPA resources
shared/
  rest_api/          CMS REST API (pages, content, site settings)
  httproot/app/      CMS UI (Vue as ES modules, includes a rich-text editor)
  templates/default/ Plain example template that shows all features
  templates/ingrid/  Example template based on Bootstrap
scripts/             build.sh / build.cmd
```

## Build and run

Requirement: Go (see `go.mod` for the version).

```bash
./scripts/build.sh        # or scripts\build.cmd on Windows
```

This produces `bin/spa-server`, `bin/builder` and `bin/site-server`.

```bash
# 1. Start the CMS (default: http://127.0.0.1:9998)
bin/spa-server -data ~/website-builder

# 2. Edit content in the browser, then generate the site
#    (the CMS must keep running, the builder fetches assets from it)
bin/builder -data ~/website-builder -tpl shared/templates/ingrid -out /path/to/output

# 3. Preview the result (default: http://localhost:8080)
bin/site-server -root /path/to/output
```

Options of `spa-server` and `builder`:

| Flag            | Meaning                                          | Default                       |
|-----------------|--------------------------------------------------|-------------------------------|
| `-data`         | Data directory (database, asset storage)         | `~/website-builder`           |
| `-out`          | Output directory (created if missing)            | `<data>/tmp/generated`        |
| `-tpl`          | Template directory (required for `builder`)      | –                             |
| `-spa`          | Serve SPA from a directory instead of embedded   | embedded                      |
| `-listen`       | Server address (`ip:port`)                       | `127.0.0.1:9998`              |
| `-cert`, `-key` | TLS certificate and key                          | –                             |

Options of `site-server`: `-root` (directory to serve, required) and `-listen` (default `localhost:8080`).

Currently the builder fetches assets from the running `spa-server` at `http://localhost:9998` (hard-coded, see below).

## Creating templates

A template is a directory with an `index.html`. Two working examples live in `shared/templates/`: `default` shows every feature without styling, `ingrid` is a Bootstrap-based layout. The builder produces one output file per page from the template's `index.html`.

### How the output is produced

1. The **whole template directory** is copied into the output directory (CSS, JavaScript, images, …). The template's own `index.html` is copied too, but is overwritten by the generated home page.
2. For every page, `index.html` is loaded and `{{ Name }}` placeholders are replaced with site settings (see below).
3. Content is inserted at `.page-content`, the navigation at `.main-nav`.
4. Post-processing resolves assets and links and fills `.toc`.
5. The result is written to `<output>/<page path>/index.html`; the home page goes to `<output>/index.html`.

Because pages live in subdirectories, reference files of the template with **root-relative paths** (`/css/app.css`, `/images/logo.png`), not relative ones. The same applies to the generated links: they are root-relative, so the site has to be served from the root of a domain.

### Site settings: `{{ Name }}`

Settings are maintained in the CMS (route `/config`) as name/value pairs. A fresh database starts with `Title` and `SubTitle`. In a template, `{{ Title }}` is replaced with the value before anything else is parsed, anywhere in the file (also inside `<title>`, `<meta>` or attributes):

```html
<title>{{ Title }} | {{ SubTitle }}</title>
<img src="{{ Logo }}" alt="{{ Title }}">
```

Placeholders without a matching setting are left untouched. Avoid setting names that are also used by the template data below (`Page`, `Meta`, `Content`, `Children`, …), otherwise they would be replaced too early.

### Template syntax

Dynamic parts use vuego, modeled on Vue: `{{ expression }}`, `v-for`, `v-if`/`v-else`, `v-html`, `:attribute="expression"`, `<template>` and filters such as `value | asset`. Which data is available depends on where the template is used:

| Where                  | Data                                                                 |
|------------------------|----------------------------------------------------------------------|
| `.tpl.<Type>` (content)| `Page`, `Meta`, `Content`                                            |
| `.main-nav`            | `Children` (page tree)                                               |
| `.toc`                 | `Children` (headings) and `Page`                                     |

Available filters: `asset` (copies an asset into the output and returns its path) and `replace` (`{{ text | replace("a", "b") }}`).

### 1. `.page-content` – where a page's content goes

The element with class `page-content` is emptied at build time and filled with the page's content blocks, in the order set in the CMS.

### 2. `.tpl.<ContentType>` – one template per content type

For every content block, the builder looks for an element with the classes `tpl` and the name of the content type (e.g. `tpl Text`) and uses its **innerHTML** as the template. Place these templates inside `.page-content`, so they don't end up in the output themselves. An element can serve several types by listing them all: `<div class="tpl Text Html">`.

Inside a content template, three objects are available:

| Name      | Content                                                                 |
|-----------|-------------------------------------------------------------------------|
| `Content` | The fields of the block, as filled in by the editor (see table below)   |
| `Meta`    | The block itself, e.g. `Meta.Name`, `Meta.ContentType`                  |
| `Page`    | The page being generated, e.g. `Page.Name`, `Page.Sort`                 |

| Content type | Editor    | Fields in `Content`                    |
|--------------|-----------|----------------------------------------|
| `Html`       | Rich text | `text` (HTML, use `v-html`)            |
| `Text`       | Plain text| `text` (plain text, escaped)           |
| `Image`      | Image     | `src`, `alt`, `class`, `href`          |

(The CMS also contains a `Test` editor, which is meant for development only.)

```html
<div class="tpl Html">
    <h2>{{ Meta.Name }}</h2>
    <div v-html="Content.text"></div>
</div>

<div class="tpl Text">
    <p style="white-space: pre-line">{{ Content.text }}</p>
</div>

<div class="tpl Image">
    <a v-if="Content.href" :href="Content.href">
        <img :src="Content.src | asset" :alt="Content.alt" :class="Content.class"/>
    </a>
    <img v-else :src="Content.src | asset" :alt="Content.alt" :class="Content.class"/>
</div>
```

If there is no template for a content type, the block is output as a `<pre>` with its raw content, so the problem stays visible. If binding fails, the error is shown in a `<pre>` as well.

### 3. `.main-nav` – navigation

The inner HTML of every `.main-nav` element is used as a template and filled with the page tree. `Children` holds the top-level pages; each item has:

| Field         | Content                                                  |
|---------------|----------------------------------------------------------|
| `Name`        | Name of the page                                         |
| `Path`        | Link to the page (root-relative)                         |
| `PathName`    | Directory name of the page                               |
| `IsHomePage`  | `true` for the home page (use `'/'` as link)             |
| `Page`        | The page itself (`Page.Id`, `Page.Sort`, …)              |
| `CurrentPage` | The page currently being generated (e.g. to highlight it)|
| `Children`    | Sub-pages, same structure                                |

Recursive templates are not supported yet, so nesting depth has to be written out (the default template shows three levels). There may be several `.main-nav` elements, for example a mobile and a desktop variant.

```html
<div class="main-nav">
    <ul>
        <li v-for="child in Children">
            <a :href="child.IsHomePage ? '/' : child.Path">{{ child.Name }}</a>
        </li>
    </ul>
</div>
```

### 4. `.toc` – table of contents (optional)

If the template contains an element with class `toc`, the builder collects all headings (`h1` to `h6`) of the generated page, gives each an `id` derived from its text and fills `.toc` with the result. Headings are nested according to their level. Each item has:

| Field      | Content                                          |
|------------|--------------------------------------------------|
| `Caption`  | Text of the heading                              |
| `Name`     | The generated `id`, use it as anchor (`'#' + child.Name`) |
| `Level`    | Heading level (1 to 6)                           |
| `Number`   | Running number of the heading on the page        |
| `Page`     | The page the heading belongs to                  |
| `Children` | Nested headings                                  |

As with the navigation, nesting has to be written out; the default template shows five levels:

```html
<div class="toc">
    <ul>
        <li v-for="child in Children">
            <a :href="'#' + child.Name">{{ child.Caption }}</a>
            <ul>
                <li v-for="child in child.Children">
                    <a :href="'#' + child.Name">{{ child.Caption }}</a>
                </li>
            </ul>
        </li>
    </ul>
</div>
```

Note that all headings of the final document are included, also those outside `.page-content`.

### Assets and links

After the content has been inserted, every `src` and `href` attribute in the document is post-processed:

- Values pointing to `/assets/<hash>` are copied from the asset storage into the output's `assets/` directory and rewritten. The `asset` filter does the same explicitly (`Content.src | asset`) and is needed for values that are not plain attributes.
- `page://<ID>` is replaced with the path of the target page, so links survive renaming and moving pages.
- Absolute paths within the output directory are turned into root-relative links.

### Minimal example

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>{{ Title }}</title>
</head>
<body>
    <nav class="main-nav">
        <ul>
            <li v-for="child in Children">
                <a :href="child.IsHomePage ? '/' : child.Path">{{ child.Name }}</a>
            </li>
        </ul>
    </nav>

    <main class="page-content">
        <div class="tpl Html"><h2>{{ Meta.Name }}</h2><div v-html="Content.text"></div></div>
        <div class="tpl Text"><p>{{ Content.text }}</p></div>
        <div class="tpl Image"><img :src="Content.src | asset" :alt="Content.alt"></div>
    </main>
</body>
</html>
```

## Current state and known gaps

Incomplete list:

- Only one template per build; per-page template selection is prepared but not implemented.
- The home page is simply the first page in the tree.
- The builder's asset storage URL is hard-coded to `http://localhost:9998`, so `spa-server` has to run on that address while building.
- Generated links are root-relative; the site cannot be served from a sub-path.
- Navigation and table of contents need their nesting depth written out in the template.
- The test in `test/` is outdated and currently does not compile.
- Error handling, validation and documentation are still thin.

## Third-party components

This project builds on the following components. They are not covered by this project's license; each one remains under its own license, which you can find in the linked repository.

### Go modules (direct dependencies)

| Component | Used for | Repository |
|-----------|----------|------------|
| goquery | Parsing and manipulating the HTML templates | https://github.com/PuerkitoBio/goquery |
| vuego | Rendering the Vue-like template syntax | https://github.com/titpetric/vuego |
| Gin | HTTP server and routing | https://github.com/gin-gonic/gin |
| otiai10/copy | Copying the template directory into the output | https://github.com/otiai10/copy |
| golang.org/x/net | HTML parser | https://github.com/golang/net |
| golang.org/x/text | Unicode normalization of file names | https://github.com/golang/text |

### Libraries by the same author

| Component | Used for | Repository |
|-----------|----------|------------|
| asset-storage | Asset storage, its REST API and UI parts | https://github.com/c8121/asset-storage |
| go_db_entities | Simple entity persistence | https://github.com/c8121/go_db_entities |
| go_vue_spa | Serving the Vue SPA | https://github.com/c8121/go_vue_spa |

### Frontend (CMS)

| Component | Version | Repository |
|-----------|---------|------------|
| Vue | 3.5.41 | https://github.com/vuejs/core |
| Vue Router | 5.2.0 | https://github.com/vuejs/router |
| Bootstrap | 5.3.8 | https://github.com/twbs/bootstrap |
| Bootstrap Icons | 1.13.1 | https://github.com/twbs/icons |
| Neiki Editor (rich-text editor, bundled in `shared/httproot/app/js/neiki-editor/`) | 1.0 | https://github.com/neikiri/neiki-editor |

Vue, Vue Router, Bootstrap and Bootstrap Icons are referenced by `website-builder.html` but are not part of this repository; they are delivered through the `go_vue_spa` and `asset-storage` modules.

### Templates

| Component | Version | Repository |
|-----------|---------|------------|
| Bootstrap (CSS and JavaScript, bundled in `shared/templates/ingrid/css/`) | 5.3.8 | https://github.com/twbs/bootstrap |

### Transitive dependencies

Go modules pulled in indirectly are listed in `go.mod` and `go.sum`, among them Bleve (full-text index, https://github.com/blevesearch/bleve) and the pure-Go SQLite driver `modernc.org/sqlite` (https://gitlab.com/cznic/sqlite). Their licenses apply as well.

## License

AGPLv3, see [`LICENSE`](LICENSE).
