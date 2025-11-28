

## topologie logique ou physique


topologie : architecture conception
- logique :
  - on peut avoir des réseaux virtuels
- physique :
  - ce que l'on a physiquement, l'équipement (réseau, machines...)


CIDR : Nombre de bits pour le masque de sous-réseau


## Architechture réseau + avancée

Service RH      Service Compta
     \          / 
      \        / 
       \      /
        \    /
        SWITCH ------------------ Wi-Fi ---------- Appareils Mobiles
           |   \
           |       \
           |          \
           |             \
      LAMP Server      Clients




SI on ne veut pas que le service compta puisse accéder au form mais que la RH si,

on peut : 
- mettre en place un filtrage à ARP sur le switch (compliqué, il faudrait un switch spécial) -> associe @MAC et @IP des machines
- Mettre le LAMP dans  un réseau VLAN privé + Un VLAN pour la compta + un VLAN pour la RH (voir même pour le Wifi et les clients) On segmente tous les services de l'entreprise
   -  Chaque VLAN a son IP
   -  Cette méthode est mieux que la première car elle permet de ne pas surcharger les réseau
   -  On rajoute un routeur physique
       - On crée des sous-interface (??) sur le routeur
           - on veut permetttre les communications VLAN Client -> VLAN Serveur LAMP et VLAN Serveur LAMP -> Client
           - On met donc en place  des passerelles :
               - Gi0/0.100 pour le VLAN CLIENT dont l'IP est 192.168.100.1.254/24
               - Gi0/0.10 pour le VLAN Serveur LAMP dont l'IP est 192.168.10.1.254/24


- Protocole 802 1.Q : encapsulation dot 1q
- Du routage interVLAN


On augmente a difficulté : On veut que COmpta puisse communiquer à tout SAIF le Serveur LAMP

- On ajoute une règle de filtrage dans le switch de niveai 3 OU un Switch de niveau 2 + un routeur

CISCO : 
   - 2960 : Niveau 2 (premier chiffre)
   - 3560 : Niveau 3 (premier chiffre)
HP:
    - 2824 : Niveau 2
    - 3... : Niveau 3


ex : Cisco Nexu 48 ports Gibabit 5 SFP 3750 POE (Niveau 3) -> 1800€ neuf - 110 € d'occasion : Toujours prendre d'occsaion pro (ou faire les bennes)
Attention SFP != SFP+ (10Gbit)
Reconditionnement d'équipements (3 boites à Amsterdam) -> Garantie à vie
Les établissements public ne peuvent pas prendre du reconditionné

Pour les assos on peut demadner de récupérer le matériel des entreprises (dons) - ça leur éviter de payer des frais de reconditionnement super chers au poids



## Pfsense

         WAN
          |
          |
        Pfsense
       /        \
   LAN/          \ DMZ (Uniquement ce qui est public, pas la BDD, le WAN peut uniquement accéder à la DMZ et pas au LAN)
     /            \
   BDD           Serveur Web


à partir de 16Go de RAM on peut tenter APcket Tracer  GN3



## Répartition de charge
- avec docker : ansible


Sinon: archi réseau : 



pc ------------ swtitch ------HA Proxy (load balancing)
                                |
                                |
                              Swtich
                              /  |  \
                           /     |    \
                        /        |      \
                     LAMP 1    LAMP 2   LAMP 3