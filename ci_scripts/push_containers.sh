# !/bin/bash

source ci_scripts/variables.sh

docker build -t $IMAGE_NAME .
docker tag $IMAGE_NAME:latest akrainoenea/$IMAGE_NAME:$IMAGE_TAG-$AARCH
docker push akrainoenea/$IMAGE_NAME:$IMAGE_TAG-$AARCH

