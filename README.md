# lavazares
`lavazares` is the backend for [typephil](https://github.com/codephil-columbia/typephil).

## Local environment setup
1. `git pull origin master`
2. Ensure you have all `go` dependencies using `dep ensure`. [dep](https://github.com/golang/dep) is the best way to get them. 
3. Make sure you have the secrets.json file in the projects root direcotry. 
4. In order to run locally on the prod db (rip), run `go run main.go -local=false`

## Running on AWS 
1. Rebuild the app. `go build .`
2. Kill the previous docker container. You can find the running containers using `sudo docker ps`. Copy the container ID of the container and run `sudo docker stop <container-id>`. 
3. Build new container. Run `sudo docker build .` This will output the ID of the newly built container, copy that since you'll need it in order to run it.
4. Run container. `sudo docker run -d -p 8081:5000 <id-you-previously-copied>` The `-d` flag runs the container in the background, and `-p` the open port on AWS to the port that the app listens on.
