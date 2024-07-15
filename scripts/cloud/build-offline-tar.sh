#!/bin/bash
set -x

ARCH=${ARCH:-"amd64"}
CLOUD_VERSION=${CLOUD_VERSION:-"latest"}

# pull and save images
mkdir -p output/tars

images=(
  docker.io/nebstudio/sealos-cloud:$CLOUD_VERSION
  docker.io/nebstudio/kubernetes:v1.27.11
  docker.io/nebstudio/helm:v3.14.1
  docker.io/nebstudio/cilium:v1.14.8
  docker.io/nebstudio/cert-manager:v1.14.6
  docker.io/nebstudio/openebs:v3.10.0
  docker.io/nebstudio/victoria-metrics-k8s-stack:v1.96.0
  docker.io/nebstudio/ingress-nginx:v1.9.4
  docker.io/nebstudio/kubeblocks:v0.8.2
  docker.io/nebstudio/kubeblocks-redis:v0.8.2
  docker.io/nebstudio/kubeblocks-mongodb:v0.8.2
  docker.io/nebstudio/kubeblocks-postgresql:v0.8.2
  docker.io/nebstudio/kubeblocks-apecloud-mysql:v0.8.2
  docker.io/nebstudio/cockroach:v2.12.0
  docker.io/nebstudio/metrics-server:v0.6.4
)

for image in "${images[@]}"; do
  sealos pull --platform "linux/$ARCH" "$image"
  filename=$(echo "$image" | cut -d':' -f1 | tr / -)
  if [[ ! -f "output/tars/${filename}.tar" ]]; then
    sealos save -o "output/tars/${filename}.tar" "$image"
  fi
done


# get and save cli
mkdir -p output/cli

VERSION="v5.0.0"

wget https://github.com/nebstudio/sealos/releases/download/${VERSION}/sealos_${VERSION#v}_linux_${ARCH}.tar.gz \
   && tar zxvf sealos_${VERSION#v}_linux_${ARCH}.tar.gz sealos && chmod +x sealos && mv sealos output/cli

# get and save install scripts
echo "
#!/bin/bash
bash scripts/load-images.sh
bash scripts/install.sh --cloud-version=$CLOUD_VERSION

" > output/install.sh

mkdir -p output/scripts

echo '
#!/bin/bash

cp cli/sealos /usr/local/bin

for file in tars/*.tar; do
  /usr/local/bin/sealos load -i $file
done

'  > output/scripts/load-images.sh

curl -sfL https://raw.githubusercontent.com/nebstudio/sealos/${CLOUD_VERSION}/scripts/cloud/install.sh -o output/scripts/install.sh

# tar output to a tar.gz
mv output sealos-cloud
tar czfv sealos-cloud.tar.gz sealos-cloud

# md5sum output tar.gz
md5sum sealos-cloud.tar.gz | cut -d " " -f1 > sealos-cloud.tar.gz.md5
