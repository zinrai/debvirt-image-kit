# Localization
d-i debian-installer/locale string C.UTF-8
d-i debian-installer/language string en
d-i debian-installer/country string JP
d-i keyboard-configuration/xkb-keymap select us

# Network configuration
d-i netcfg/choose_interface select auto
d-i netcfg/get_hostname string debian
d-i netcfg/get_domain string unassigned-domain

# Mirror settings
d-i mirror/country string manual
d-i mirror/http/hostname string deb.debian.org
d-i mirror/http/directory string /debian
d-i mirror/http/proxy string

# Account setup
d-i passwd/root-login boolean false
d-i passwd/user-fullname string Debian User
d-i passwd/username string {{.SSHUsername}}
d-i passwd/user-password password {{.SSHPassword}}
d-i passwd/user-password-again password {{.SSHPassword}}

# Clock and time zone setup
d-i clock-setup/utc boolean true
d-i time/zone string UTC

# Partitioning
d-i partman-auto/method string regular
d-i partman-auto/choose_recipe select atomic
d-i partman-partitioning/confirm_write_new_label boolean true
d-i partman/choose_partition select finish
d-i partman/confirm boolean true
d-i partman/confirm_nooverwrite boolean true

# Package selection
tasksel tasksel/first multiselect none
d-i pkgsel/include string openssh-server sudo
d-i pkgsel/upgrade select full-upgrade

# Bootloader installation
d-i grub-installer/only_debian boolean true
d-i grub-installer/bootdev string /dev/vda

# Finishing up the installation
d-i finish-install/reboot_in_progress note

# Run custom commands
d-i preseed/late_command string \
    in-target apt-get update; \
    in-target apt-get upgrade -y; \
    in-target apt-get clean
