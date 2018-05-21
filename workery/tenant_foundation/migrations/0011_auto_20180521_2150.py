# Generated by Django 2.0.4 on 2018-05-21 21:50

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('tenant_foundation', '0010_remove_order_payment_date'),
    ]

    operations = [
        migrations.AddField(
            model_name='associate',
            name='is_archived',
            field=models.BooleanField(db_index=True, default=True, help_text='Indicates whether associate was archived.', verbose_name='Is Archived'),
        ),
        migrations.AddField(
            model_name='customer',
            name='is_archived',
            field=models.BooleanField(db_index=True, default=True, help_text='Indicates whether customer was archived.', verbose_name='Is Archived'),
        ),
        migrations.AddField(
            model_name='order',
            name='is_archived',
            field=models.BooleanField(db_index=True, default=True, help_text='Indicates whether order was archived.', verbose_name='Is Archived'),
        ),
        migrations.AddField(
            model_name='partner',
            name='is_archived',
            field=models.BooleanField(db_index=True, default=True, help_text='Indicates whether partner was archived.', verbose_name='Is Archived'),
        ),
        migrations.AddField(
            model_name='staff',
            name='is_archived',
            field=models.BooleanField(db_index=True, default=True, help_text='Indicates whether staff was archived.', verbose_name='Is Archived'),
        ),
    ]
