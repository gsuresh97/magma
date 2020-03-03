# Introduction
This README serves as a guide to build and install binaries for a robust 
version of magma in the magma dev VM. There are also instructions on how to 
run integration tests in order to ensure that magma is working the way it 
should. Finally, there are some pointers on what to pay attention to when 
making the rest of the Attach procedure in Magma robust.

# Setup Steps
1. Install Magma
2. Install Zookeeper
3. Configure Zookeeper
4. Building with Zookeeper
5. Installing Magma Test Environment
6. Run Tests

# Setup
## Install Magma
1. Follow the steps under the sections "Development Tools" and "Build/Deploy Tooling" [here](https://facebookincubator.github.io/magma/docs/basics/prerequisites) to download all the prerequisites.
2. Follow the steps under the sections "Provisioning the Environment" and "Initial Run" [here](https://facebookincubator.github.io/magma/docs/basics/quick_start_guide) to build magma and ensure that the service starts.

## Install Zookeeper
1. Run `vagrant up magma && vagrant ssh magma` to `ssh` into the VM being used to run the magma service.
2. Zookeeper will have to be installed locally in this VM so that the gateway can use the Zookeeper client binaries to communicate with the Zookeeper server. 
3. Install the JDK on the VM using `sudo apt-get install default-jdk`
4. Download the apache zookeeper files from [here](https://downloads.apache.org/zookeeper/stable/)

        # Replace X with the appropriate version number
        wget https://downloads.apache.org/zookeeper/stable/apache-zookeeper-X.X.X.tar.gz

    Do not install the `.*-bin.tar.gz` files, they do not come with the necessary header files for C clients
5. If you're installing the Zookeeper server in a different machine than the magma VM, follow the instructions [here](https://linuxconfig.org/how-to-install-and-configure-zookeeper-in-ubuntu-18-04).
6. Unarchive the tar file into `/opt/`
7. `cd` into `/opt/apache-zookeeper-[version]` and run `mvn clean install` to compile the binaries
8. To build the C client library and generate the docs for the library, follow steps 2-6 under INSTALLATION [here](https://github.com/apache/zookeeper/blob/master/zookeeper-client/zookeeper-client-c/README). To generate the docs, the packages `graphviz` and `doxygen` may also be needed. A copy of the C client API for zookeeper (though not formatted as nicely as the generated docs) are available [here](https://zookeeper.dpldocs.info/deimos.zookeeper.zookeeper.html).

## Configuring Zookeeper
1. Create a directory for the data to be stored at `/var/lib/zookeeper`. Then, create a Zookeeper config file at `/opt/apache-zookeeper-[version]/conf/zoo.cfg` and set it to the following:

        # The number of milliseconds of each tick
        tickTime=2000
        # The number of ticks that the initial 
        # synchronization phase can take
        initLimit=10
        # The number of ticks that can pass between 
        # sending a request and getting an acknowledgement
        syncLimit=5
        # the directory where the snapshot is stored.
        # do not use /tmp for storage, /tmp here is just 
        # example sakes.
        dataDir=/var/lib/zookeeper
        # the port at which the clients will connect
        clientPort=2181
        # the maximum number of client connections.
        # increase this if you need to handle more clients
        maxClientCnxns=60
        #
        # Be sure to read the maintenance section of the 
        # administrator guide before turning on autopurge.
        #
        # http://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_maintenance
        #
        # The number of snapshots to retain in dataDir
        #autopurge.snapRetainCount=3
        # Purge task interval in hours
        # Set to "0" to disable auto purge feature
        #autopurge.purgeInterval=1
        4lw.commands.whitelist=*
2. Start the Zookeeper server locally in magma VM by executing `sudo /opt/apache-zookeeper-[version]/bin/zkServer.sh start`
3. To test that the server is running, exeute the following:

        $ sudo /opt/apache-zookeeper[version]/bin/zkCli.sh
        > create /test hello
        > get /test 
    This should print "hello"

        > deletall /test
4. To reset the data in Zookeeper, you can either call the `deleteall` in the Zookeeper client shell on all the paths or you can delete the `dataDir` set in `zoo.cfg` In this case, that's `/var/lib/zookeeper`.

## Building with Zookeeper.
1. To link successfully with Zookeeper while building magma, ensure that the path of the zookeeper library included in line 250 of `lte/gateway/c/oai/tasks/nas/CMakeLists.txt` exists and is valid. Currently it is `/opt/apache-zookeeper-3.5.6/zookeeper-client/zookeeper-client-c/.libs/`. Change the version number as needed and make sure this path exists.
2. Build magme as before, by running `make run` in the magma root directory.

## Installing the Magma Test Environment
1. Follow the steps under the "Test VM setup" section [here](https://github.com/facebookincubator/magma/blob/master/docs/readmes/lte/s1ap_tests.md) This will download and install another VM that emulates the UE and communicates with the original magma dev VM.

## Running Tests
1. Make sure the magma service is running in the magma dev VM by following the instructions given [here](https://github.com/facebookincubator/magma/blob/master/docs/readmes/lte/s1ap_tests.md) under the "Gateway VM Setup" section.
2. The actual LTE procedures will be running in the magma dev VM. The testing VM only issues the calls. 
3. The magma logs are located at `/var/log/mme.log`. The functions being called as well as the error messages are logged here quite clearly. The logs for the attach modifications made in this repo are also here.
4. To run the tests, follow the instructions under the "Run Tests" section [here](https://github.com/facebookincubator/magma/blob/master/docs/readmes/lte/s1ap_tests.md). 
5. The reliable attach test is in `s1aptests/reliable_attach.py`. 
6. To check that the reliably implemented portion of attach is running properly, run `make integ-_test TESTS=s1aptests/reliable_attach.py` from `$MAGMA_ROOT/lte/gateway/python/integ_tests` on the test VM. This will log the contents of the hashtable that is stored in the log file prefixed by "Hashtable Dump". 
7. Then, run `make restart` in `magma/lte/gateway` of the dev VM to restart the magma service and clear the state.
8. Run step 6 again, to see "Fetching hastable from Zookeeper success" in the logs, followed by a dump of the hash table that was fetched prefixed again by "Hashtable Dump"

