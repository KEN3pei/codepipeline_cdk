#!bin/sh

archi=$(uname -m)
if [ $archi = arm64 ]; then
    image="public.ecr.aws/codebuild/local-builds:aarch64"
elif [ $archi = x86_64 ]; then
    image="public.ecr.aws/codebuild/local-builds:latest"
fi
echo image: $image

printf "buildspec projct path?: " & read -r BUILDSPEC_PROJECT_PATH

./codebuild_build.sh -i public.ecr.aws/codebuild/amazonlinux2-x86_64-standard:4.0 -l $image -a test/ -s $BUILDSPEC_PROJECT_PATH -e env.local -p dev-iam-admin -c
