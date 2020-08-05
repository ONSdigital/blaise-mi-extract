
# run this from the project directory as scripts/deployExtractFunction

gcloud pubsub topics create ons-blaise-dev-pds-20-extract-topic

gcloud functions deploy ExtractFunction --runtime go113 --trigger-topic ons-blaise-dev-pds-20-extract-topic `
  --env-vars-file scripts/.env.yaml
