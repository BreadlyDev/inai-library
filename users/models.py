from django.db import models
from django.contrib.auth.models import AbstractUser, BaseUserManager


ROLES = (
    ("Admin", "Admin"),
    ("Librarian", "Librarian"),
    ("Student", "Student")
)


# def validate_phone(phone):
#     if phone[:1].isdigit() and 8 < len(phone) < 10:
#         return True
#     return False


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
    firstname = models.CharField(max_length=150, unique=True)
    lastname = models.CharField(max_length=150, unique=True)
    email = models.EmailField(unique=True)
    phone = models.CharField(max_length=15)
    role = models.CharField(max_length=150, choices=ROLES, default=ROLES[2][1])
    group = models.ForeignKey(Group, on_delete=models.SET_NULL, null=True, blank=True, unique=False)

    username = None
    date_joined = None
    last_login = None
    groups = None
    user_permissions = None
    first_name = None
    last_name = None

    USERNAME_FIELD = 'email'
    REQUIRED_FIELDS = []

    class Meta:
        db_table = "users"

    objects = CustomUserManager()

    def __str__(self):
        return f"{self.role} {self.first_name} {self.last_name}"
