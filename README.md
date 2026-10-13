# How to deploy google reverse-proxy on render.
This repo is a simple example used as a starting point to deploy go web on [Render](https://render.com/).

# Deployment

1. Fork [fromxt/google-proxy-render](https://github.com/fromxt/google-proxy-render) on GitHub.

2. Create a new **Web Service** on Render, and give Render permission to access your new repo.

3. Select **GO** for the environment and use the following values during creation:

   > Build Command          **`go build -tags netgo -ldflags '-s -w' -o app`**
 
   > Start Command          **`./app`**

   > Environment           **`Add Environment Variable Key:PROXY_PASSWORD value:your password`**
 4. Your app will be live on your Render URL as soon as the build  finishes.
    
5. Test the Results

     > Once Render finishes redeploying successfully, open your proxy URL in your browser.

     > A native system login prompt will immediately pop up.

     > Username: You can type anything you like (e.g., admin or leave it blank).

     > Password: Enter the PROXY_PASSWORD you set in the Render dashboard (e.g., 123456).

     > Check the browser's "Remember password" box, and the login prompt will never pop up again for future visits.
