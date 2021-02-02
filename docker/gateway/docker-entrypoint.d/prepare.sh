# 根据环境变量替换swagger.json的地址
sed -i "s|https://petstore.swagger.io/v2/swagger.json|${SWAGGER_JSON}|" /usr/share/nginx/html/swagger-ui/index.html
