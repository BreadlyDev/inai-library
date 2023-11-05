from django.db import models
from django.contrib.auth.models import AbstractUser, BaseUserManager


USER_STATUS = (
    ("Admin", "Admin"),
    ("Librarian", "Librarian"),
    ("Student", "Student")
)


def validate_phone(phone):
    if phone[:1].isdigit() and 8 < len(phone) < 10:
        return True
    return False


class CustomUserManager(BaseUserManager):
    def create_user(self, email, password=None, **extra_fields):
        if not email:
            raise ValueError("Email is required field")
        email = self.normalize_email(email)
        status = extra_fields.get("status", "Student")

        if status not in dict(USER_STATUS).keys():
            raise ValueError("Invalid user status")

        user = self.model(email=email, **extra_fields)
        user.status = status
        user.set_password(password)
        user.save(using=self._db)
        return user

    def create_superuser(self, email, password=None, **extra_fields):
        extra_fields.setdefault("is_staff", True)
        extra_fields.setdefault("is_superuser", True)
        extra_fields["status"] = "Admin"

        return self.create_user(email, password, **extra_fields)


class User(AbstractUser):

    username = None
    password = models.CharField(max_length=150, null=True)
    first_name = models.CharField(max_length=150, default="", blank=True)
    last_name = models.CharField(max_length=150, default="", blank=True)
    email = models.EmailField(unique=True, null=True)
    phone = models.CharField(max_length=150, default="", blank=True, validators=[validate_phone])
    status = models.CharField(max_length=150, choices=USER_STATUS, default=USER_STATUS[2][0])
    date_joined = None

    USERNAME_FIELD = 'email'
    REQUIRED_FIELDS = []

    groups = models.ManyToManyField(
        'auth.Group',
        related_name='custom_users',
        blank=True,
        help_text='The groups this user belongs to. A user will get all permissions granted to each of their groups.',
        verbose_name='groups',
    )

    user_permissions = models.ManyToManyField(
        'auth.Permission',
        related_name='custom_users',
        blank=True,
        help_text='Specific permissions for this user.',
        verbose_name='user permissions',
    )

    class Meta:
        db_table = "users"

    objects = CustomUserManager()
