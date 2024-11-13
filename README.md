#BUILD 

docker buildx build --no-cache -o type=docker -f production.Dockerfile -t huecomundo.jeremysaputra.com:5000/miniwebserver:prod .
docker push huecomundo.jeremysaputra.com:5000/miniwebserver:prod