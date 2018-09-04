# -*- coding: utf-8 -*-
from datetime import date, datetime, timedelta
from django.core.management import call_command
from django_tenants.test.cases import TenantTestCase
from django_tenants.test.client import TenantClient
from django.urls import reverse
from django.utils import timezone
from rest_framework import status

from shared_foundation import constants
from shared_foundation.models import SharedFranchise
from shared_foundation.models import SharedUser
from shared_foundation.utils import get_jwt_token_and_orig_iat
from tenant_foundation.models import Associate


TEST_SCHEMA_NAME = "london"
TEST_USER_EMAIL = "bart@workery.ca"
TEST_USER_USERNAME = "bart@workery.ca"
TEST_USER_PASSWORD = "123P@$$w0rd"
TEST_USER_TEL_NUM = "123 123-1234"
TEST_USER_TEL_EX_NUM = ""
TEST_USER_CELL_NUM = "123 123-1234"


class TestTenantAssociateViews(TenantTestCase):
    """
    Console:
    python manage.py test tenant_associate.tests
    """

    #------------------#
    # Setup Unit Tests #
    #------------------#

    def setup_tenant(tenant):
        """Tenant Schema"""
        tenant.schema_name = TEST_SCHEMA_NAME
        tenant.name='Over 55 (London) Inc.',
        tenant.alternate_name="Over55",
        tenant.description="Located at the Forks of the Thames in ...",
        tenant.address_country="CA",
        tenant.address_locality="London",
        tenant.address_region="Ontario",
        tenant.post_office_box_number="", # Post Offic #
        tenant.postal_code="N6H 1B4",
        tenant.street_address="78 Riverside Drive",
        tenant.street_address_extra="", # Extra line.

    @classmethod
    def setUpClass(cls):
        """
        Run at the beginning before all the unit tests run.
        """
        super().setUpClass()

    def setUp(self):
        """
        Run at the beginning of every unit test.
        """
        super(TestTenantAssociateViews, self).setUp()

        # Setup our app and account.
        call_command('init_app', verbosity=0)
        call_command('populate_tenant_content', TEST_SCHEMA_NAME, verbosity=0)
        call_command('populate_tenant_sample_db', TEST_SCHEMA_NAME, verbosity=0)

        # Get user and credentials.
        user = SharedUser.objects.get(email="bart+executive@workery.ca")
        token, orig_iat = get_jwt_token_and_orig_iat(user)

        # Setup our clients.
        self.anon_c = TenantClient(self.tenant)
        self.auth_c = TenantClient(self.tenant, HTTP_AUTHORIZATION='JWT {0}'.format(token))
        self.auth_c.login(
            username = "bart+executive@workery.ca",
            password = "123P@$$w0rd"
        )

    def tearDown(self):
        """
        Run at the end of every unit test.
        """
        # Delete previous data.
        SharedUser.objects.all().delete()

        # Delete our clients.
        del self.anon_c
        del self.auth_c

        # Finish teardown.
        super(TestTenantAssociateViews, self).tearDown()

    def test_associate_summary_page(self):
        response = self.auth_c.get(self.tenant.reverse('workery_tenant_member_summary'))
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # self.assertIn('Staff', str(response.content))

    def test_associate_member_create_page(self):
        response = self.auth_c.get(self.tenant.reverse('workery_tenant_member_create'))
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # self.assertIn('Staff', str(response.content))

    def test_associate_member_confirm_create_page(self):
        response = self.auth_c.get(self.tenant.reverse('workery_tenant_member_confirm_create'))
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # self.assertIn('Staff', str(response.content))

    def test_associate_member_list_page(self):
        response = self.auth_c.get(self.tenant.reverse('workery_tenant_member_list'))
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # self.assertIn('Staff', str(response.content))

    def test_associate_member_search_page(self):
        response = self.auth_c.get(self.tenant.reverse('workery_tenant_member_search'))
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # self.assertIn('Staff', str(response.content))

    def test_associate_member_search_results_page(self):
        response = self.auth_c.get(self.tenant.reverse('workery_tenant_member_search_results')+"?keyword=test")
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # self.assertIn('Staff', str(response.content))

    def test_associate_lite_retrieve_page(self):
        associate = Associate.objects.all().first()
        a_url = self.tenant.reverse(
            reverse_id='workery_tenant_member_lite_retrieve',
            reverse_args=['summary', int(associate.id)]
        )
        response = self.auth_c.get(a_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_associate_member_full_retrieve_page(self):
        associate = Associate.objects.all().first()
        a_url = self.tenant.reverse(
            reverse_id='workery_tenant_member_full_retrieve',
            reverse_args=['summary', int(associate.id)]
        )
        response = self.auth_c.get(a_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_associate_member_retrieve_for_comment_list_page(self):
        associate = Associate.objects.all().first()
        a_url = self.tenant.reverse(
            reverse_id='workery_tenant_member_retrieve_for_comment_list',
            reverse_args=['summary', int(associate.id)]
        )
        response = self.auth_c.get(a_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_associate_member_retrieve_for_activity_sheet_list_page(self):
        associate = Associate.objects.all().first()
        a_url = self.tenant.reverse(
            reverse_id='workery_tenant_member_retrieve_for_activity_sheet_list',
            reverse_args=['summary', int(associate.id)]
        )
        response = self.auth_c.get(a_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_associate_member_retrieve_for_jobs_list_page(self):
        associate = Associate.objects.all().first()
        a_url = self.tenant.reverse(
            reverse_id='workery_tenant_member_retrieve_for_jobs_list',
            reverse_args=['summary', int(associate.id)]
        )
        response = self.auth_c.get(a_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_associate_member_update_page(self):
        associate = Associate.objects.all().first()
        a_url = self.tenant.reverse(
            reverse_id='workery_tenant_member_update',
            reverse_args=['summary', int(associate.id)]
        )
        response = self.auth_c.get(a_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
