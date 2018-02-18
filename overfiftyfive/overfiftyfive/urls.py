"""overfiftyfive URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/1.11/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  url(r'^$', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  url(r'^$', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.conf.urls import url, include
    2. Add a URL to urlpatterns:  url(r'^blog/', include('blog.urls'))
"""
from django.conf import settings
from django.conf.urls import include, url
from django.conf.urls.static import static
from django.contrib import admin
from django.conf.urls.i18n import i18n_patterns
from django.contrib.sitemaps.views import sitemap
# from overfiftyfive.sitemaps import StaticViewSitemap


# sitemaps = {
#     'static': StaticViewSitemap,
# }

# Custom errors.
# handler403 = "shared_foundation.views.http_403_page"
handler404 = "shared_home.views.page_not_found"
handler500 = "shared_home.views.server_error"


# Base URLs.
urlpatterns = [
    url(r'^admin/', admin.site.urls),
    url(r'^i18n/', include('django.conf.urls.i18n')),
    # url('^', include('django.contrib.auth.urls')),
    url(r'^', include('shared_api.urls')),
    url(r'^', include('shared_foundation.urls')),

    #  # Sitemap
    # url(r'^sitemap\.xml$', sitemap, {'sitemaps': sitemaps}, name='django.contrib.sitemaps.views.sitemap'),
    #
    # Django-RQ
    url(r'^django-rq/', include('django_rq.urls')),
]

# Serving static and media files during development
# urlpatterns += static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
# urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)

# Serving our URLs.
urlpatterns += i18n_patterns(
    # Public specific URLs.
    url(r'^', include('shared_api.urls')),
    url(r'^', include('shared_auth.urls')),
    url(r'^', include('shared_home.urls')),

    # Tenant specific URLs.
    url(r'^', include('tenant_api.urls')),
    url(r'^', include('tenant_associate.urls')),
    url(r'^', include('tenant_customer.urls')),
    url(r'^', include('tenant_dashboard.urls')),
    url(r'^', include('tenant_order.urls')),
)
