# jo

Yet another (experimental) package manager for JavaScript. Although it has a simpler model than npm, bower etc.

**Warning: this is an experimental project and may not play nice with the rest of JavaScript ecosystem! You'll be better off using jspm in a production environment.**

First each package is defined relative to some root directory.

For example:

    |-- bundle.js
    |-- index.html
    `-- js
        |-- Application.js
        `-- page
            |-- items
            |   |-- basic
            |   |   `-- paragraph.js
            |   `-- registry.js
            |-- page.js
            `-- view.js

The root directory in this case is the `js` directory. To generate the `bundle.js` file run inside the `example` directory:
    
    jo -jopath="js" build bundle.js

The first parameter is the root directory and the second parameter is the output file.

Each package is defined by their folder: e.g. there are packages `page`, `page/items`, `page/items/basic`. The package is comprised of all files in that folder. Each file can import a package via `import "package/name"` and export variables via `export VariableName`. The exported variables then will be accessible via:

    // in person/name/value.js
    var Value = "Some Name";
    export Value;
    
    // in page/other.js
    import "person/name";
    console.log(name.Value);

All exported variables/functoins will be accessible via the imported packages folder name.

Of course there can be no circular dependencies - enforcing such constraint makes code much simpler. There is no builtin versioning, all dependencies should downloaded and fixed.

## Known Issues

* imports and variables at package level can accidentally shadow each other (import package names have to be rewritten to avoid this mistake)