# Generated by Django 2.0.4 on 2018-05-05 19:21

from django.db import migrations, models
import tenant_foundation.models.order


class Migration(migrations.Migration):

    dependencies = [
        ('tenant_foundation', '0001_initial'),
    ]

    operations = [
        migrations.AddField(
            model_name='order',
            name='start_date',
            field=models.DateField(blank=True, default=tenant_foundation.models.order.get_todays_date, help_text='The date that this order will begin.', verbose_name='Start Date'),
        ),
    ]
