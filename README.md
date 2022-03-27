# Get Started
### One time setup
Pre-req `gcloud`, `docker`, `terraform`
```
gcloud auth application-default login
gcloud auth configure-docker
```
### Build the image
```
docker build -t gcr.io/PROJECT/NAME .
docker tag gcr.io/PROJECT/NAME gcr.io/PROJECT/NAME:ENV
```
for example I build into my takehome project. Tagging prod just for example, do dev during development
```
docker build -t gcr.io/petal-takehome/reverseserver
docker tag gcr.io/petal-takehome/reverseserver gcr.io/petal-takehome/reverseserver:prod
```
### Run locally
replace image namespace and name as needed
```
docker run -p 8090:8090 gcr.io/petal-takehome/reverseserver
```
### Deploy to gcloud
Note, I build on a M1 chip macbook so cross architecture build is required when deploying to a AMD architecture GCP. You can skip if you have a AMD already
```
docker buildx build --platform linux/amd64 -t gcr.io/petal-takehome/reverseserver .
docker tag gcr.io/petal-takehome/reverseserver gcr.io/petal-takehome/reverseserver:prod
```
Push this image and deploy it to gcp cloud run. This will not work for you because project is configured to my account's project.
Modify that in terraform to your own project or create a new workspace for your project in terraform.
```
docker push gcr.io/petal-takehome/reverseserver:prod
cd terraform
terraform init
terraform apply
```
# Test it out
We find the URL google gave us, in my case it is 
```
https://cloudrun-srv-nwgurm3uia-uc.a.run.app
```
we can hit it like this
```bash
) terraform curl -H 'Content-Type:application/json' -d '{"data":"hello petal"}' https://cloudrun-srv-nwgurm3uia-uc.a.run.app/v1
LATEP OLLEH
```

# Future considerations
Under 3 hour time constraint I'm not able to make too much but here's some considerations
1. monitoring for standard HTTP server latency, throughput, response codes
2. setup monitoring for upstream shoutcloud.io latency and availability
3. use a real HTTP framework for the golang code. I just used the bare minimum `net.http` for simplicity
4. have a real logging system that goes to ELK or google cloud logging. Right now it's minimum `log.Print`
5. cloud cost monitoring, I set my max scale to 3 so my bill stays reasonablely trivial. But in a production environment we'd need monitoring for how much resource we are consuming