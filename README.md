# Automate Kubernetes Debian Base Images

## Dependencies

* Google Container Builder
* [Remote Builder Container](https://github.com/rphillips/cloud-builders-community/tree/k8s_builds/remote-builder)

## Design

This build suite leverages Google Container Builder to perform the following:

* spawn a cloud instance
* provisions the instance
  * latest packages
  * docker
* checks out kubernetes
* runs a build within kubernetes/build/debian-base for all supported architectures
* saves the image to the filesystem
* copies the image back to Container Builder for import into Google Container Repository

## Building the remote-builder image

Note: The destination project will need to be modified.

``` sh
git clone https://github.com/rphillips/cloud-builders-community && cd cloud-builders-community
git checkout k8s_builds
gcloud container builds submit -t gcr.io/rphillips-dev/remote-builder:v0.3.1 . 
```

## Setup the Remote Builder

[Instructions](https://github.com/rphillips/cloud-builders-community/tree/fixes/cleanup_and_doc_fix/remote-builder#quick-start)

## Building debian base images for all architectures

Within this repo (`kubernetes-base-image-builder`) run the following command:

``` sh
gcloud container builds submit --config cloudbuild.yaml --timeout=1h .
```

Resulting containers will be pushed to Google Container Repository.
