# Generated by Django 2.0.13 on 2019-09-22 14:19

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('tenant_foundation', '0050_privatefileupload'),
    ]

    operations = [
        migrations.RenameField(
            model_name='privatefileupload',
            old_name='binary_file',
            new_name='data_file',
        ),
    ]
