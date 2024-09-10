package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var (
	debianVersion    string
	debianArch       string
	outputDir        string
	diskSize         string
	memorySize       string
	sshUsername      string
	sshPassword      string
	isoBaseURL       string
	checksumFileName string
	isoFileName      string
	preseedFile      string
	genOption        string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "debvirt-image-kit [command]",
		Short: "Create Debian images for KVM using Packer",
		Long:  `debvirt-image-kit is a tool to create Debian images for KVM virtualization using HashiCorp Packer as a backend.`,
		Run:   runGenerator,
	}

	rootCmd.Flags().StringVarP(&debianVersion, "version", "v", "11.6.0", "Debian version")
	rootCmd.Flags().StringVarP(&debianArch, "arch", "a", "amd64", "Debian architecture")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "output", "Output directory")
	rootCmd.Flags().StringVar(&diskSize, "disk-size", "20000M", "Disk size (e.g., 5000M, 10G)")
	rootCmd.Flags().StringVar(&memorySize, "memory", "1024", "Memory size (e.g., 2048)")
	rootCmd.Flags().StringVar(&sshUsername, "ssh-username", "debian", "SSH username")
	rootCmd.Flags().StringVar(&sshPassword, "ssh-password", "", "SSH password (if not provided, a random password will be generated)")
	rootCmd.Flags().StringVar(&isoBaseURL, "iso-base-url", "https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/", "Base URL for ISO download")
	rootCmd.Flags().StringVar(&checksumFileName, "checksum-file", "SHA256SUMS", "Checksum file name")
	rootCmd.Flags().StringVar(&isoFileName, "iso-file", "", "ISO file name (e.g., debian-11.6.0-amd64-netinst.iso)")
	rootCmd.Flags().StringVar(&preseedFile, "preseed-file", "preseed.cfg.tpl", "Path to the preseed template file")
	rootCmd.Flags().StringVar(&genOption, "gen", "", "Generate option: 'preseed', 'packer', or 'all' (default: build image)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runGenerator(cmd *cobra.Command, args []string) {
	fmt.Println("Starting debvirt-image-kit...")

	// Check if Packer is installed
	if err := checkPackerInstallation(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate random password if not provided
	if sshPassword == "" {
		sshPassword = generateRandomPassword(16)
		fmt.Printf("Generated random SSH password: %s\n", sshPassword)
	}

	switch genOption {
	case "preseed":
		if err := generatePreseedFile(); err != nil {
			fmt.Printf("Error generating preseed file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Preseed file generated successfully.")
	case "packer":
		if _, err := generatePackerTemplate(); err != nil {
			fmt.Printf("Error generating Packer template: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Packer template generated successfully.")
	case "all":
		if err := generatePreseedFile(); err != nil {
			fmt.Printf("Error generating preseed file: %v\n", err)
			os.Exit(1)
		}
		if _, err := generatePackerTemplate(); err != nil {
			fmt.Printf("Error generating Packer template: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Preseed file and Packer template generated successfully.")
	case "":
		// Default behavior: generate both and build image
		if err := generatePreseedFile(); err != nil {
			fmt.Printf("Error generating preseed file: %v\n", err)
			os.Exit(1)
		}
		packerTemplateFile, err := generatePackerTemplate()
		if err != nil {
			fmt.Printf("Error generating Packer template: %v\n", err)
			os.Exit(1)
		}

		// Install Packer plugins
		fmt.Println("Installing Packer plugins...")
		initCmd := exec.Command("packer", "init", packerTemplateFile)
		initCmd.Stdout = os.Stdout
		initCmd.Stderr = os.Stderr
		if err := initCmd.Run(); err != nil {
			fmt.Printf("Error installing Packer plugins: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Running Packer to build the image...")
		// Run Packer
		packerCmd := exec.Command("packer", "build", packerTemplateFile)
		packerCmd.Stdout = os.Stdout
		packerCmd.Stderr = os.Stderr
		if err := packerCmd.Run(); err != nil {
			fmt.Printf("Error running Packer: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("debvirt-image-kit: Debian image generated successfully!")
	default:
		fmt.Printf("Invalid gen option: %s. Use 'preseed', 'packer', 'all', or omit for default behavior.\n", genOption)
		os.Exit(1)
	}
}

func checkPackerInstallation() error {
	_, err := exec.LookPath("packer")
	if err != nil {
		return fmt.Errorf("Packer is not installed or not in the system PATH. Please install Packer and try again. Error: %v", err)
	}
	return nil
}

func generateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	password := make([]rune, length)
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))]
	}
	return string(password)
}

func generatePreseedFile() error {
	if err := os.MkdirAll("http", 0755); err != nil {
		return fmt.Errorf("error creating http directory: %v", err)
	}

	source, err := os.Open(preseedFile)
	if err != nil {
		return err
	}
	defer source.Close()

	content, err := io.ReadAll(source)
	if err != nil {
		return err
	}

	tmpl, err := template.New("preseed").Parse(string(content))
	if err != nil {
		return err
	}

	destination, err := os.Create(filepath.Join("http", "preseed.cfg"))
	if err != nil {
		return err
	}
	defer destination.Close()

	data := struct {
		SSHUsername string
		SSHPassword string
	}{
		SSHUsername: sshUsername,
		SSHPassword: sshPassword,
	}

	err = tmpl.Execute(destination, data)
	if err != nil {
		return err
	}

	return nil
}

func generatePackerTemplate() (string, error) {
	packerTemplateFile := fmt.Sprintf("debian-%s-%s.pkr.hcl", debianVersion, debianArch)

	tmpl := template.Must(template.New("packer").Parse(`
packer {
  required_plugins {
    qemu = {
      version = ">= 1.0.0"
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
`))

	f, err := os.Create(packerTemplateFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// If ISO file name is not provided, construct it
	if isoFileName == "" {
		isoFileName = fmt.Sprintf("debian-%s-%s-netinst.iso", debianVersion, debianArch)
	}

	isoURL := fmt.Sprintf("%s%s", isoBaseURL, isoFileName)
	isoChecksum := fmt.Sprintf("file:%s%s", isoBaseURL, checksumFileName)

	err = tmpl.Execute(f, struct {
		DebianVersion string
		DebianArch    string
		OutputDir     string
		DiskSize      string
		MemorySize    string
		SSHUsername   string
		SSHPassword   string
		ISOURL        string
		ISOChecksum   string
	}{
		DebianVersion: debianVersion,
		DebianArch:    debianArch,
		OutputDir:     outputDir,
		DiskSize:      diskSize,
		MemorySize:    memorySize,
		SSHUsername:   sshUsername,
		SSHPassword:   sshPassword,
		ISOURL:        isoURL,
		ISOChecksum:   isoChecksum,
	})

	if err != nil {
		return "", err
	}

	return packerTemplateFile, nil
}
