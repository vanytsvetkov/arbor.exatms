# exaTMS

exaTMS – exo [Arbor Networks©](https://www.netscout.com/product/arbor-threat-mitigation-system) Threat Mitigation System tool used to patch mitigation announces using **exabgp** and a custom golang binary before exporting announces to BIRD.

### Usage scenario
The main scenario in which `exatms` is helpful is depicted in the figure below. Admin may need to modify the attributes of a mitigation announce (ex. BGP.as_path) received from TMS Appliance before exporting (reannounce) it to BIRD's RIB.
```text
   ------------------------------------------------------------------------
      Host1                                         Host2                 
                             |                                            
   -----------+              |    +----------+                  +----------
       TMS    |              |    |  ExaBGP  | :1790       :179 |   BIRD   
    Appliance | 179:BGP:1790 ===> |  Daemon  | <=== exaTMS ===> |  Daemon  
   -----------+              |    +----------+                  +----------
                             |                                            
                    FIREWALL (only port T:1790 is open)                   
                                                                          
   ------------------------------------------------------------------------
```


### Examples
#### Ex. 1 Full cycle of testing
First of all, let's build the sources and run the binary:
```shell
$ bash -c 'if [ -d exasrc ]; then cd exasrc; else echo "Command should be run at the top level of this repository" && exit; fi 
go get; go build -o ../exatms;
if [ ! -p /etc/exabgp/exabgp.in ]; then mkfifo /etc/exabgp/exabgp.in; fi
sleep 1e+9 > /etc/exabgp/exabgp.in & &> /dev/null;
../exatms --neighbor "127.0.0.1" --peer_as "65500" --loglevel "DEBUG" --logfile "/var/log/exabgp/exatms.log" < /etc/exabgp/exabgp.in'
```
You can then send few **exabgp's** messages over the pipe to `exatms` standard input:
```shell
$ echo '{ "exabgp": "4.0.1", "time": 1678483744.0192661, "host" : "test-host", "pid" : 2952173, "ppid" : 1, "counter": 1, "type": "update", "neighbor": { "address": { "local": "172.16.0.1", "peer": "172.16.0.2" }, "asn": { "local": 65500, "peer": 65501 } , "direction": "receive", "message": { "update": { "attribute": { "origin": "igp", "local-preference": 100, "community": [ [ 0, 65500 ], [0, 65501] ] }, "announce": { "ipv4 unicast": { "172.16.0.2": [ { "nlri": "10.0.0.0/24" }, { "nlri": "10.0.1.0/24" } ] } } } } } }' > /etc/exabgp/exabgp.in
$ echo '{ "exabgp": "4.0.1", "time": 1678544034.886104, "host" : "test-host", "pid" : 3167756, "ppid" : 1, "counter": 4, "type": "update", "neighbor": { "address": { "local": "172.16.0.1", "peer": "172.16.0.2" }, "asn": { "local": 65500, "peer": 65501 } , "direction": "receive", "message": { "update": { "attribute": { "origin": "igp", "local-preference": 100, "community": [ [ 0, 65500 ], [0, 65501] ] }, "withdraw": { "ipv4 unicast": [ { "nlri": "10.0.0.0/24" } ] } } } } }' > /etc/exabgp/exabgp.in
```
Now check stdout of `exatms` and log file:
```shell
$ tail -f /var/log/exabgp/exatms.log
```
#### Ex. 2 Integration into ExaBGP's & BIRD's .conf files 
Below is an example **exabgp's** configuration file:
```shell
process exatms {
    run /etc/exabgp/exatms/exatms --neighbor "127.0.0.1" --peer_as "65500" --loglevel "INFO" --logfile "/var/log/exabgp/exatms.log";
    encoder json;
}

neighbor 172.16.0.2 {
	description "TMS Appliance";

	router-id 172.16.0.1;
	local-as 65500;
	local-address 172.16.0.1;
	peer-as 65500;
    
        auto-flush true;
        capability {
            route-refresh enable;
        }

	family {
		ipv4 unicast;
	}

	api {
              processes [exatms];
              receive {
                parsed;
                update;
                keepalive;
              }
        }
}


neighbor 127.0.0.1 {
	description "BIRD Daemon";

	router-id 127.0.0.1;
	local-as 65500;
	local-address 127.0.0.1;
	peer-as 65500;

	family {
		ipv4 unicast;
	}

        api {
              processes [exatms];
        }
}
```
And the obvious part of the BGP session configuration between **exabgp** and **BIRD**:
```shell
ipv4 table T_as65500_exaTMS;

protocol bgp as65500_exaTMS {
        description "exaTMS";
        local as 65500;
        neighbor 127.0.0.1 port 1790 as 65500;
        rr client;
        ipv4 {
                table T_as65500_exaTMS;
                import all;
                export none;
        };
}
```
### Reference

+ [Exa-Networks/exabgp](https://github.com/Exa-Networks/exabgp)
+ [alice-lg/birdwatcher](https://github.com/alice-lg/birdwatcher)
+ [czerwonk/bird_socket](https://github.com/czerwonk/bird_socket)
+ [abh/bgpapi](https://github.com/abh/bgpapi)