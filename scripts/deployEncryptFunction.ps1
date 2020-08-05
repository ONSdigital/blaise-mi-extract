
# run this from the project directory as scripts/deployExtractFunction

gcloud functions deploy EncryptFunction --runtime go113 --trigger-resource ons-blaise-dev-pds-20-mi-encrypt `
  --trigger-event google.storage.object.finalize --env-vars-file scripts/.env.yaml
