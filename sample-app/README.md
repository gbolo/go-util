# Sample App Designed for Technical Interviews


1. download kubectl and configure it with provided kubeconfig
2. create two namespaces called dev and prod
3. deploy the sample app to both namespaces, ensure the application (not DB) is HA
4. create an ingress for the application in both namespaces with FQDN:
  - prod-blah.com
  - dev-blah.com
5. There is a list of clients in csv format provided by sales team, Ensure all clients are onboarded
6. modify a client


4. create a deployment of nginx in the prod namespace with 2 replicas named nginx-deployment
. create a service that points to this deployment in the prod namespace
e. exec into nginx-pod1 (dev namespace) and curl the nginx service (in prod namespace) to make sure it works
f. create a network policy that blocks dev nginx-pod1 from accessing this prod nginx service
