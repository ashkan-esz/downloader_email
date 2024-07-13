# Email

email sender service of the [downloader_api](https://github.com/ashkan-esz/downloader_api) project.

## Motivation

making email service a microservice and handle the entire project better

## How to use
docker repository is ashkanaz2828/downloader_email

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

| Prop                                   | Description                                                                              | Required | Default Value |
|----------------------------------------|------------------------------------------------------------------------------------------|----------|---------------|
| **`PORT`**                             | server port                                                                              | `false`  | 3000          |
| **`CORS_ALLOWED_ORIGINS`**             | address joined by `---` example: https://download-admin.com---https:download-website.com | `false`  |               |
| **`RABBITMQ_URL`**                     |                                                                                          | `true`   |               |
| **`INITIAL_WAIT_FOR_MAIL_SERVER_SEC`** |                                                                                          | `false`  | localhost     |
| **`MAIN_SERVER_ADDRESS`**              | the url of the downloader_api (main server)                                              | `true`   |               |
| **`MAILSERVER_HOST`**                  |                                                                                          | `false`  | localhost     |
| **`MAILSERVER_PORT`**                  |                                                                                          | `false`  | 587           |
| **`MAILSERVER_USERNAME`**              |                                                                                          | `false`  |               |
| **`MAILSERVER_PASSWORD`**              |                                                                                          | `false`  |               |
| **`USER_SESSION_PAGE`**                |                                                                                          | `false`  |               |
| **`SENTRY_DNS`**                       | see [sentry.io](https://sentry.io)                                                       | `false`  |               |
| **`PRINT_ERRORS`**                     |                                                                                          | `false`  | false         |


## Local Smtp Server
- mailServer:
    1. https://hub.docker.com/r/boky/postfix
    2. https://docker-mailserver.github.io/docker-mailserver/latest/usage/
    3. create subdomain 'mail' point to server ip. example:: A   mail   SERVER_IP
    4. create record 'MX' with name of domain and point to subdomain. example:: MX   movietracker.site   mail.movietracker.site  DNS only
    5. add rDNS or PTR record to point to domain. example:: PTR   SERVER_IP   movietracker.site  DNS only
    6. add rDNS or PTR in server to point to domain. exmaple:: movietracker.site
    7. add SPF record to dns. example:: TXT   movietracker.site   v=spf1 ip4:SERVER_IP include:movietracker.site +all  DNS only

## Future updates

- [x]  Fast and light.
- [ ]  Documentation.
- [ ]  Write test.

## Contributing

Contributions are always welcome!

See `contributing.md` for ways to get started.

Please adhere to this project's `code of conduct`.

## Support

Contributions, issues, and feature requests are welcome!
Give a ⭐️ if you like this project!

## Related

- [downloader_api](https://github.com/ashkan-esz/downloader_api)
- [downloader_gochat](https://github.com/ashkan-esz/downloader_gochat)

## Author

**Ashkan Esz**

- [Profile](https://github.com/ashkan-esz "Ashkan esz")
- [Email](mailto:ashkanaz2828@gmail.com?subject=Hi "Hi!")
