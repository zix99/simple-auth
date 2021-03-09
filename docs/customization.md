# Customization

[[toc]]

## Changing Basic Content

For starters, you probably want to enable basic whitelabeling functionality. This includes things like the company name, copyright, etc.

The relevant configuration section is `metadata`:

```yaml
metadata:
    company: "My Company"
    footer: ""
    tagline: null
    bucket: {} # Not used by default. Can be used to customize
```

## Stylesheet Changes

Any file placed in the `/static` working directory will be served as a static file to the web.  By default, *simple-auth* will attempt to load `common.css` in this directory to modify any in-line styles.

::: warning
In order to allow reading from disk, you need to set the config value `staticfromdisk: true`
:::

### Adding a Background

In the below example, place a file `bg.jpg` in the `/static/` directory, and add the following to `common.css`:

```css
body {
  background: linear-gradient(to bottom, rgba(0,0,0,0), rgba(0,0,0,.8)), url('bg.jpg');
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
}
```

## Customize Template

While most of *simple-auth*'s frontend is written in VueJS, it is still wrapped in a static template for the content.

You can override this template by providing your own `template/web/layout.tmpl` in the working directory of *simple-auth*
and setting `staticfromdisk: true`.

::: details Default Layout
<<< @/../templates/web/layout.tmpl
:::
