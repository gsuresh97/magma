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

# Attach Notes
* The call graph located [here](https://github.com/gsuresh97/magma/blob/master/docs/readmes/lte/Attach_call_flow_in_Magma.png) provides a lot of detail. All the lines in green color are places where the context is added and the lines in red signify areas where the context is changed. Since the context is only synced with zookeeper before and after it's used, some of these can be ignored but they should all be checked nevertheless. 
* Contexts used in multiple places during the attach procedure call include  `mme_ue_context_t` in `mme_app_ue_context.h` and `emm_context_s` in `emm_data.h`. 
  * These contexts contain mostly hashtables and other pointers structs.
  * To effectively store this in zookeeper, the pointers will need to be dereferenced and serialized seperately. There is a chance those may also consist of pointers so in order to serialize these contexts efficiently, a profile of these contexts could first be constructed and the serialization could be dome recrusively using the profile as a guideline.
* A context is also stored in a few linked lists such as the one in `emm_context->emm_procedures->emm_common_procs` in the struct `emm_procedures_t` in `lte/gateway/c/oai/tasks/nas/nas_procedures.h` These will need to be serialized seperately, since the nodes are represented by pointers that cannot be serialized using `memcpy()`.
* There are also some instances where it seems that magma creates and uses containers that store callback functions that are stored as function pointers across multiple calls to the Attach procedure. These will also need to be handled seperately. 
  * One option may be to create a constant global hashtable that encodes known values to function pointers, and serialize the key corresponding to the function pointer when the context is stored in Zookeeper. 
  * This, while making the implementation simpler could be messy to maintain. Another approach could be to refactor the code and design it such that function pointers need not be stored in Zookeeper.
  * One example of function pointers being stored is in `lte/gateway/c/oai/tasks/nas/nas_procedures.c` in the function `nas_new_authentication_procedure`.
  * This function stores a wrapper of type `nas_emm_common_procedure_t` to a list. If we drill down far enough, we see that `wrapper->proc.emm_proc.base_proc.success_notif` has the folowing typdef: `typedef int (*success_cb_t)(struct emm_context_s *);`