# Generated by Django 2.0.4 on 2018-05-05 23:14

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('tenant_foundation', '0007_auto_20180505_2106'),
    ]

    operations = [
        migrations.RenameField(
            model_name='taskitem',
            old_name='created',
            new_name='created_at',
        ),
        migrations.RenameField(
            model_name='taskitem',
            old_name='last_modified',
            new_name='last_modified_at',
        ),
    ]
