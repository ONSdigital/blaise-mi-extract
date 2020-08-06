#!/bin/bash

gcloud pubsub topics create ons-blaise-dev-pds-20-extract-topic

gcloud functions deploy ExtractFunction --runtime go113 --trigger-topic ons-blaise-dev-pds-20-extract-topic \
  --set-env-vars ENCRYPT_LOCATION=ons-blaise-dev-pds-20-mi-encrypt,DB_SERVER='ons-blaise-dev-pds-20:europe-west2:blaise-dev-28475bb5',DB_USER='blaise',DB_PASSWORD='h/REpcUoEPksBt5y',DB_DATABASE='blaise',DB_SOCKET_DIR='/cloudsql'


