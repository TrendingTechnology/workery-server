# -*- coding: utf-8 -*-
import logging
import django_rq
from datetime import datetime, timedelta
from dateutil import tz
from django.conf import settings
from django.contrib.auth.models import Group
from django.contrib.auth import authenticate
from django.db.models import Q, Prefetch
from django.utils.translation import ugettext_lazy as _
from django.utils import timezone
from django.utils.http import urlquote
from rest_framework import exceptions, serializers
from rest_framework.response import Response
from rest_framework.validators import UniqueValidator
from shared_foundation.models.franchise import SharedFranchise


logger = logging.getLogger(__name__)


class SharedFranchiseListCreateSerializer(serializers.ModelSerializer):

    # OVERRIDE THE MODEL FIELDS AND ENFORCE THE FOLLOWING CUSTOM VALIDATION RULES.
    schema_name = serializers.CharField(
        required=True,
        allow_blank=False,
        validators=[UniqueValidator(queryset=SharedFranchise.objects.all())],
    )

    postal_code = serializers.CharField(
        required=True,
        allow_blank=False,
    )

    name = serializers.CharField(
        required=True,
        allow_blank=False,
    )

    alternate_name = serializers.CharField(
        required=True,
        allow_blank=False,
    )

    class Meta:
        model = SharedFranchise
        fields = (
            # Thing
            'created',
            'last_modified',
            'alternate_name',
            'description',
            'name',
            'url',

            # # ContactPoint
            # 'area_served',
            # 'available_language',
            # 'contact_type',
            # 'email',
            # 'fax_number',
            # 'hours_available',
            # 'product_supported',
            # 'telephone',
            # 'telephone_type_of',
            # 'telephone_extension',
            # 'other_telephone',
            # 'other_telephone_type_of',
            # 'other_telephone_extension',

            # Postal ddress
            'address_country',
            'address_locality',
            'address_region',
            'post_office_box_number',
            'postal_code',
            'street_address',
            'street_address_extra',

            # Tenancy
            'schema_name'
        )

    def setup_eager_loading(cls, queryset):
        """ Perform necessary eager loading of data. """
        queryset = queryset.prefetch_related()
        return queryset

    def create(self, validated_data):
        """
        Override the `create` function to add extra functinality.
        """
        #-----------------------------
        # Get our inputs.
        #-----------------------------
        alternate_name = validated_data.get('alternate_name', None)
        name = validated_data.get('name', None)
        description = validated_data.get('description', False)
        address_country = validated_data.get('address_country', None)
        address_locality = validated_data.get('address_locality', None)
        address_region = validated_data.get('address_region', None)
        post_office_box_number = validated_data.get('post_office_box_number', None)
        postal_code = validated_data.get('postal_code', None)
        street_address = validated_data.get('street_address', None)
        street_address_extra = validated_data.get('street_address_extra', None)
        schema_name = validated_data.get('schema_name', None)
        until_date = validated_data.get('until_date', None)
        until_date = validated_data.get('until_date', None)
        logger.info("Input data:", str(validated_data))

        # #-----------------------------
        # # Create our `AwayLog` object.
        # #-----------------------------
        from shared_franchise.tasks import create_franchise_func
        django_rq.enqueue(create_franchise_func, validated_data)

        # # Create our log.
        # log = AwayLog.objects.create(
        #     associate=associate,
        #     reason=reason,
        #     until_further_notice=until_further_notice,
        #     until_date=until_date,
        #     created_by=self.context['created_by'],
        #     last_modified_by=self.context['created_by'],
        # )
        # logger.info("Created AwayLog")
        #
        # # Save our away information to the associate.
        # associate.away_log = log
        # associate.save()
        # logger.info("Assigned AwayLog to associate.")
        #
        # # Return our validated data.
        # validated_data['id'] = log.id
        return validated_data
