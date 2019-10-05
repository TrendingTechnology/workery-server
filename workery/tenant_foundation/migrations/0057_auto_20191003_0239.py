# Generated by Django 2.0.13 on 2019-10-03 02:39

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('tenant_foundation', '0056_workorderinvoice_revision_version'),
    ]

    operations = [
        migrations.AlterField(
            model_name='workorderinvoice',
            name='invoice_associate_tax',
            field=models.CharField(blank=True, help_text='The associate tax number for this invoice document.', max_length=18, null=True, verbose_name='Invoice Associate Tax #'),
        ),
    ]