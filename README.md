# letsencrypt-deploy

letsencrypt certificates deployment tool. Certificates can be created with [letsencrypt-lambda](https://github.com/lscheidler/letsencrypt-lambda).

## Requirements

If You want to deploy the client with terraform, You need following tools:

- [ansible](https://www.ansible.com/) (>= 2.9.13)
- [terraform](https://www.terraform.io) (>= 0.13.5)

## Usage

```
Usage of build/linux_amd64/letsencrypt-deploy (0.3):
  -H value
    	run hook after certificates has updated
  -c string
    	set config file (default "/etc/letsencrypt-deploy/config.json")
  -configFile string
    	set config file (default "/etc/letsencrypt-deploy/config.json")
  -d value
    	domains
  -delay int
    	deploy certificates with a delay after creation (days)
  -domain value
    	domains
  -dynamodbTableName string
    	dynamodb table name
  -e string
    	account email
  -email string
    	account email
  -fortios
    	deploy certificate on fortios firewall
  -fortios-base-url string
    	TODO
  -hook value
    	run hook after certificates has updated
  -local
    	deploy certificate on local machine
  -o string
    	output location for certificates (default "/etc/ssl/private")
  -outputLocation string
    	output location for certificates (default "/etc/ssl/private")
  -prefix string
    	file prefix for letsencrypt certificates (default "letsencrypt.")
  -t string
    	dynamodb table name

hooks:

  exec;<command>
        execute <command> after updating certificates
  sns;<sns-topic>[;<sns-subject>[;<sns-message>]]
        publish a auto-generated message to <sns-topic> after updating certificates,
        if sns-message is set, use this message for publishing
```

## Example

```
./letsencrypt-deploy -email me@example.com -domain example.com,*.example.com -passphraseFile /tmp/deploy.passphrase -o certificates/
```

### Production

To delay the deployment of a renewed certificate for 10 days, e.g. for production server (deploying first to testing):

```
./letsencrypt-deploy -delay 10 -email me@example.com -domain example.com,*.example.com -passphraseFile /tmp/deploy.passphrase -o certificates/
```
