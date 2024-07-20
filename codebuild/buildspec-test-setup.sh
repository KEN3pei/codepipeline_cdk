#!bin/sh

docker pull public.ecr.aws/codebuild/amazonlinux2-x86_64-standard:4.0 | tee logs/dockerlog.txt

archi=$(uname -m)
if [ $archi = arm64 ]; then
    docker pull public.ecr.aws/codebuild/local-builds:aarch64 | tee -a logs/dockerlog.txt
elif [ $archi = x86_64 ]; then
    docker pull public.ecr.aws/codebuild/local-builds:latest | tee -a logs/dockerlog.txt
fi

curl -O  https://raw.githubusercontent.com/aws/aws-codebuild-docker-images/master/local_builds/codebuild_build.sh | tee -a logs/dockerlog.txt
chmod +x codebuild_build.sh | tee -a logs/dockerlog.txt

