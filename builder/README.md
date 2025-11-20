# kwctl builder image

The `Containerfile` in this directory creates a Debian "Trixie" environment that includes the dependencies required to build `kwui`. This environment can be used to cross-compile `kwui` for the Raspberry Pi.

To  build for Raspberry Pi, you can create the image by running this command from the `builder` directory:

```
podman build --platform linux/arm64 -t kwbuilder .
```

And then compile the code by running this command from the top level of the repository:

```
podman run -it --rm --platform linux/arm64 \
  -v "$PWD":/src:z -w /src -v gocache:/cache kwbuilder make
```

The `-v gocache:/cache` creates a persistent Go build and module cache. You will most definitely want this because building under emulation is substantially slower than building natively, and using the cache will reduce the build time (after the first build).
