# Woke documentation site

Woke uses the [material](https://squidfunk.github.io/mkdocs-material/) MkDocs theme.

You can find the live site at <https://docs.getwoke.tech>.

See [mkdocs.yml](../mkdocs.yml) for details.

## Running locally

```bash
docker build -t woke-mkdocs -f docs/Dockerfile .
docker run --rm -v `pwd`:/workdir -w /workdir -p 8000:8000 woke-mkdocs
```

You can now access the site locally via `http://localhost:8000`. Auto-reload is enabled

## Deploying

Any changes made to the `docs` directory will be automatically deployed on merges to `main`.
