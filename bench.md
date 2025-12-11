```bash
echo '{"url":"https://pogoda.yandex.ru"}' > /tmp/url.json
ab -k -c 10 -n 20000 http://localhost:8080/api/url/ &                                                
ab -k -c 10 -n 20000 http://localhost:8080/api/url/5001/history &                                                
ab -k -c 10 -n 20000 http://localhost:8080/api/url/5002/history &                                                
ab -k -c 10 -n 20000 http://localhost:8080/api/url/5003/history &                                                
ab -k -c 10 -n 20000 -p /tmp/url.json -T 'application/json' http://localhost:8080/api/url &        
ab -k -c 10 -n 20000 -p /tmp/url.json -T 'application/json' http://localhost:8080/api/url/activate  &                                         
ab -k -c 10 -n 20000 -p /tmp/url.json -T 'application/json' http://localhost:8080/api/url/deactivate & 

wait 
```
