# Generated by Django 2.0.13 on 2019-12-19 00:23

from django.conf import settings
from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
        ('tenant_foundation', '0071_auto_20191109_1857'),
    ]

    operations = [
        migrations.CreateModel(
            name='UnifiedSearchItem',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('text', models.CharField(blank=True, db_index=True, help_text='The searchable content text used by the keyword searcher function.', max_length=511, null=True, unique=True, verbose_name='Text')),
                ('created_from', models.GenericIPAddressField(blank=True, help_text='The IP address of the creator.', null=True, verbose_name='Created from')),
                ('created_from_is_public', models.BooleanField(default=False, help_text='Is creator a public IP and is routable.', verbose_name='Is the IP ')),
                ('last_modified_from', models.GenericIPAddressField(blank=True, help_text='The IP address of the modifier.', null=True, verbose_name='Last modified from')),
                ('last_modified_from_is_public', models.BooleanField(default=False, help_text='Is modifier a public IP and is routable.', verbose_name='Is the IP ')),
                ('associate', models.OneToOneField(blank=True, help_text='The associate of this search item.', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='unified_search_item', to='tenant_foundation.Associate')),
                ('created_by', models.ForeignKey(blank=True, help_text='The user whom created this object.', null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='created_unified_search_items', to=settings.AUTH_USER_MODEL)),
                ('customer', models.OneToOneField(blank=True, help_text='The customer of this search item.', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='unified_search_item', to='tenant_foundation.Customer')),
                ('file', models.OneToOneField(blank=True, help_text='The file of this search item.', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='unified_search_item', to='tenant_foundation.PrivateFileUpload')),
                ('job', models.OneToOneField(blank=True, help_text='The work-order of this search item.', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='unified_search_item', to='tenant_foundation.WorkOrder')),
                ('last_modified_by', models.ForeignKey(blank=True, help_text='The user whom modified this object last.', null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='last_modified_unified_search_items', to=settings.AUTH_USER_MODEL)),
                ('partner', models.OneToOneField(blank=True, help_text='The partner of this search item.', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='unified_search_item', to='tenant_foundation.Partner')),
                ('staff', models.OneToOneField(blank=True, help_text='The staff of this search item.', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='unified_search_item', to='tenant_foundation.Staff')),
                ('tags', models.ManyToManyField(blank=True, help_text='The tags with this unified search item.', related_name='unified_search_items', to='tenant_foundation.Tag')),
            ],
            options={
                'verbose_name': 'Unified Search Item',
                'verbose_name_plural': 'Unified Search Items',
                'db_table': 'workery_unified_search_items',
                'permissions': (),
                'default_permissions': (),
            },
        ),
    ]
