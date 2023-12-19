from django.db import models
from django.contrib.auth.models import AbstractUser, BaseUserManager
from rest_framework.exceptions import ValidationError

ROLES = (
    ("Admin", "Admin"),
    ("Librarian", "Librarian"),
    ("Student", "Student")
)


def validate_phone(phone):
    if not phone.isdigit():
        return ValidationError("Номер телефона должен состоять из цифр")


class Group(models.Model):
    name = models.CharField(max_length=150)

    class Meta:
        db_table = "groups"

    def __str__(self):
        return f"{self.name} group"


class CustomUserManager(BaseUserManager):
    def create_user(self, email, password=None, **extra_fields):
        if not email:
            raise ValueError("Email is required field")

        email = self.normalize_email(email)
        role = extra_fields.get("role", "Student")

        if role not in dict(ROLES).keys():
            raise ValueError("Invalid user status")

        user = self.model(email=email, **extra_fields)
        user.role = role
        user.set_password(password)
        user.save(using=self._db)
        return user

    def create_superuser(self, email, password=None, **extra_fields):
        extra_fields.setdefault("is_staff", True)
        extra_fields.setdefault("is_superuser", True)
        extra_fields["role"] = "Admin"

        return self.create_user(email, password, **extra_fields)


class User(AbstractUser):
    password = models.CharField(max_length=128)
    firstname = models.CharField(max_length=150)
    lastname = models.CharField(max_length=150)
    email = models.EmailField(unique=True)
    phone = models.CharField(max_length=15, validators=[validate_phone])
    role = models.CharField(max_length=150, choices=ROLES, default=ROLES[2][1])
    group = models.ForeignKey(Group, on_delete=models.SET_NULL, null=True, blank=True, unique=False)

    username = None
    date_joined = None
    last_login = None
    groups = None
    user_permissions = None
    first_name = None
    last_name = None

    USERNAME_FIELD = "email"
    REQUIRED_FIELDS = []

    class Meta:
        db_table = "users"

    objects = CustomUserManager()

    def __str__(self):
        return f"{self.role} {self.firstname} {self.lastname}"

    def save(self, *args, **kwargs):
        if not self.password:
            self.set_password(self.password)
