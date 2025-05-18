# shelf | HTTPS_Hub

### A proxy server that automates HTTPS connections to local IP addresses

## Quickstart

1. **Clone the repository:**

```shell
git clone https://github.com/projects-shelf/HTTPS_Hub.git
cd HTTPS_Hub
```

2. **Create a Cloudflare API Token:**

Generate an "Edit zone DNS" API token at Cloudflare (ensure the token includes root domain access).

3. **Edit Environment Variables:**
   
Edit environment variables at `docker-compose.yml`

4. **Configure Routing:**

Edit `config.yml` to define the routing rules.

Example configuration:

```yml
ports:
  - rssr=:50046
  - rssg=192.168.0.2:50048
```

This will route:

- `https://rssr.local.example.com` → `http://{LOCAL_IP}:50046`
- `https://rssg.local.example.com` → `http://192.168.0.2:50048`

5. **Run the Server:**

RUN ```docker-compose up -d```

## Features

### Custom Routing

Easily configure domain-to-IP mappings through config.yml.

### HTTPS Proxying

Converts HTTP requests to HTTPS and handles SSL/TLS certificates automatically.

### Cloudflare DNS Management

Automatically updates DNS records using Cloudflare's API.

### Let's Encrypt Support

Automatic SSL certificate generation and renewal using Let's Encrypt.


## License

HTTPS_Hub is licensed under [MIT License](https://github.com/projects-shelf/HTTPS_Hub/blob/main/LICENSE).

## Author

Developed by [PepperCat](https://github.com/PepperCat-YamanekoVillage).
