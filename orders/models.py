from django.db import models
from django.core.validators import ValidationError
from django.utils import timezone
from datetime import timedelta, date
from users.models import User
from books.models import Book


def validate_due_date(value):
    current_date = date.today()
    next_month = current_date.replace(day=1) + timedelta(days=32)
    if value > next_month:
        raise ValidationError("Дата сдачи не может быть позже конца следующего месяца.")


def validate_inventory_number(inventory_number):
    if (inventory_number[:7] != "INAI.KG"
            and inventory_number[7:].isdigit()
            or inventory_number[:5] != "КГФИ."
            and inventory_number[5:].isdigit()
    ):
        raise ValidationError("Неправильный формат инвентарного номера. "
                              "Инвентарный номер должен начинаться с"
                              "INAI.KG или КГФИ. и заканчиваться цифрами")



ORDER_STATUS = (
    ("Ожидает проверки", "Ожидает проверки"),
    ("В обработке", "В обработке"),
    ("Выполнен", "Выполнен"),
    ("Отклонено", "Отклонено"),
)


class Order(models.Model):
    owner = models.ForeignKey(User, on_delete=models.CASCADE)
    book = models.ForeignKey(Book, on_delete=models.CASCADE)
    inventory_number = models.CharField(max_length=150, null=True, blank=True, validators=[validate_inventory_number])
    status = models.CharField(max_length=50, choices=ORDER_STATUS, default=ORDER_STATUS[0][0])
    comment = models.TextField(default="", blank=True)
    created_time = models.DateTimeField(auto_now_add=True)
    due_time = models.DateField(validators=[validate_due_date])

    class Meta:
        ordering = ["-created_time"]
        db_table = "orders"
        verbose_name = "Order"
        verbose_name_plural = "Orders"

    def formatted_created_time(self):
        local_time = timezone.localtime(self.created_time)
        return local_time.strftime("%d-%m-%y %H:%M:%S")

    def __str__(self):
        return f"Order by {self.owner} at {self.formatted_created_time()}"
