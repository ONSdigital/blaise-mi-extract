#!/bin/bash

gcloud functions deploy ZipFunction --runtime go113 --trigger-resource ons-blaise-dev-pds-20-mi-zip \
  --trigger-event google.storage.object.finalize  --env-vars-file scripts/.env.yaml


