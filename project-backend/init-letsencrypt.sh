#!/bin/bash

domains=(azdent.alitech.biz)
rsa_key_size=4096
data_path="./certbot"
email="isayev.a@alitech.biz" # ← замени на свой рабочий email
staging=0

if [ -d "$data_path" ]; then
  read -p "⚠️ SSL уже существует. Заменить? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi

mkdir -p "$data_path/www"
mkdir -p "$data_path/conf"

docker compose up -d nginx

sleep 5

domain_args=""
for domain in "${domains[@]}"; do
  domain_args="$domain_args -d $domain"
done

email_arg="--email $email"
staging_arg=""
if [ $staging != "0" ]; then
  staging_arg="--staging"
fi

docker compose run --rm certbot certonly --webroot -w /var/www/certbot \
  $staging_arg \
  $email_arg \
  $domain_args \
  --rsa-key-size $rsa_key_size \
  --agree-tos \
  --force-renewal

docker compose down
