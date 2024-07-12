# Export Your Subscriptions with Google Takeout

## Go to Google Takeout:
    Visit Google Takeout.

## Select YouTube Data:
    - Deselect all products.
    - Scroll down and select "YouTube and YouTube Music".
    - Click on "All YouTube data included" and make sure only "subscriptions" is checked.
    - Click "OK".

## Choose Export Options:
    - Click "Next step".
    - Choose your export method, file type, and delivery method.
    - Click "Create export".

## Download the Exported Data:
    - Google will prepare the export and notify you when it’s ready.
    - Download the export file (usually a ZIP file).

## Extract the Exported Data:
    - Extract the contents of the ZIP file.
    - Inside, you’ll find a CSV file containing your subscriptions list.

## Copy subscriptions.csv to Project Repo
    - Copy subscriptions.csv to root project

# Export client_secret.json Step-by-Step Guide:

## Create a New Project (if you haven't already):
    1. Go to the Google Cloud Console.
    2. Click on the project dropdown at the top of the page and select "New Project".
    3. Give your project a name and click "Create".

## Enable the YouTube Data API:
    1. In the Google Cloud Console, navigate to the "APIs & Services" dashboard.
    2. Click on "Enable APIs and Services".
    3. Search for "YouTube Data API v3" and select it.
    4. Click "Enable" to enable the API for your project.

## Create OAuth 2.0 Credentials:
    1. Go to the "APIs & Services" > "Credentials" page.
    2. Click on "Create Credentials" and select "OAuth client ID".
    3. If prompted to set up the OAuth consent screen, follow the instructions to configure it. You will need to provide some basic information about your application. Input "youtube.force-ssl" & "youtube.read-only" to scope
    4. For the application type, select "Desktop app".
    5. Give your OAuth client ID a name (e.g., "YouTube API Client").
    6. Click "Create".

## Download the client_secret.json File:
    1. After creating the OAuth client ID, you will see a dialog with your client ID and client secret. Click "Download JSON".
    2. This will download the client_secret.json file to your computer. Save this file in a secure location.

## Move the client_secret.json File to Your Working Directory:
    1. Move the downloaded client_secret.json file to the directory where you will be running your Python script. Ensure that the file name matches what you reference in your script.
    2. Rename it with client_secret.json

## Build Project & Run

Read Here:
https://go.dev/doc/tutorial/compile-install

For Development:
$ go run import_subscriptions.go

If Prompted, Follow instruction to get token
