# bundlejs

This is allows to make JavaScript packages similar how Go does them and it bundles them to a single file.

This means a package import is always relative to a root directory. And the package consists of all files in the directory. See the "example" directory for an actual example.

Warning: this is an experimental project and doesn't play nice with the rest of JavaScript ecosystem!

## Known Issues

* imports and variables at package level can accidentally shadow each other (import package names have to be rewritten to avoid this mistake)