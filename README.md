# lavazares
`lavazares` is the backend for [typephil](https://github.com/codephil-columbia/typephil).

## Local environment setup
1. `git pull origin master`
2. Ensure you have all `go` dependencies using `dep ensure`. [dep](https://github.com/golang/dep) is the best way to get them. 
3. Make sure you have the secrets.json file in the projects root directory. 
4. In order to run locally on the prod db (you probably shouldn't do this), run `go run main.go -local=false`
5. In order to run locally with a local db, set it up using the instructions below, and then run `go run main.go`


### Running Local DB setup 
If you are running the local version of the db, you're going to need to Docker and a python version > 2.7. Then
1. Create a virtualenviroment (will probably be `virtualenv env`)
2. Start virtualenv `source env/bin/activate`
3. Install python dependencies `pip install -r requirements.txt`
4. Run the build script `./build-local`

The build script starts a local copy of the database using the schema defined in `/docker/local/entities.sql` and prepopulates
it with the values in `/docker/local/test_data.sql` running at port 5432. Feel free to do whatever you want with the db. 
If you want to clean it and populate it with the data it started with, rerun the build script. 

## Running on AWS 
1. Rebuild the app. `go build .`
2. Kill the previous docker container. You can find the running containers using `sudo docker ps`. Copy the container ID of the container and run `sudo docker stop <container-id>`. 
3. Build new container. Run `sudo docker build .` This will output the ID of the newly built container, copy that since you'll need it in order to run it.
4. Run container. `sudo docker run -d -p 8081:5000 <id-you-previously-copied>` The `-d` flag runs the container in the background, and `-p` the open port on AWS to the port that the app listens on.
