FROM stephenmiracle/frankenwp:beta-1.0.0


ENV XDEBUG_REMOTE_ENABLE=1

# Access the phpinfo to show xdebug is installed.
COPY ./phpinfo.php /var/www/html/phpinfo.php
COPY ./test.php /var/www/html/test.php
COPY ./dev.ini $PHP_INI_DIR/conf.d/dev.ini

# install xdebug
RUN pecl install xdebug && docker-php-ext-enable xdebug




