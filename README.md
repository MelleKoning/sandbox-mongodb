# This is a sandbox project to have a look at a mongodb package

Origin: https://github.com/mongodb/mongo-go-driver

# Adding the mongo-go-driver

following lines have added the mongo-go-driver to the package.

dep init
dep ensure -add github.com/mongodb/mongo-go-driver/mongo

# above lines will pull the driver into the vendor folder

# Run the driver build in tests

First ensure you have a working version of MongoDB running. This can be a mongo database that is running in a docker instance.

For a full setup of mongo including a Primary and a replica set, 
we can follow the installation instructions:

https://www.sohamkamani.com/blog/2016/06/30/docker-mongo-replica-set/
this needed to include the first comment at the bottom of the post, as the instructions
in the post itself did not even work :(

If you do not like reading the post, 
ok, the docker commands to set this up are copied here:

setup the network:
* docker network create my-mongo-cluster
  
setup the mongo1 instance 
this is without detached mode so you will have to use
new console for each to run

* exposed ports need to be different, just as internal ports
docker run -d --net my-mongo-cluster -p 27017:27017 --name mongo1 mongo mongod --replSet my-mongo-set --port 27017
docker run -d --net my-mongo-cluster -p 27018:27018 --name mongo2 mongo mongod --replSet my-mongo-set --port 27018
docker run -d --net my-mongo-cluster -p 27019:27019 --name mongo3 mongo mongod --replSet my-mongo-set --port 27019

# setup the config by attaching to the mongo instance
example, attach to mongo within docker container:
# enter mongo1 mongo shell
docker exec -it mongo1 mongo
* execute following commands to create a test database
* and setup the replica set configuration
>db = (new Mongo('localhost:27017')).getDB('test')
>config={ "_id" : "my-mongo-set", "members" : [ { "_id" : 0, "host" : "mongo1:27017" }, { "_id":1,"host":"mongo2:27018"},{"_id":2,"host":"mongo3:27019"}]}
* following tells the test-database about the three different instances of mongo1, mongo2, mongo3
 >rs.initiate(config) 
* show the detailed configuration
>rs.status()

* on windows add to hostfile to recognize the mongo1, mongo2, mongo3 as domain names:
127.0.0.1 mongo1 
127.0.0.1 mongo2
127.0.0.1 mongo3

* test connection on bash command prompt:
>mongo mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set

* Test connection when mongo1 is down:
$ docker stop mongo1
$ mongo mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set

or test via browser:
http://localhost:27017/_replSet

# Goal of this
The goal was to see if we can follow changes in a replica set.
So, an external system can do a change to the collection, and the 'WATCH' command
sees this change occuring. This can be used to act upon collection changes.
see main_test.go for the implementation that watches a cursor/collection for any changes.
Put a breakpoint on the for loop that checks for the cursor, and from
another client (like Robo3T) do an update to the collection, you will 
see the cursor sees that external change using the changestream
