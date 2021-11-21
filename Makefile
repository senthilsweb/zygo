.PHONY: assets go_binary build docker
NAME='senthilsweb/zygo'
TAG=$$(git log -1 --pretty=%!H(MISSING))
IMG=${NAME}:${TAG}
LATEST=${NAME}:latest
 
build:
  # @docker build -t ${IMG} .
  # @docker tag ${IMG} ${LATEST}
  echo "${IMG} ${LATEST}"
push:
  @docker push ${NAME}
 
login:
  @docker log -u ${DOCKER_USER} -p ${DOCKER_PASS}
	