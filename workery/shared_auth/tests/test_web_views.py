# -*- coding: utf-8 -*-
from datetime import date, datetime, timedelta
from django.core.management import call_command
from django_tenants.test.cases import TenantTestCase
from django_tenants.test.client import TenantClient
from django.urls import reverse
from django.utils import timezone
from rest_framework import status
from shared_foundation.models import SharedFranchise
from shared_foundation.models import SharedUser
from shared_foundation.models import SharedUser
from shared_foundation.utils import get_jwt_token_and_orig_iat


TEST_USER_EMAIL = "bart@workery.ca"
TEST_USER_USERNAME = "bart@workery.ca"
TEST_USER_PASSWORD = "123P@$$w0rd"
TEST_USER_TEL_NUM = "123 123-1234"
TEST_USER_TEL_EX_NUM = ""
TEST_USER_CELL_NUM = "123 123-1234"


class TestSharedAuthWebViews(TenantTestCase):
    """
    Class used to test the web views.

    Console:
    python manage.py test shared_auth.tests.test_web_views
    """

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
        super(TestSharedAuthWebViews, self).setUp()

        # Setup our app and account.
        call_command('init_app', verbosity=0)
        call_command(
           'create_shared_account',
           TEST_USER_EMAIL,
           TEST_USER_PASSWORD,
           "Bart",
           "Mika",
           verbosity=0
        )

        # Get user and credentials.
        user = SharedUser.objects.get()
        token, orig_iat = get_jwt_token_and_orig_iat(user)

        # Setup our clients.
        self.anon_c = TenantClient(self.tenant)
        self.auth_c = TenantClient(self.tenant, HTTP_AUTHORIZATION='JWT {0}'.format(token))
        self.auth_c.login(
            username=TEST_USER_USERNAME,
            password=TEST_USER_PASSWORD
        )

        # Attach our user(s) to our test tenant organization.
        user.franchise = self.tenant
        user.save()

    def tearDown(self):
        """
        Run at the end of every unit test.
        """
        # Delete our clients.
        del self.anon_c
        del self.auth_c

        # Delete previous data.
        SharedUser.objects.all().delete()

        # Finish teardown.
        super(TestSharedAuthWebViews, self).tearDown()

    def test_get_index_page(self):
        response = self.anon_c.get(reverse('workery_login_master'))
        self.assertEqual(response.status_code, status.HTTP_200_OK)

    def test_user_login_redirector_page_with_anonymous_user(self):
        response = self.anon_c.get(reverse('workery_login_redirector'))
        self.assertEqual(response.status_code, status.HTTP_302_FOUND)

    def test_user_login_redirector_page_with_authenticated_user(self):
        response = self.auth_c.get(
            reverse('workery_login_redirector'),
        )
        self.assertEqual(response.status_code, status.HTTP_302_FOUND)

    # def test_send_reset_password_email_master_page(self):
    #     response = self.anon_c.get(reverse('workery_send_reset_password_email_master'))
    #     self.assertEqual(response.status_code, status.HTTP_200_OK)
    #
    # def test_send_reset_password_email_submitted_page(self):
    #     response = self.anon_c.get(reverse('workery_send_reset_password_email_submitted'))
    #     self.assertEqual(response.status_code, status.HTTP_200_OK)
    #
    # def test_rest_password_master_page_with_success(self):
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #     url = reverse('workery_reset_password_master', args=[me.pr_access_code])
    #     response = self.anon_c.get(url)
    #     self.assertEqual(response.status_code, status.HTTP_200_OK)
    #
    # def test_rest_password_master_page_with_bad_pr_access_code(self):
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #     url = reverse('workery_reset_password_master', args=['some-bad-pr-access-code'])
    #     response = self.anon_c.get(url)
    #     self.assertEqual(response.status_code, status.HTTP_302_FOUND)
    #
    # def test_rest_password_master_page_with_expired_pr_access_code(self):
    #     # Get the user profile.
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #
    #     # Set the expiry date to be old!
    #     today = timezone.now()
    #     today_minus_1_year = today - timedelta(minutes=1)
    #     me.pr_expiry_date = today_minus_1_year
    #     me.save()
    #
    #     # Run our test...
    #     url = reverse('workery_reset_password_master', args=[me.pr_access_code])
    #     response = self.anon_c.get(url)
    #
    #     # Verify the results.
    #     self.assertEqual(response.status_code, status.HTTP_302_FOUND)
    #
    # def test_user_activation_detail_page_with_success(self):
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #     url = reverse('workery_user_activation_detail', args=[me.pr_access_code])
    #     response = self.anon_c.get(url)
    #     self.assertEqual(response.status_code, status.HTTP_200_OK)
    #
    # def test_rest_user_activation_detail_page_with_bad_pr_access_code(self):
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #     url = reverse('workery_user_activation_detail', args=['some-bad-pr-access-code'])
    #     response = self.anon_c.get(url)
    #     self.assertNotEqual(response.status_code, status.HTTP_200_OK)
    #
    # def test_rest_user_activation_detail_page_with_expired_pr_access_code(self):
    #     # Get the user profile.
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #
    #     # Set the expiry date to be old!
    #     today = timezone.now()
    #     today_minus_1_year = today - timedelta(minutes=1)
    #     me.pr_expiry_date = today_minus_1_year
    #     me.save()
    #
    #     # Run our test...
    #     url = reverse('workery_user_activation_detail', args=[me.pr_access_code])
    #     response = self.anon_c.get(url)
    #
    #     # Verify the results.
    #     self.assertNotEqual(response.status_code, status.HTTP_200_OK)
    #
    # def test_user_logout_redirector_master_page_with_redirect(self):
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #     url = reverse('workery_logout_redirector')
    #     response = self.anon_c.get(url)
    #     self.assertEqual(response.status_code, 302)
    #
    # def test_user_logout_redirector_master_page_with_success(self):
    #     me = SharedUser.objects.get(email=TEST_USER_EMAIL)
    #     url = reverse('workery_logout_redirector')
    #     response = self.auth_c.get(url)
    #     self.assertEqual(response.status_code, 302)
