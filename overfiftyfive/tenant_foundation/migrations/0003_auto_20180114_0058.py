# Generated by Django 2.0 on 2018-01-14 00:58

from django.db import migrations
import djmoney.models.fields


class Migration(migrations.Migration):

    dependencies = [
        ('tenant_foundation', '0002_auto_20180114_0052'),
    ]

    operations = [
        migrations.AlterField(
            model_name='order',
            name='service_fee',
            field=djmoney.models.fields.MoneyField(blank=True, decimal_places=2, default=None, default_currency='CAD', help_text='The service fee that the customer was charged by the associate..', max_digits=10, null=True, verbose_name='Service Fee'),
        ),
    ]
