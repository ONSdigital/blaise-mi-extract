
# run this from the project directory as scripts/deployZipFunction

gcloud functions deploy ZipFunction --runtime go113 `
  --trigger-resource ons-blaise-dev-pds-20-mi-zip `
  --region=europe-west2 --trigger-event google.storage.object.finalize `
  --set-env-vars ENCRYPT_LOCATION=ons-blaise-dev-pds-20-mi-encrypt

