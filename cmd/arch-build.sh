#!/bin/bash
name=$1
check=$(./check-uf.sh $name)
if [ "$check" != "$name" ]; then
	echo $check
	exit 1
fi
cd $name
port=$(cat UF_PORT)

export CGO_ENABLED=0 
export GOOS=linux 
GOARCH=amd64 go build -ldflags="-s -w" -o build/"$name"_amd64
GOARCH=arm64 go build -ldflags="-s -w" -o build/"$name"_arm64
chmod +x build/"$name"_amd64
chmod +x build/"$name"_arm64

#upx --best --lzma build/amd64/"$name"
#upx --best --lzma build/arm64/"$name"

base="localhost:5000/$name"
podman manifest rm $base

buildah manifest create $base

for arch in arm64 amd64; do
	current="${base}_${arch}:latest"
	bin_name=${name}_${arch}

	echo $base
	echo $current
	echo $bin_name
	podman rmi -f $current
	podman rm -f mycontainer

	#buildah from --name mycontainer mirror.gcr.io/library/alpine
	buildah from --name mycontainer mirror.gcr.io/library/alpine
	echo ==============================================
	buildah config --workingdir /bin mycontainer
	buildah copy mycontainer build/$bin_name "/bin/$bin_name"
	buildah copy mycontainer /fridayy/taxes/full.csv /fridayy/taxes/full.csv
	#buildah copy mycontainer ../../localhost+8-key.pem "/bin/"
	#buildah copy mycontainer ../../localhost+8.pem "/bin/"
	#buildah copy mycontainer ../../.env.prod "/bin/.env"

	buildah config --entrypoint "/bin/$bin_name" mycontainer
	#buildah config --created-by "textFridayy" $latest
	#buildah config --author "Taylor" $latest
	##buildah commit $latest $latest

	buildah commit --squash mycontainer $current
	buildah manifest add --arch $arch $base $current
done

buildah manifest push $base docker://$base --all

echo "built service: $base"
