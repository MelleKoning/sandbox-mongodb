# This is a sandbox project to have a look at a mongo-go-driver  package

Origin: https://github.com/mongodb/mongo-go-driver

# Why have a look at this package?

The mgo driver we currently are using in Veritas is not being updated for some time.
Also, that mgo driver could not be used for 'watching' changes on a mongo-replica set.

The goal was to see if with this package:
- we can setup a replica set using docker and then
- see if we can follow changes to a collection using a mongodb-cursor

So, what happens is, an external system or one of the running services might do a change to
a document in the collection. Can your service see this change and act upon it?
The changestream-'WATCH' command sees this change occuring. This can be used to act
upon collection changes.

# Adding the mongo-go-driver

if you first get this project you need to get the driver in your vendor folder. Following lines have added the mongo-go-driver to the package to vendor folder.

* dep init 
* dep ensure -add github.com/mongodb/mongo-go-driver/mongo

Above lines will pull the driver into the vendor folder

# Run the driver build-in tests

First ensure you have a working instance of MongoDB running. This can be a mongo database that is running in a docker instance.

* Replica set
  
For a full setup of mongo including a Primary and a replica set, 
we can follow the installation instructions:

https://www.sohamkamani.com/blog/2016/06/30/docker-mongo-replica-set/
this needed to include the first comment at the bottom of the post, as the instructions
in the post itself did not even work :(

If you do not like reading the post, take steps in the next chapter

# Setup the replica set

ok, the docker commands to set this up are copied here:

setup the network:
* docker network create my-mongo-cluster
  
setup the mongo1 instance 
this is without detached mode so you will have to use
new console for each to run

exposed ports need to be different, just as internal ports

* docker run -d --net my-mongo-cluster -p 27017:27017 --name mongo1 mongo mongod --replSet my-mongo-set --port 27017
* docker run -d --net my-mongo-cluster -p 27018:27018 --name mongo2 mongo mongod --replSet my-mongo-set --port 27018
* docker run -d --net my-mongo-cluster -p 27019:27019 --name mongo3 mongo mongod --replSet my-mongo-set --port 27019

setup the config by attaching to the mongo instance
example, attach to mongo within docker container:

Enter the mongo1 mongo shell
* docker exec -it mongo1 mongo

execute following commands to create a test database
and setup the replica set configuration

>db = (new Mongo('localhost:27017')).getDB('test')

>config={ "_id" : "my-mongo-set", "members" : [ { "_id" : 0, "host" : "mongo1:27017" }, { "_id":1,"host":"mongo2:27018"},{"_id":2,"host":"mongo3:27019"}]}

following tells the test-database about the three different instances of mongo1, mongo2, mongo3

 >rs.initiate(config) 

show the detailed configuration

>rs.status()

to exit out of the shell type
* exit
  
# windows hostfile to see mongo instances
on windows add to hostfile to recognize the mongo1, mongo2, mongo3 as domain names:

* 127.0.0.1 mongo1 
* 127.0.0.1 mongo2
* 127.0.0.1 mongo3

test connection on bash command prompt, type:

>mongo mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set

Test connection when mongo1 is down:
* docker stop mongo1
* mongo mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=my-mongo-set

or see reply via browser:
http://localhost:27017/_replSet

# Conclusion: See result of this

See main_test.go for the (simplified) implementation that watches a
cursor/collection for any changes.

While debugging put a breakpoint on the for loop that checks for the cursor, and from
another client (like Robo3T) do an update to the collection, you will 
see the Watch() method cursor does see that external change using the changestream, so that we can act upon that, for example update a local read-cache.
