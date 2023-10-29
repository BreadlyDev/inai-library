from django.db import models
from django.core.validators import ValidationError
from django.utils import timezone
from datetime import timedelta
from users.models import User
from books.models import Book


def validate_due_date(value):
    current_date = timezone.now().date()
    next_month = current_date.replace(day=1) + timedelta(days=32)
    if value.date() > next_month:
        raise ValidationError("Дата сдачи не может быть позже конца следующего месяца.")


ORDER_STATUS = (
    ("Не рассмотрено", "Не рассмотрено"),
    ("Подтверждено", "Подтверждено"),
    ("Обрабатывается", "Обрабатывается"),
    ("Готово", "Готово"),
    ("Ошибка", "Ошибка"),
)


class Order(models.Model):
    owner = models.ForeignKey(User, on_delete=models.CASCADE)
    books = models.ManyToManyField(Book)
    status = models.CharField(max_length=50, choices=ORDER_STATUS, default=ORDER_STATUS[0][0])
    comment = models.TextField(default="", blank=True)
    created_time = models.DateTimeField(auto_now_add=True)
    due_time = models.DateTimeField(validators=[validate_due_date])

    class Meta:
        ordering = ["-created_time"]
        db_table = "orders"

    def __str__(self):
        return f"Order by {self.owner} at {self.created_time}"
