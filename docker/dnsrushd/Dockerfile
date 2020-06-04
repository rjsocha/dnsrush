FROM php:apache
RUN a2dismod status && a2enmod rewrite
COPY cnf/php.ini /usr/local/etc/php/php.ini
COPY app/ /var/www/html/
COPY scripts/ /usr/local/bin/
RUN chmod +x /usr/local/bin/* && chown -R www-data:www-data /var/www/html/deploy
