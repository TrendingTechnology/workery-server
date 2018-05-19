#!/bin/bash
#
# setup.env.sh
# The purpose of this script is to setup sample environment variables for this project. It is up to the responsibility of the developer to change these values once they are generated for production use.
#

# Step 1: Clear the file.
clear;
cat > workery/workery/.env << EOL
#--------#
# Django #
#--------#
SECRET_KEY=l7y)rwm2(@nye)rloo0=ugdxgqsywkiv&#20dqugj76w)s!!ns
DEBUG=True
ALLOWED_HOSTS='*'
ADMIN_NAME='Bartlomiej Mika'
ADMIN_EMAIL=bart@mikasoftware.com

#----------#
# Database #
#----------#
DATABASE_URL=postgis://django:123password@localhost:5432/workery_db
DB_NAME=workery_db
DB_USER=django
DB_PASSWORD=123password
DB_HOST=localhost
DB_PORT="5432"

#-------#
# Email #
#-------#
DEFAULT_TO_EMAIL=bart@mikasoftware.com
DEFAULT_FROM_EMAIL=postmaster@mover55london.ca
EMAIL_BACKEND=django.core.mail.backends.console.EmailBackend
MAILGUN_ACCESS_KEY=<YOU_NEED_TO_PROVIDE>
MAILGUN_SERVER_NAME=over55london.ca

#----------------#
# Django-Htmlmin #
#----------------#
HTML_MINIFY=True
KEEP_COMMENTS_ON_MINIFYING=False

#--------#
# Sentry #
#--------#
SENTRY_RAVEN_CONFIG_DSN=https://xxxx:yyyyy@sentry.io/zzzzzzzz

#--------------------------------#
# Application Specific Variables #
#--------------------------------#
O55_LOGLEVEL=INFO
O55_APP_HTTP_PROTOCOL=http://
O55_APP_HTTP_DOMAIN=workery.ca
O55_APP_DEFAULT_MONEY_CURRENCY=CAD
O55_GITHUB_WEBHOOK_SECRET=None
EOL

# Developers Note:
# (1) Useful article about setting up environment variables with travis:
#     https://stackoverflow.com/a/44850245
