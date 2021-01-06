### 利用acme 生成https协议

##### 手动https
1.curl  https://get.acme.sh | sh  
2.cd ~/.acme.sh/  
3.alias acme.sh=~/.acme.sh/acme.sh  
4.acme.sh --issue --force  -d a.test.com --webroot /var/www/html/a/  
5.mkdir -p /etc/nginx/ssl_cert/a.test.com  
6.acme.sh --install-cert -d a.test.com --key-file /etc/nginx/ssl_cert/a.test.com/key.pem --fullchain-file /etc/nginx/ssl_cert/a.test.com/cert.pem --reloadcmd     "service nginx force-reload"

##### 手动dns
1.acme.sh --issue -d a.test.com --dns --yes-I-know-dns-manual-mode-enough-go-ahead-please  
2.nslookup -q=TXT _acme-challenge.a.test.com 
3.acme.sh --renew -d a.test.com --yes-I-know-dns-manual-mode-enough-go-ahead-please  
4.mkdir -p /etc/nginx/ssl_cert/a.test.com 
5.acme.sh --install-cert -d a.test.com --key-file /etc/nginx/ssl_cert/a.test.com/key.pem --fullchain-file /etc/nginx/ssl_cert/a.test.com/cert.pem --reloadcmd     "service nginx force-reload"  
###### Nginx 配置：
6.listen 443;  
  ssl_certificate /etc/nginx/ssl_cert/a.test.com/cert.pem;  
  ssl_certificate_key /etc/nginx/ssl_cert/a.test.com/key.pem;  
  ssl_protocols TLSv1 TLSv1.1 TLSv1.2;  
  ssl_prefer_server_ciphers on;  