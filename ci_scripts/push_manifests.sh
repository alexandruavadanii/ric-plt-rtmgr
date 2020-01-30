source ci_scripts/variables.sh

export DOCKER_CLI_EXPERIMENTAL=enabled

docker manifest create --amend akrainoenea/$IMAGE_NAME:$IMAGE_TAG akrainoenea/$IMAGE_NAME:$IMAGE_TAG-$(uname -m)
docker manifest push --purge akrainoenea/$IMAGE_NAME:$IMAGE_TAG
