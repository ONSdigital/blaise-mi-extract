#!/bin/bash

gcloud functions deploy ExtractFunction --runtime go113 --trigger-resource ons-blaise-dev-pds-20-mi-extract \
  --trigger-event google.storage.object.finalize  --env-vars-file scripts/.env.yaml

