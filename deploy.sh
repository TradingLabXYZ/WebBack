#!/bin/sh
gcloud config set project complimenti
gcloud app deploy app.yaml --stop-previous-version
