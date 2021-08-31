#!/bin/bash -e
set -o pipefail

this=`basename $0`

usage () {
cat << EOF
Usage: $this [-h] RELEASE_VERSION

Options:
  -h         show this help and exit
EOF
}

#
# Parse command line
#
while getopts "h" opt; do
    case $opt in
        h)  usage
            exit 0
            ;;
        *)  usage
            exit 1
            ;;
    esac
done
shift "$((OPTIND - 1))"

# Check that no extra args were provided
if [ $# -gt 1 ]; then
    echo -e "ERROR: unknown arguments: $@\n"
    usage
    exit 1
fi

release=$1
container_image=k8s.gcr.io/nfd/node-feature-discovery:$release

#
# Check/parse release number
#
if [ -z "$release" ]; then
    echo -e "ERROR: missing RELEASE_VERSION\n"
    usage
    exit 1
fi

if [[ $release =~ ^(v[0-9]+\.[0-9]+)(\..+)?$ ]]; then
    docs_version=${BASH_REMATCH[1]}
else
    echo -e "ERROR: invalid RELEASE_VERSION '$release'"
    exit 1
fi

if [ -z "$assets_only" ]; then
    # Patch docs configuration
    echo Patching docs/_config.yml
    sed -e s"/release:.*/release: $release/"  \
        -e s"/version:.*/version: $docs_version/" \
        -e s"!container_image:.*!container_image: k8s.gcr.io/nfd/node-feature-discovery:$release!" \
        -i docs/_config.yml

    # Patch README
    echo Patching README.md to refer to $release
    sed s"!node-feature-discovery/v.*/!node-feature-discovery/$release/!" -i README.md

    # Patch deployment templates
    echo Patching kustomize templates to use $container_image
    sed -E -e s",^([[:space:]]+)image:.+$,\1image: $container_image," \
           -e s",^([[:space:]]+)imagePullPolicy:.+$,\1imagePullPolicy: IfNotPresent," \
           -i deployment/base/*/*yaml

    # Patch Helm chart
    echo "Patching Helm chart"
    sed -e s"/appVersion:.*/appVersion: $release/" -i deployment/helm/node-feature-discovery/Chart.yaml
    sed -e s"/pullPolicy:.*/pullPolicy: IfNotPresent/" \
        -e s"!gcr.io/k8s-staging-nfd/node-feature-discovery!k8s.gcr.io/nfd/node-feature-discovery!" \
        -i deployment/helm/node-feature-discovery/values.yaml
    sed -e s"!kubernetes-sigs.github.io/node-feature-discovery/master!kubernetes-sigs.github.io/node-feature-discovery/$docs_version!" \
        -i deployment/helm/node-feature-discovery/README.md

    # Patch e2e test
    echo Patching test/e2e/node_feature_discovery.go flag defaults to k8s.gcr.io/nfd/node-feature-discovery and $release
    sed -e s'!"nfd\.repo",.*,!"nfd.repo", "k8s.gcr.io/nfd/node-feature-discovery",!' \
        -e s"!\"nfd\.tag\",.*,!\"nfd.tag\", \"$release\",!" \
      -i test/e2e/node_feature_discovery.go
fi

#
# Create release assets to be uploaded
#
helm package deployment/helm/node-feature-discovery/ --version $semver

chart_name="node-feature-discovery-chart-$semver.tgz"
mv node-feature-discovery-$semver.tgz $chart_name
sign_helm_chart $chart_name

cat << EOF

*******************************************************************************
*** Please manually upload the following generated files to the Github release
*** page:
***
***   $chart_name
***   $chart_name.prov
***
*******************************************************************************
EOF
