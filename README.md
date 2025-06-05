# Introduction
This project aim to provide a flexible solution to add / remove Cyberark's servers into loadbalancer pools.

It consists in a Windows or Linux service that, when running, performs healtchecks and if all is OK listens to a port you specify.
Your loadbalancer check if that port is open to add (or not) the server on the pool.

If you need to remove a server from your loadbalancer pool, simply stop the service.

The healthcheck consists in three checks :
* Are the service you defined are running ?
* The connectivity to the Vault is OK ?
* The configured disks free space is more than 10%

If any check fail, the service stop listening to the network port specified.
When the checks pass again, the service resume the listening.

# Installation
 ## Windows installation
 * Run the provided msi file
 * Tune the configuration file with the necessary info
 * Start the "Cyberark Supervision"
 * That's it !

 ## Linux Installation
 TBD

# Configuration

* List the applicative services you want to be monitor.
* Then, indicate the location of the Vault.ini configuration file OR specify directly your Vault(s) IP Address(es)
* Choose the port on which the service should listen.

Sample configuration file :

``` yaml
    # Services that must be running
    services:
      - Cyber-Ark Privileged Session Manager
      - TermService
    # Manually set your Vault.ini location
    vaultIniLocation: C:\Program Files\Cyberark\PSM\Vault\Vault.ini
    # OR manually set the Vault IP address(es)
    vaultsIP: 127.0.0.1,10.20.30.40
    # Disk (or mount points) you want to check
    disks:
      - C
    # The listening port when everything is OK
    listeningPort: 38001 
    # Healthcheck is performed every X seconds
    healthCheckInterval: 10  # seconds
    # Connection to Vault will be considered timeout after X seconds
    connectionTimeout: 5     # seconds for vault connections
```

# Build / Compilation
## Build the Windows installer
A wix installer is provided but has to be compiled with Visual Studio with wix modules installed or wix directly (see build.bat)

## Build the binaries
### Windows build
```shell
GOOS=windows go build
```

### Linux build
```shell
GOOS=linux go build
```

# Code Logic
* Listen.go has the main function
* First
  * Run the function initialize to parse the config values and check everything is properly set in the file
* Periodically
  * Run checks in check_status.go