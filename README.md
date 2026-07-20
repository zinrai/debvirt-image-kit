# debvirt-image-kit

debvirt-image-kit creates Debian images for KVM virtualization using HashiCorp Packer as a backend.

It renders two build inputs from local templates and drives Packer over them:

- `preseed.cfg`, the answer file the installer fetches over HTTP for an unattended Debian install
- a Packer template that boots the Debian netinst ISO under QEMU/KVM and produces a qcow2 image

## Model

The work is split into two commands so the rendered inputs can be inspected before a multi-minute build runs.

1. `render` writes the build inputs into a `build-debian-<version>-<arch>` directory
2. `build` runs `packer init` and `packer build` against the image with that `--version` and `--arch`

```
debvirt-image-kit render --version 13.6.0
debvirt-image-kit build --version 13.6.0
```

Both commands identify the image by `--version` and `--arch`, so `build` needs nothing else. The per-version directory lets several images coexist in one working directory. Re-run `render` to change anything.

Run `debvirt-image-kit render -h` for flags.

## Templates

`preseed.cfg.tpl` and `debian.pkr.hcl.tpl` are loaded from the working directory by default. Point `--preseed-file` and `--packer-template` at your own copies to customize the install or the Packer source.

## Requirements

- `packer` on `PATH` (needed by `build`, not by `render`)
- Permission to run QEMU/KVM

If `--ssh-password` is omitted, a random password is generated and printed. It is set on the created user, and Packer uses it to connect over SSH during the build.

## ISO source

`--version` and `--arch` alone determine the ISO source. The tool appends
`<version>/<arch>/iso-cd/` to `--iso-base-url`, which defaults to the
`debian-cd` tree. That tree only carries the current release, so building an
older point release means pointing the prefix at the archive:

```
debvirt-image-kit render --version 12.7.0 \
  --iso-base-url https://cdimage.debian.org/cdimage/archive/
```

Because the base URL is only a prefix, any Debian CD mirror with the same
layout works the same way. Use this to route around a slow or unreachable
default host:

```
debvirt-image-kit render --version 13.6.0 \
  --iso-base-url https://ftp.jaist.ac.jp/pub/Linux/debian-cd/
```

## License

This project is licensed under the [MIT License](./LICENSE).
