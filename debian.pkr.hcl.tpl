packer {
  required_plugins {
    qemu = {
      version = ">= 1.1.0"
      source  = "github.com/hashicorp/qemu"
    }
  }
}

source "qemu" "debian" {
  iso_url      = "{{.ISOURL}}"
  iso_checksum = "{{.ISOChecksum}}"
  output_directory = "{{.OutputDir}}"
  shutdown_command = "echo '{{.SSHPassword}}' | sudo -S /sbin/shutdown -hP now"
  disk_size        = "{{.DiskSize}}"
  memory           = "{{.MemorySize}}"
  format           = "qcow2"
  accelerator      = "kvm"
  http_directory   = "http"
  ssh_username     = "{{.SSHUsername}}"
  ssh_password     = "{{.SSHPassword}}"
  ssh_timeout      = "20m"
  vm_name          = "debian-{{.DebianVersion}}-{{.DebianArch}}"
  net_device       = "virtio-net"
  disk_interface   = "virtio"
  boot_wait        = "5s"
  boot_command     = [
    "<esc><wait>",
    "auto ",
    "url=http://{{ "{{ .HTTPIP }}" }}:{{ "{{ .HTTPPort }}" }}/preseed.cfg ",
    "<enter>"
  ]
}

build {
  sources = ["source.qemu.debian"]
}
