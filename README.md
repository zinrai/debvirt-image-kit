# debvirt-image-kit

debvirt-image-kit is a tool to create Debian images for KVM virtualization using HashiCorp Packer as a backend.

## Features

* Generates a random SSH password if not provided
* Copies the specified preseed file to the `http` directory
* Generates a Packer template based on the provided parameters
* Supports loading of external Packer template files
* Installs Packer plugins
* Executes Packer to build the Debian image
* Supports generation of individual components (preseed, Packer template) or both without building the image

## Notes

- Ensure that you have sufficient permissions to run QEMU/KVM on your system.
- The generated image will be in qcow2 format, suitable for use with KVM.
- Customize the preseed file according to your needs for automated Debian installation.
- The default Packer template file is `debian.pkr.hcl.tpl`. You can customize this template or use your own.
- The default preseed template file is `preseed.cfg.tpl`. You can customize this template or use your own.

## Installation

Build the tool:

```
$ go build
```

## Usage

1. Prepare a `preseed.cfg.tpl` file with your desired Debian installation configurations.

2. Optionally, prepare a `debian.pkr.hcl.tpl` file with your desired Packer template configurations.

3. Run the tool:

   ```
   $ ./debvirt-image-kit --version 12.7.0
   ```

   This will generate both the preseed file and Packer template, then build the image.

4. To generate only specific components:

   - Generate only the preseed file:
     ```
     $ ./debvirt-image-kit --gen preseed
     ```

   - Generate only the Packer template:
     ```
     $ ./debvirt-image-kit --gen packer
     ```

   - Generate both preseed and Packer template without building the image:
     ```
     $ ./debvirt-image-kit --gen all
     ```

5. Additional options:

   - Specify a custom SSH username:
     ```
     $ ./debvirt-image-kit --ssh-username myuser
     ```

   - Specify a custom SSH password (if not provided, a random password will be generated):
     ```
     $ ./debvirt-image-kit --ssh-password mypassword
     ```

   - Specify custom disk size and memory:
     ```
     $ ./debvirt-image-kit --disk-size 30000M --memory 4096
     ```

   - Use a custom Packer template file:
     ```
     $ ./debvirt-image-kit --packer-template my-custom-template.pkr.hcl.tpl
     ```

   - Use a custom preseed template file:
     ```
     $ ./debvirt-image-kit --preseed-file my-custom-preseed.cfg.tpl
     ```

## Example Output

```
$ ./debvirt-image-kit --version 12.7.0
Starting debvirt-image-kit...
Generated random SSH password: pBdJ4WPqcYK7zIOH
Installing Packer plugins...
Running Packer to build the image...
qemu.debian: output will be in this color.

==> qemu.debian: Retrieving ISO
==> qemu.debian: Trying https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.7.0-amd64-netinst.iso
...
Build 'qemu.debian' finished after 5 minutes 31 seconds.

==> Wait completed after 5 minutes 31 seconds

==> Builds finished. The artifacts of successful builds are:
--> qemu.debian: VM files in directory: output
debvirt-image-kit: Debian image generated successfully!
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
