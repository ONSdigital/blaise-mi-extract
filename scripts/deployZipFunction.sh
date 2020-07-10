
gcloud functions deploy ZipFunction --runtime go113 --trigger-resource ons-blaise-dev-pds-18-mi-zip --trigger-event google.storage.object.finalize \
  --set-env-vars "MI_BUCKET_NAME=ons-blaise-dev-pds-18-mi-incoming" \
  --set-env-vars "DEBUG=True" \
  --set-env-vars "FILE_PROVIDER=Google" \
  --set-env-vars "ZIP_LOCATION=ons-blaise-dev-pds-18-mi-zip" \
  --set-env-vars "ENCRYPTED_LOCATION=ons-blaise-dev-pds-18-mi-encrypted" \
  --set-env-vars "ENCRYPT_LOCATION=ons-blaise-dev-pds-18-mi-encrypt"

gcloud alpha functions add-iam-policy-binding ZipFunction --member=allUsers --role=roles/cloudfunctions.invoker