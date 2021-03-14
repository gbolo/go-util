API_ENDPOINT="http://127.0.0.1:60081"

# ENSURE ALL CLIENTS ARE DELETED
for id in 1 2 3; do
  http -v DELETE ${API_ENDPOINT}/api/v1/client/${id}
done

# ADD CLIENTS
http -v POST ${API_ENDPOINT}/api/v1/client id=1 name=linuxctl url=https://linuxctl.com/ip
http -v POST ${API_ENDPOINT}/api/v1/client id=2 name=google url=https://google.com
http -v POST ${API_ENDPOINT}/api/v1/client id=3 name=reddit url=https://reddit.com

# RETRIEVE CLIENTS
http -v GET ${API_ENDPOINT}/api/v1/client

# UPDATE CLIENT
http -v PUT ${API_ENDPOINT}/api/v1/client id=1 name=linuxctl url=https://linuxctl.com/404

# RETRIEVE CLIENTS
http -v GET ${API_ENDPOINT}/api/v1/client

# DELETE ALL CLIENTS
for id in 1 2 3; do
  http -v DELETE ${API_ENDPOINT}/api/v1/client/${id}
done

# RETRIEVE CLIENTS
http -v GET ${API_ENDPOINT}/api/v1/client