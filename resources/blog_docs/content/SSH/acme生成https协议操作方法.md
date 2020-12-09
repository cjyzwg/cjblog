### 利用acme 生成https协议
1.curl  https://get.acme.sh | sh  
2.cd ~/.acme.sh/  
3.alias acme.sh=~/.acme.sh/acme.sh  
4.acme.sh --issue --force  -d a.test.com --webroot /var/www/html/a/  
5.mkdir -p /etc/nginx/ssl_cert/a.test.com  
6.acme.sh --install-cert -d a.test.com --key-file /etc/nginx/ssl_cert/a.test.com/key.pem --fullchain-file /etc/nginx/ssl_cert/a.test.com/cert.pem --reloadcmd     "service nginx force-reload"