This project aim to provide a flexible solution to add / remove Cyberark's servers into loadbalancer pools which is always a pain :)


First list the applicative services you want to be running 
Then, indicate the location of the Vault.ini configuration file
Finally, give the port on which the service should listen.

A wix installer is provided but has to be compiled with Visual Studio with wix modules installed


Sample configuration file :

    services:
      - Cyber-Ark Privileged Session Manager
      - TermService
    vaultIniLocation: C:\Program Files\Cyberark\PSM\Vault\Vault.ini
    disks:
      - C
    listeningPort: 38001 




