# A GCS Event Trigger to Change TSV File to CSV File

This Cloud Functions is for when user upload a `.tsv` file, will trigger function to change from `.tsv` to `.csv`.

## Deploy

```
$ gcloud config set project <YOUR_GCP_PROJECT_ID>
```

```
$ gcloud functions deploy GcsTSVtoCSV --runtime go111 --trigger-resource <YOUR_GCS_BUCKET_NAME> --trigger-event google.storage.object.finalize
```

## Use
Upload a `.tsv` file to this bucket, and wait some times, then you will get a `.csv` file.
